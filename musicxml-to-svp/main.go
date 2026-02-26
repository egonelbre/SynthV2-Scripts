package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/egonelbre/synthv2-scripts/internal/phonemes"
	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

const blicksPerQuarter = 705600000

func pitchToMIDI(p *musicxml.Pitch) (midi int, detune int) {
	stepMap := map[musicxml.Step]int{
		"C": 0, "D": 2, "E": 4, "F": 5, "G": 7, "A": 9, "B": 11,
	}
	semitone := stepMap[p.Step]
	midi = (p.Octave+1)*12 + semitone

	if p.Alter != "" {
		alter, err := strconv.ParseFloat(p.Alter, 64)
		if err == nil {
			wholeSemitones := int(math.Round(alter))
			midi += wholeSemitones
			detune = int(math.Round((alter - float64(wholeSemitones)) * 100))
		}
	}
	return midi, detune
}

func durationToBlicks(duration, divisions int) int64 {
	return int64(duration) * blicksPerQuarter / int64(divisions)
}

func extractLyric(note *musicxml.Note) string {
	for _, lyric := range note.Lyric {
		if len(lyric.Text) > 0 && lyric.Text[0].EnclosedText != "" {
			return lyric.Text[0].EnclosedText
		}
	}
	return ""
}

func beatUnitToQuarters(beatUnit musicxml.NoteTypeValue, hasDot bool) float64 {
	q := 1.0
	switch beatUnit {
	case "whole":
		q = 4.0
	case "half":
		q = 2.0
	case "quarter":
		q = 1.0
	case "eighth":
		q = 0.5
	case "16th":
		q = 0.25
	}
	if hasDot {
		q *= 1.5
	}
	return q
}

func noteTieTypes(note *musicxml.Note) (hasStart, hasStop bool) {
	for _, t := range note.Tie {
		switch t.Type {
		case "start":
			hasStart = true
		case "stop":
			hasStop = true
		}
	}
	return
}

func noteHasArticulation(note *musicxml.Note, name string) bool {
	for _, n := range note.Notations {
		for _, a := range n.Articulations {
			switch name {
			case "staccato":
				if len(a.Staccato) > 0 {
					return true
				}
			case "staccatissimo":
				if len(a.Staccatissimo) > 0 {
					return true
				}
			case "tenuto":
				if len(a.Tenuto) > 0 {
					return true
				}
			case "accent":
				if len(a.Accent) > 0 {
					return true
				}
			case "strong-accent":
				if len(a.StrongAccent) > 0 {
					return true
				}
			}
		}
	}
	return false
}

