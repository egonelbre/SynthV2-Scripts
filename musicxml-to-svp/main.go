package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/egonelbre/synthv2-scripts/internal/voice"
	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

const blicksPerQuarter = 705600000

// SVP output structs

type SVPProject struct {
	Version      int              `json:"version"`
	Time         SVPTime          `json:"time"`
	Library      []*SVPGroup      `json:"library"`
	Tracks       []*SVPTrack      `json:"tracks"`
	RenderConfig SVPRenderConfig  `json:"renderConfig"`
	ProjectMixer SVPProjectMixer  `json:"projectMixer"`
	UUID         string           `json:"uuid"`
}

type SVPTime struct {
	Meters []*SVPMeter `json:"meter"`
	Tempos []*SVPTempo `json:"tempo"`
}

type SVPMeter struct {
	Index       int `json:"index"`
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}

type SVPTempo struct {
	Position int64   `json:"position"`
	BPM      float64 `json:"bpm"`
}

type SVPGroup struct {
	Name       string        `json:"name"`
	UUID       string        `json:"uuid"`
	Notes      []*SVPNote    `json:"notes"`
	Parameters SVPParameters `json:"parameters"`
}

type SVPNote struct {
	Onset    int64    `json:"onset"`
	Duration int64    `json:"duration"`
	Lyrics   string   `json:"lyrics"`
	Phonemes string   `json:"phonemes"`
	Pitch    int      `json:"pitch"`
	Detune   int      `json:"detune"`
	Takes    SVPTakes `json:"attributes"`
}

type SVPTakes struct {
	EvenSyllableDuration bool         `json:"evenSyllableDuration"`
	Muted                bool         `json:"muted"`
	TF0Offset            float64      `json:"tF0Offset"`
	SystemPitchDelta     SVPParamMode `json:"systemPitchDelta"`
	Takes                []SVPTake    `json:"dur,omitempty"`
}

type SVPParamMode struct {
	Mode string `json:"mode"`
}

type SVPTake struct {
	ID    int       `json:"id"`
	Liked bool      `json:"liked"`
	Seeds SVPSeeds  `json:"seeds"`
}

type SVPSeeds struct {
	SingingSeed  int `json:"singingSeed"`
	BackingSeed  int `json:"backingSeed"`
}

type SVPParameters struct {
	PitchDelta   SVPParamCurve `json:"pitchDelta"`
	VibratoEnv   SVPParamCurve `json:"vibratoEnv"`
	Loudness     SVPParamCurve `json:"loudness"`
	Tension      SVPParamCurve `json:"tension"`
	Breathiness  SVPParamCurve `json:"breathiness"`
	Voicing      SVPParamCurve `json:"voicing"`
	Gender       SVPParamCurve `json:"gender"`
	ToneShift    SVPParamCurve `json:"toneShift"`
	MouthOpening SVPParamCurve `json:"mouthOpening"`
}

type SVPParamCurve struct {
	Mode   string    `json:"mode"`
	Points []float64 `json:"points"`
}

type SVPTrack struct {
	Name      string       `json:"name"`
	DispColor string       `json:"dispColor"`
	DispOrder int          `json:"dispOrder"`
	Mixer     SVPMixer     `json:"mixer"`
	MainGroup SVPGroupRef  `json:"mainGroup"`
	MainRef   SVPGroupRef  `json:"mainRef"`
	Groups    []SVPGroupRef `json:"groups"`
	UUID      string       `json:"uuid"`
}

type SVPGroupRef struct {
	GroupID        string       `json:"groupID"`
	BlickOffset    int64        `json:"blickOffset"`
	PitchOffset    int          `json:"pitchOffset"`
	IsInstrumental bool         `json:"isInstrumental"`
	Database       *SVPDatabase `json:"database,omitempty"`
	Voice          *SVPVoice    `json:"voice,omitempty"`
	UUID           string       `json:"uuid"`
}

type SVPDatabase struct {
	Name             string `json:"name"`
	Language         string `json:"language"`
	Phoneset         string `json:"phoneset"`
	LanguageOverride string `json:"languageOverride"`
	PhonesetOverride string `json:"phonesetOverride"`
	BackendType      string `json:"backendType"`
	Version          string `json:"version"`
}

type SVPVoice struct {
	VocalModeInherited     bool              `json:"vocalModeInherited"`
	VocalModePreset        string            `json:"vocalModePreset"`
	VocalModeParams        map[string]float64 `json:"vocalModeParams"`
	ChoirSeatingSeparation float64           `json:"choirSeatingSeparation,omitempty"`
	ChoirNumStems          int               `json:"choirNumStems,omitempty"`
	ChoirPartName          string            `json:"choirPartName,omitempty"`
}

type SVPMixer struct {
	GainDecibel    float64 `json:"gainDecibel"`
	Pan            float64 `json:"pan"`
	Mute           bool    `json:"mute"`
	Solo           bool    `json:"solo"`
	Display        bool    `json:"display"`
}

type SVPRenderConfig struct {
	Destination      string `json:"destination"`
	Filename         string `json:"filename"`
	NumChannels      int    `json:"numChannels"`
	AspirationFormat string `json:"aspirationFormat"`
	BitDepth         int    `json:"bitDepth"`
	SampleRate       int    `json:"sampleRate"`
	ExportMixDown    bool   `json:"exportMixDown"`
}

type SVPProjectMixer struct {
	GainDecibel      float64 `json:"gainDecibel"`
	Pan              float64 `json:"pan"`
	Mute             bool    `json:"mute"`
	Solo             bool    `json:"solo"`
	LinkRoomSettings bool    `json:"linkRoomSettings"`
}

// Helper functions

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func newShortID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func newEmptyParamCurve() SVPParamCurve {
	return SVPParamCurve{Mode: "cubic", Points: []float64{}}
}

func newEmptyParameters() SVPParameters {
	return SVPParameters{
		PitchDelta:  newEmptyParamCurve(),
		VibratoEnv:  newEmptyParamCurve(),
		Loudness:    newEmptyParamCurve(),
		Tension:     newEmptyParamCurve(),
		Breathiness: newEmptyParamCurve(),
		Voicing:     newEmptyParamCurve(),
		Gender:       newEmptyParamCurve(),
		ToneShift:    newEmptyParamCurve(),
		MouthOpening: newEmptyParamCurve(),
	}
}

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

// dynamicsToLoudness maps a MusicXML dynamics element to a dB loudness value.
// It uses InnerXML to detect which child elements are present, since empty
// elements like <p/> parse to empty strings and can't be detected via field values.
func dynamicsToLoudness(d *musicxml.Dynamics) (float64, bool) {
	xml := d.InnerXML

	// Check from most specific to least specific to avoid prefix matching issues.
	// E.g. check <pppppp before <ppppp before <pppp etc.
	switch {
	case strings.Contains(xml, "<ffffff"):
		return 15, true
	case strings.Contains(xml, "<fffff"):
		return 12, true
	case strings.Contains(xml, "<ffff"):
		return 9, true
	case strings.Contains(xml, "<fff"):
		return 6, true
	case strings.Contains(xml, "<ff"):
		return 3, true

	case strings.Contains(xml, "<pppppp"):
		return -36, true
	case strings.Contains(xml, "<ppppp"):
		return -33, true
	case strings.Contains(xml, "<pppp"):
		return -30, true
	case strings.Contains(xml, "<ppp"):
		return -27, true
	case strings.Contains(xml, "<pp"):
		return -24, true

	case strings.Contains(xml, "<sffz"):
		return 3, true
	case strings.Contains(xml, "<sfzp"):
		return 0, true
	case strings.Contains(xml, "<sfpp"):
		return 0, true
	case strings.Contains(xml, "<sfz"):
		return 3, true
	case strings.Contains(xml, "<sfp"):
		return 0, true
	case strings.Contains(xml, "<sf"):
		return 3, true

	case strings.Contains(xml, "<mp"):
		return -12, true
	case strings.Contains(xml, "<mf"):
		return -6, true

	case strings.Contains(xml, "<fp"):
		return 0, true
	case strings.Contains(xml, "<fz"):
		return 3, true
	case strings.Contains(xml, "<f"):
		return 0, true

	case strings.Contains(xml, "<rfz"):
		return 0, true
	case strings.Contains(xml, "<rf"):
		return 0, true

	case strings.Contains(xml, "<pf"):
		return -6, true
	case strings.Contains(xml, "<p"):
		return -18, true

	case strings.Contains(xml, "<n"):
		return -48, true
	}
	return 0, false
}