func parseDuration(s string) int {
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	voiceFlag := flag.String("voice", "", "assign voices: choir1, choir2, choir3, or soloists")
	panFlag := flag.String("pan", "default", "panning scheme: default, spread, center")
	langFlag := flag.String("lang", "", "convert lyrics to phonemes: estonian, karelian")
	relaxedFlag := flag.Bool("relaxed", false, "enable relaxed consonant pronunciation")
	outputFlag := flag.String("o", "", "output file path (default: input with .svp extension)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: musicxml-to-svp [flags] <input.musicxml>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := flag.Arg(0)
	outputPath := *outputFlag
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = inputPath[:len(inputPath)-len(ext)] + ".svp"
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}

	var score musicxml.ScorePartwise
	if err := xml.Unmarshal(data, &score); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing MusicXML: %v\n", err)
		os.Exit(1)
	}

	// Build part name lookup
	partNames := map[string]string{}
	for _, sp := range score.PartList.GroupScorePart.ScorePart {
		name := ""
		if sp.PartName != nil {
			name = sp.PartName.EnclosedText
		}
		partNames[sp.Id] = name
	}

	// Pass 1: Extract global timing from first part
	meters := []*SVPMeter{}
	tempos := []*SVPTempo{}
	type measureInfo struct {
		startBlicks int64
		divisions   int
	}
	var measureInfos []measureInfo

	if len(score.Part) > 0 {
		firstPart := score.Part[0]
		divisions := 4
		cursor := int64(0)
		measureIdx := 0

		for _, measure := range firstPart.Measure {
			measureInfos = append(measureInfos, measureInfo{
				startBlicks: cursor,
				divisions:   divisions,
			})

			measureDuration := int64(0)
			for _, el := range measure.Element {
				switch value := el.Value.(type) {
				case *musicxml.Attributes:
					if value.Divisions != 0 {
						divisions = value.Divisions
						measureInfos[len(measureInfos)-1].divisions = divisions
					}
					for _, t := range value.Time {
						meters = append(meters, &SVPMeter{
							Index:       measureIdx,
							Numerator:   t.Beats,
							Denominator: t.BeatType,
						})
					}
				case *musicxml.Direction:
					// Extract tempo from Sound element
					if value.Sound != nil && value.Sound.Tempo != "" {
						bpm, err := strconv.ParseFloat(value.Sound.Tempo, 64)
						if err == nil {
							tempos = append(tempos, &SVPTempo{
								Position: cursor + measureDuration,
								BPM:      bpm,
							})
						}
					}
					// Also check metronome in direction-type
					if len(tempos) == 0 || value.Sound == nil {
						for _, dt := range value.DirectionType {
							if dt.Metronome != nil && dt.Metronome.PerMinute != nil {
								bpm, err := strconv.ParseFloat(dt.Metronome.PerMinute.EnclosedText, 64)
								if err == nil {
									q := beatUnitToQuarters(dt.Metronome.BeatUnit, dt.Metronome.BeatUnitDot != "")
									bpm = bpm * q // adjust if beat unit is not quarter
									tempos = append(tempos, &SVPTempo{
										Position: cursor + measureDuration,
										BPM:      bpm,
									})
								}
							}
						}
					}
				case *musicxml.Note:
					if value.Grace != nil {
						continue
					}
					dur := parseDuration(value.Duration)
					if dur > 0 {
						if value.Chord == "" {
							blicks := durationToBlicks(dur, divisions)
							measureDuration += blicks
						}
					}
				case *musicxml.Backup:
					dur := parseDuration(value.Duration)
					if dur > 0 {
						measureDuration -= durationToBlicks(dur, divisions)
					}
				case *musicxml.Forward:
					dur := parseDuration(value.Duration)
					if dur > 0 {
						measureDuration += durationToBlicks(dur, divisions)
					}
				}
			}

			// Advance cursor by the full measure duration
			// Use time signature to compute expected measure duration
			meterNum := 4
			meterDen := 4
			for _, m := range meters {
				if m.Index <= measureIdx {
					meterNum = m.Numerator
					meterDen = m.Denominator
				}
			}
			expectedDuration := int64(meterNum) * blicksPerQuarter * 4 / int64(meterDen)
			cursor += expectedDuration

			measureIdx++
		}
	}

	// Defaults if none found
	if len(meters) == 0 {
		meters = append(meters, &SVPMeter{Index: 0, Numerator: 4, Denominator: 4})
	}
	if len(tempos) == 0 {
		tempos = append(tempos, &SVPTempo{Position: 0, BPM: 120})
	}

	// Pass 2: Notes per part
	var library []*SVPGroup
	var tracks []*SVPTrack

	dispColors := []string{
		"#6699cc", "#cc6699", "#99cc66", "#cc9966",
		"#9966cc", "#66cc99", "#cc6666", "#6666cc",
	}

	for partIdx, part := range score.Part {
		partName := partNames[part.Id]
		if partName == "" {
			partName = fmt.Sprintf("Part %d", partIdx+1)
		}

		divisions := 4
		transpose := 0 // semitone offset from MusicXML <transpose>
		cursor := int64(0)
		var notes []*SVPNote
		var dynEvents []dynEvent
		var accents []accentEvent
		pendingTies := map[int]*SVPNote{} // keyed by MIDI pitch
		var prevOnset int64

		measureIdx := 0
		for _, measure := range part.Measure {
			// Reset cursor to measure start if we have timing info
			if measureIdx < len(measureInfos) {
				cursor = measureInfos[measureIdx].startBlicks
			}

			for _, el := range measure.Element {
				switch value := el.Value.(type) {
				case *musicxml.Attributes:
					if value.Divisions != 0 {
						divisions = value.Divisions
					}
					for _, tr := range value.Transpose {
						chromatic, _ := strconv.Atoi(tr.Chromatic)
						transpose = chromatic + tr.OctaveChange*12
					}
				case *musicxml.Direction:
					for _, dt := range value.DirectionType {
						// Dynamics markings (p, mf, f, ...)
						for _, dyn := range dt.Dynamics {
							if lvl, ok := dynamicsToLevel(dyn); ok {
								dynEvents = append(dynEvents, dynEvent{
									position: cursor,
									kind:     dynLevel,
									loudness: lvl.loudness,
									tension:  lvl.tension,
								})
							}
						}
						// Wedge hairpins (crescendo/diminuendo)
						if dt.Wedge != nil {
							num := dt.Wedge.Number
							if num == 0 {
								num = 1
							}
							switch dt.Wedge.Type {
							case "crescendo":
								dynEvents = append(dynEvents, dynEvent{
									position: cursor,
									kind:     dynCrescStart,
									number:   num,
								})
							case "diminuendo":
								dynEvents = append(dynEvents, dynEvent{
									position: cursor,
									kind:     dynDimStart,
									number:   num,
								})
							case "stop":
								dynEvents = append(dynEvents, dynEvent{
									position: cursor,
									kind:     dynWedgeStop,
									number:   num,
								})
							}
						}
						// Text-based cresc./dim.
						for _, w := range dt.Words {
							if isTextCresc(w.EnclosedText) {
								dynEvents = append(dynEvents, dynEvent{
									position: cursor,
									kind:     dynCrescStart,
									number:   -1, // unpaired text marking
								})
							} else if isTextDim(w.EnclosedText) {
								dynEvents = append(dynEvents, dynEvent{
									position: cursor,
									kind:     dynDimStart,
									number:   -1,
								})
							}
						}
					}
				case *musicxml.Note:
					if value.Grace != nil {
						continue
					}

					dur := parseDuration(value.Duration)
					if dur <= 0 {
						continue
					}
					blicks := durationToBlicks(dur, divisions)

					// Chord: use previous note's onset
					onset := cursor
					if value.Chord != "" {
						onset = prevOnset
					}

					if value.Rest != nil {
						// Rest: advance cursor only
						if value.Chord == "" {
							prevOnset = cursor
							cursor += blicks
						}
						continue
					}

					if value.Pitch == nil {
						// No pitch, no rest — skip (e.g. unpitched)
						if value.Chord == "" {
							prevOnset = cursor
							cursor += blicks
						}
						continue
					}

					midi, detune := pitchToMIDI(value.Pitch)
					midi += transpose
					tieStart, tieStop := noteTieTypes(value)

					if tieStop {
						// Extend pending tied note
						if pending, ok := pendingTies[midi]; ok {
							pending.Duration += blicks
							if !tieStart {
								delete(pendingTies, midi)
							}
						}
						if value.Chord == "" {
							prevOnset = cursor
							cursor += blicks
						}
						continue
					}

					lyric := extractLyric(value)

					note := &SVPNote{
						Onset:    onset,
						Duration: blicks,
						Lyrics:   lyric,
						Phonemes: "",
						Pitch:    midi,
						Detune:   detune,
						Takes: SVPTakes{
							EvenSyllableDuration: true,
							SystemPitchDelta:     SVPParamMode{Mode: "cubic"},
							Takes: []SVPTake{{
								ID:    0,
								Liked: false,
								Seeds: SVPSeeds{},
							}},
						},
					}
					// Apply articulation adjustments.
					hasTenuto := noteHasArticulation(value, "tenuto")
					if !hasTenuto {
						if noteHasArticulation(value, "staccatissimo") {
							note.Duration = note.Duration / 4
						} else if noteHasArticulation(value, "staccato") {
							note.Duration = note.Duration / 2
						}
					}
					if noteHasArticulation(value, "strong-accent") {
						accents = append(accents, accentEvent{
							position: onset,
							duration: blicks,
							strong:   true,
						})
					} else if noteHasArticulation(value, "accent") {
						accents = append(accents, accentEvent{
							position: onset,
							duration: blicks,
						})
					}

					notes = append(notes, note)

					if tieStart {
						pendingTies[midi] = note
					}

					if value.Chord == "" {
						prevOnset = cursor
						cursor += blicks
					}

				case *musicxml.Backup:
					dur := parseDuration(value.Duration)
					if dur > 0 {
						cursor -= durationToBlicks(dur, divisions)
					}
				case *musicxml.Forward:
					dur := parseDuration(value.Duration)
					if dur > 0 {
						cursor += durationToBlicks(dur, divisions)
					}
				}
			}
			measureIdx++
		}

		// Mark melismatic continuation notes with "-".
		for _, n := range notes {
			if n.Lyrics == "" {
				n.Lyrics = "-"
			}
		}

		// Build loudness curve from collected dynamic events
		params := newEmptyParameters()
		if len(dynEvents) > 0 {
			params.Loudness.Points = buildCurve(dynEvents, func(e dynEvent) float64 { return e.loudness }, 6)
			params.Tension.Points = buildCurve(dynEvents, func(e dynEvent) float64 { return e.tension }, 0.15)
		}
		if len(accents) > 0 {
			params.Loudness.Points = applyAccents(params.Loudness.Points, accents, 1.5, 3)
			params.Tension.Points = applyAccents(params.Tension.Points, accents, 0.15, 0.3)
		}

		group := &SVPGroup{
			Name:       partName,
			UUID:       newUUID(),
			Notes:      notes,
			Parameters: params,
		}
		library = append(library, group)

		// Create main group (empty, for the track's own data)
		mainGroup := &SVPGroup{
			Name:       "main",
			UUID:       newUUID(),
			Notes:      []*SVPNote{},
			Parameters: newEmptyParameters(),
		}

		color := dispColors[partIdx%len(dispColors)]

		track := &SVPTrack{
			Name:      partName,
			DispColor: color,
			DispOrder: partIdx,
			UUID:      newUUID(),
			Mixer: SVPMixer{
				GainDecibel: 0,
				Pan:         0,
				Display:     true,
			},
			MainGroup: SVPGroupRef{
				GroupID: mainGroup.UUID,
				UUID:    newUUID(),
			},
			MainRef: SVPGroupRef{
				GroupID: mainGroup.UUID,
				UUID:    newUUID(),
			},
			Groups: []SVPGroupRef{{
				GroupID: group.UUID,
				UUID:    newUUID(),
			}},
		}
		tracks = append(tracks, track)
		library = append(library, mainGroup)
	}

	// Assign voices if requested.
	if *voiceFlag != "" {
		assignVoices(tracks, strings.ToLower(*voiceFlag), *relaxedFlag, *panFlag)
		setNoteAttributes(library)
	}

	// Convert lyrics to phonemes if requested.
	if *langFlag != "" {
		conv := phonemes.New(*langFlag)
		if conv == nil {
			fmt.Fprintf(os.Stderr, "unknown language: %q (options: estonian, karelian)\n", *langFlag)
			os.Exit(1)
		}
		applyPhonemes(library, conv)
	}

	project := SVPProject{
		Version: 196,
		UUID:    newUUID(),
		Time: SVPTime{
			Meters: meters,
			Tempos: tempos,
		},
		Library: library,
		Tracks:  tracks,
		RenderConfig: SVPRenderConfig{
			Destination:      "",
			Filename:         "",
			NumChannels:      1,
			AspirationFormat: "noAspiration",
			BitDepth:         16,
			SampleRate:       44100,
			ExportMixDown:    true,
		},
		ProjectMixer: SVPProjectMixer{
			LinkRoomSettings: true,
		},
	}

	out, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %s (%d parts, %d meters, %d tempos)\n", outputPath, len(tracks), len(meters), len(tempos))
	for _, t := range tracks {
		for _, g := range library {
			if g.UUID == t.Groups[0].GroupID {
				fmt.Fprintf(os.Stderr, "  %s: %d notes\n", t.Name, len(g.Notes))
			}
		}
	}
}