func extractLyric(note *musicxml.Note) string {
	if len(note.Lyric) == 0 {
		return ""
	}
	lyric := note.Lyric[0]
	text := ""
	if len(lyric.Text) > 0 {
		text = lyric.Text[0].EnclosedText
	}
	if text == "" {
		return ""
	}
	switch lyric.Syllabic {
	case "begin", "middle":
		return text + " -"
	default:
		return text
	}
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

// Dynamic event types for building loudness curves.
type dynEventKind int

const (
	dynLevel      dynEventKind = iota // instant dynamic level (p, mf, f, ...)
	dynCrescStart                     // start of crescendo hairpin or text cresc.
	dynDimStart                       // start of diminuendo hairpin or text dim.
	dynWedgeStop                      // end of hairpin
)

type dynEvent struct {
	position int64
	kind     dynEventKind
	loudness float64 // only meaningful for dynLevel
	number   int     // wedge number for matching start/stop
}

// buildLoudnessCurve takes a chronological list of dynamic events and produces
// SVP loudness curve points [pos, val, pos, val, ...].
//
// Strategy:
//   - At each dynLevel event, emit a quick step transition (two points close
//     together) from the previous level to the new level.
//   - Between a cresc/dim start and the following wedge stop, emit a linear
//     ramp. The ramp target is the loudness of the next dynLevel event after
//     the stop (if one exists at the same position), otherwise we estimate
//     +6 dB for crescendo, -6 dB for diminuendo.
//   - Between other events, hold the current level (no points needed, cubic
//     interpolation holds the last value).
const stepTransitionBlicks = int64(blicksPerQuarter / 8) // 1/8 quarter note

func buildLoudnessCurve(events []dynEvent) []float64 {
	if len(events) == 0 {
		return nil
	}

	var points []float64
	currentLevel := 0.0
	hasLevel := false

	// Pair up wedge starts and stops by number.
	type wedgeInfo struct {
		startIdx int
		stopIdx  int
		kind     dynEventKind // dynCrescStart or dynDimStart
	}
	openWedges := map[int]int{} // number -> index in events
	var wedges []wedgeInfo

	for i, ev := range events {
		switch ev.kind {
		case dynCrescStart, dynDimStart:
			openWedges[ev.number] = i
		case dynWedgeStop:
			if startIdx, ok := openWedges[ev.number]; ok {
				wedges = append(wedges, wedgeInfo{
					startIdx: startIdx,
					stopIdx:  i,
					kind:     events[startIdx].kind,
				})
				delete(openWedges, ev.number)
			}
		}
	}

	// Build a set of wedge ranges for quick lookup.
	type wedgeRange struct {
		startPos  int64
		stopPos   int64
		kind      dynEventKind
		stopIdx   int
	}
	var ranges []wedgeRange
	for _, w := range wedges {
		ranges = append(ranges, wedgeRange{
			startPos: events[w.startIdx].position,
			stopPos:  events[w.stopIdx].position,
			kind:     w.kind,
			stopIdx:  w.stopIdx,
		})
	}

	// Find the next dynLevel event at or after a given index.
	findNextLevel := func(fromIdx int) (float64, bool) {
		for j := fromIdx; j < len(events); j++ {
			if events[j].kind == dynLevel {
				return events[j].loudness, true
			}
		}
		return 0, false
	}

	addPoint := func(pos int64, val float64) {
		points = append(points, float64(pos), val)
	}

	for i, ev := range events {
		switch ev.kind {
		case dynLevel:
			// Check if this level is the target of a just-ended wedge.
			// If so, the ramp already brought us here; just update currentLevel.
			isWedgeTarget := false
			for _, wr := range ranges {
				if wr.stopPos == ev.position || (ev.position-wr.stopPos >= 0 && ev.position-wr.stopPos < stepTransitionBlicks*2) {
					if wr.stopIdx == i-1 || (i > 0 && events[i-1].kind == dynWedgeStop) {
						isWedgeTarget = true
					}
				}
			}

			if isWedgeTarget {
				// End of ramp — just place the final point.
				addPoint(ev.position, ev.loudness)
			} else if hasLevel {
				// Step transition: hold old level, then jump to new.
				transitionStart := ev.position - stepTransitionBlicks
				if transitionStart < 0 {
					transitionStart = 0
				}
				addPoint(transitionStart, currentLevel)
				addPoint(ev.position, ev.loudness)
			} else {
				// First dynamic marking — set initial level.
				addPoint(ev.position, ev.loudness)
			}

			currentLevel = ev.loudness
			hasLevel = true

		case dynCrescStart, dynDimStart:
			// Find the matching stop and determine target level.
			var stopPos int64
			var targetLevel float64
			found := false

			for _, wr := range ranges {
				if wr.startPos == ev.position && wr.kind == ev.kind {
					stopPos = wr.stopPos
					// Look for a dynLevel right after the stop.
					if lvl, ok := findNextLevel(wr.stopIdx + 1); ok {
						targetLevel = lvl
					} else {
						// Estimate.
						if ev.kind == dynCrescStart {
							targetLevel = currentLevel + 6
						} else {
							targetLevel = currentLevel - 6
						}
					}
					found = true
					break
				}
			}

			if !found {
				// Unpaired cresc/dim text — estimate over 2 measures.
				stopPos = ev.position + 2*4*blicksPerQuarter
				if ev.kind == dynCrescStart {
					targetLevel = currentLevel + 6
				} else {
					targetLevel = currentLevel - 6
				}
			}

			// Emit ramp: start point at current level, end point at target.
			if hasLevel {
				addPoint(ev.position, currentLevel)
			}
			addPoint(stopPos, targetLevel)
			currentLevel = targetLevel

		case dynWedgeStop:
			// Already handled via ranges above.
		}
	}

	return points
}

// isTextCresc checks if a words element indicates crescendo.
func isTextCresc(text string) bool {
	t := strings.ToLower(strings.TrimSpace(text))
	return strings.HasPrefix(t, "cresc") || strings.HasPrefix(t, "crésc")
}

// isTextDim checks if a words element indicates diminuendo/decrescendo.
func isTextDim(text string) bool {
	t := strings.ToLower(strings.TrimSpace(text))
	return strings.HasPrefix(t, "dim") || strings.HasPrefix(t, "decresc")
}

func main() {
	voiceFlag := flag.String("voice", "", "assign voices: choir1, choir2, choir3, or soloists")
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
		cursor := int64(0)
		var notes []*SVPNote
		var dynEvents []dynEvent
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
				case *musicxml.Direction:
					for _, dt := range value.DirectionType {
						// Dynamics markings (p, mf, f, ...)
						for _, dyn := range dt.Dynamics {
							if loudness, ok := dynamicsToLoudness(dyn); ok {
								dynEvents = append(dynEvents, dynEvent{
									position: cursor,
									kind:     dynLevel,
									loudness: loudness,
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

		// Build loudness curve from collected dynamic events
		params := newEmptyParameters()
		if len(dynEvents) > 0 {
			params.Loudness.Points = buildLoudnessCurve(dynEvents)
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
		assignVoices(tracks, strings.ToLower(*voiceFlag))
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

func assignVoices(tracks []*SVPTrack, voiceArg string) {
	// Parse track names into voice parts.
	infos := make([]voice.TrackInfo, len(tracks))
	for i, t := range tracks {
		infos[i] = voice.TrackInfo{
			Name:    t.Name,
			Part:    voice.ParseVoicePart(t.Name),
			PartNum: voice.ParsePartNum(t.Name),
		}
	}

	switch voiceArg {
	case "choir1", "1":
		applyChoirToTracks(tracks, infos, voice.Choirs[0])
	case "choir2", "2":
		applyChoirToTracks(tracks, infos, voice.Choirs[1])
	case "choir3", "3":
		applyChoirToTracks(tracks, infos, voice.Choirs[2])
	case "soloists", "solo", "4":
		applySoloistsToTracks(tracks, infos)
	default:
		fmt.Fprintf(os.Stderr, "unknown voice source: %q (options: choir1, choir2, choir3, soloists)\n", voiceArg)
		os.Exit(1)
	}
}

func applyChoirToTracks(tracks []*SVPTrack, infos []voice.TrackInfo, choir voice.ChoirInfo) {
	db := &SVPDatabase{
		Name:     choir.Name,
		Language: choir.Language,
		Phoneset: choir.Phoneset,
		BackendType: "SVR3",
		Version:     "202",
	}

	for i, t := range tracks {
		info := infos[i]
		if info.Part == voice.Unknown {
			fmt.Fprintf(os.Stderr, "  skipping %q (unknown voice part)\n", info.Name)
			continue
		}

		choirPart := voice.MapPartToChoir(info.Part, choir)
		fmt.Fprintf(os.Stderr, "  %s -> %s (%s)\n", info.Name, choir.Name, choirPart.String())

		// Set database and voice on MainRef.
		t.MainRef.Database = db
		t.MainRef.Voice = &SVPVoice{
			VocalModeInherited:     true,
			VocalModeParams:        map[string]float64{},
			ChoirSeatingSeparation: 0.7,
		}

		// Set choir part on group refs.
		for j := range t.Groups {
			partName := ""
			if choirPart != voice.Soprano {
				partName = string(choirPart)
			}
			t.Groups[j].Voice = &SVPVoice{
				VocalModeInherited:     true,
				VocalModeParams:        map[string]float64{},
				ChoirNumStems:          4,
				ChoirSeatingSeparation: 0.7,
				ChoirPartName:          partName,
			}
		}
	}
}

func applySoloistsToTracks(tracks []*SVPTrack, infos []voice.TrackInfo) {
	assignments := voice.AssignSoloists(infos)

	for i, t := range tracks {
		info := infos[i]
		soloist, ok := assignments[i]
		if !ok {
			if info.Part == voice.Unknown {
				fmt.Fprintf(os.Stderr, "  skipping %q (unknown voice part)\n", info.Name)
			}
			continue
		}

		fmt.Fprintf(os.Stderr, "  %s -> %s\n", info.Name, soloist.Name)

		t.MainRef.Database = &SVPDatabase{
			Name:        soloist.Name,
			Language:    soloist.Language,
			Phoneset:    soloist.Phoneset,
			BackendType: "SVR3",
			Version:     "201",
		}
		t.MainRef.Voice = &SVPVoice{
			VocalModeInherited: true,
			VocalModeParams:    map[string]float64{},
		}
	}
}

func parseDuration(s string) int {
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}
