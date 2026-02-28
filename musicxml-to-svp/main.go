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

func extractLyric(note *musicxml.Note, verse int) string {
	if verse > 0 {
		// Look for exact match first, then fall back to highest lyric number ≤ verse.
		bestNum := 0
		bestText := ""
		for _, lyric := range note.Lyric {
			if len(lyric.Text) == 0 || lyric.Text[0].EnclosedText == "" {
				continue
			}
			num, err := strconv.Atoi(lyric.Number)
			if err != nil {
				continue
			}
			if num == verse {
				return lyric.Text[0].EnclosedText
			}
			if num <= verse && num > bestNum {
				bestNum = num
				bestText = lyric.Text[0].EnclosedText
			}
		}
		if bestText != "" {
			return bestText
		}
	}
	// Fallback: first non-empty lyric.
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

type playedMeasure struct {
	measureIdx int // index into part.Measure
	verse      int // 1-based lyric number for this pass
}

func measureHasForwardRepeat(m *musicxml.Measure) bool {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Repeat != nil && bl.Repeat.Direction == "forward" {
				return true
			}
		}
	}
	return false
}

func measureBackwardRepeat(m *musicxml.Measure) *musicxml.Repeat {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Repeat != nil && bl.Repeat.Direction == "backward" {
				return bl.Repeat
			}
		}
	}
	return nil
}

func measureEndingStart(m *musicxml.Measure) *musicxml.Ending {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Ending != nil && bl.Ending.Type == "start" {
				return bl.Ending
			}
		}
	}
	return nil
}

func measureHasEndingStop(m *musicxml.Measure) bool {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Ending != nil && (bl.Ending.Type == "stop" || bl.Ending.Type == "discontinue") {
				return true
			}
		}
	}
	return false
}

func parseEndingNumbers(s string) []int {
	var nums []int
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if n, err := strconv.Atoi(part); err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

func intSliceContains(s []int, v int) bool {
	for _, n := range s {
		if n == v {
			return true
		}
	}
	return false
}

func unrollMeasures(measures []*musicxml.Measure) []playedMeasure {
	// Pre-scan to find repeat regions.
	type repeatRegion struct {
		start, end int
		times      int
	}
	var regions []repeatRegion
	repeatStart := 0
	for i, m := range measures {
		if measureHasForwardRepeat(m) {
			repeatStart = i
		}
		if br := measureBackwardRepeat(m); br != nil {
			times := br.Times
			if times == 0 {
				times = 2
			}
			regions = append(regions, repeatRegion{start: repeatStart, end: i, times: times})
			repeatStart = i + 1
		}
	}

	// Emit the unrolled sequence.
	var result []playedMeasure
	i := 0
	regionIdx := 0

	for i < len(measures) {
		if regionIdx < len(regions) && i == regions[regionIdx].start {
			r := regions[regionIdx]
			regionIdx++

			// Find endings and track which measures are inside which ending.
			type endingInfo struct {
				inEnding bool
				numbers  []int
			}
			endingInfos := make([]endingInfo, r.end-r.start+1)
			var currentEndingNums []int
			inEnding := false
			maxEndingNum := 0
			hasEndings := false

			for j := r.start; j <= r.end; j++ {
				if startEnding := measureEndingStart(measures[j]); startEnding != nil {
					inEnding = true
					currentEndingNums = parseEndingNumbers(startEnding.Number)
					hasEndings = true
					for _, n := range currentEndingNums {
						if n > maxEndingNum {
							maxEndingNum = n
						}
					}
				}
				if inEnding {
					endingInfos[j-r.start] = endingInfo{inEnding: true, numbers: currentEndingNums}
				}
				if measureHasEndingStop(measures[j]) {
					inEnding = false
				}
			}

			// Emit passes.
			for pass := 1; pass <= r.times; pass++ {
				for j := r.start; j <= r.end; j++ {
					ei := endingInfos[j-r.start]
					if ei.inEnding && !intSliceContains(ei.numbers, pass) {
						continue
					}
					result = append(result, playedMeasure{measureIdx: j, verse: pass})
				}
			}

			i = r.end + 1

			// If there are endings and the final pass exceeds all ending numbers,
			// consume continuation measures (implicit final volta) with the final verse.
			if hasEndings && r.times > maxEndingNum {
				finalVerse := r.times
				for i < len(measures) {
					if regionIdx < len(regions) && i == regions[regionIdx].start {
						break
					}
					result = append(result, playedMeasure{measureIdx: i, verse: finalVerse})
					i++
				}
			}
			continue
		}

		result = append(result, playedMeasure{measureIdx: i, verse: 1})
		i++
	}

	return result
}

type measureInfo struct {
	startBlicks int64
	divisions   int
}

// buildStructure unrolls repeats and computes measure start positions, meters, and tempos
// from the first part.
func buildStructure(firstPart *musicxml.Part) ([]playedMeasure, []measureInfo, []MeterChange, []TempoChange) {
	unrolled := unrollMeasures(firstPart.Measure)

	var meters []MeterChange
	var tempos []TempoChange
	var infos []measureInfo

	divisions := 4
	cursor := int64(0)

	for measureIdx, pm := range unrolled {
		measure := firstPart.Measure[pm.measureIdx]

		infos = append(infos, measureInfo{
			startBlicks: cursor,
			divisions:   divisions,
		})

		measureDuration := int64(0)
		for _, el := range measure.Element {
			switch value := el.Value.(type) {
			case *musicxml.Attributes:
				if value.Divisions != 0 {
					divisions = value.Divisions
					infos[len(infos)-1].divisions = divisions
				}
				for _, t := range value.Time {
					meters = append(meters, MeterChange{
						MeasureIndex: measureIdx,
						Numerator:    t.Beats,
						Denominator:  t.BeatType,
					})
				}
			case *musicxml.Direction:
				if value.Sound != nil && value.Sound.Tempo != "" {
					bpm, err := strconv.ParseFloat(value.Sound.Tempo, 64)
					if err == nil {
						tempos = append(tempos, TempoChange{
							Position: cursor + measureDuration,
							BPM:      bpm,
						})
					}
				}
				if len(tempos) == 0 || value.Sound == nil {
					for _, dt := range value.DirectionType {
						if dt.Metronome != nil && dt.Metronome.PerMinute != nil {
							bpm, err := strconv.ParseFloat(dt.Metronome.PerMinute.EnclosedText, 64)
							if err == nil {
								q := beatUnitToQuarters(dt.Metronome.BeatUnit, dt.Metronome.BeatUnitDot != "")
								bpm = bpm * q
								tempos = append(tempos, TempoChange{
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

		meterNum := 4
		meterDen := 4
		for _, m := range meters {
			if m.MeasureIndex <= measureIdx {
				meterNum = m.Numerator
				meterDen = m.Denominator
			}
		}
		expectedDuration := int64(meterNum) * blicksPerQuarter * 4 / int64(meterDen)
		cursor += expectedDuration
	}

	return unrolled, infos, meters, tempos
}

func noteArticulations(note *musicxml.Note) Articulation {
	var a Articulation
	if noteHasArticulation(note, "staccato") {
		a |= ArticulationStaccato
	}
	if noteHasArticulation(note, "staccatissimo") {
		a |= ArticulationStaccatissimo
	}
	if noteHasArticulation(note, "tenuto") {
		a |= ArticulationTenuto
	}
	if noteHasArticulation(note, "accent") {
		a |= ArticulationAccent
	}
	if noteHasArticulation(note, "strong-accent") {
		a |= ArticulationStrongAccent
	}
	return a
}

func makeGraceNote(g *musicxml.Note, transpose, verse int) GraceNote {
	midi, detune := 0, 0
	if g.Pitch != nil {
		midi, detune = pitchToMIDI(g.Pitch)
		midi += transpose
	}
	notatedType := ""
	if g.Type != nil {
		notatedType = string(g.Type.EnclosedText)
	}
	return GraceNote{
		Pitch:        midi,
		Detune:       detune,
		Lyric:        extractLyric(g, verse),
		NotatedType:  notatedType,
		Acciaccatura: g.Grace.Slash == "yes",
	}
}

// flushTailGraces attaches pending grace notes as trailing graces to the last note,
// computing durations and stealing time from the parent note.
func flushTailGraces(notes []Note, pendingGraces []*musicxml.Note, transpose, verse int) []Note {
	if len(notes) == 0 || len(pendingGraces) == 0 {
		return notes
	}
	last := &notes[len(notes)-1]
	var graces []GraceNote
	for _, g := range pendingGraces {
		if g.Pitch == nil {
			continue
		}
		graces = append(graces, makeGraceNote(g, transpose, verse))
	}
	if len(graces) == 0 {
		return notes
	}
	graceDurs, totalGrace := capGraceDurs(graces, last.Duration/2)
	for i := range graces {
		graces[i].Duration = graceDurs[i]
	}
	last.Duration -= totalGrace
	last.TrailingGraces = append(last.TrailingGraces, graces...)
	return notes
}

// buildNotes extracts notes from a part, resolving ties and attaching grace notes.
func buildNotes(part *musicxml.Part, unrolled []playedMeasure, infos []measureInfo) []Note {
	divisions := 4
	transpose := 0
	cursor := int64(0)
	var notes []Note
	pendingTies := map[int]int{} // MIDI pitch -> index in notes
	var prevOnset int64
	var pendingGraces []*musicxml.Note
	var graceIsTail bool
	var lastNoteIdx int = -1
	var lastVerse int

	for unrolledIdx, pm := range unrolled {
		lastVerse = pm.verse
		measure := part.Measure[pm.measureIdx]

		if unrolledIdx < len(infos) {
			cursor = infos[unrolledIdx].startBlicks
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
				// Skip directions in buildNotes; handled by buildDynamics.
			case *musicxml.Note:
				if value.Grace != nil {
					if len(pendingGraces) == 0 {
						graceIsTail = lastNoteIdx >= 0
					}
					pendingGraces = append(pendingGraces, value)
					continue
				}

				dur := parseDuration(value.Duration)
				if dur <= 0 {
					continue
				}
				blicks := durationToBlicks(dur, divisions)

				onset := cursor
				if value.Chord != "" {
					onset = prevOnset
				}

				if value.Rest != nil || value.Pitch == nil {
					if graceIsTail && len(pendingGraces) > 0 && lastNoteIdx >= 0 {
						notes = flushTailGraces(notes, pendingGraces, transpose, pm.verse)
					}
					pendingGraces = nil
					lastNoteIdx = -1
					if value.Chord == "" {
						prevOnset = cursor
						cursor += blicks
					}
					continue
				}

				// Process pending grace notes.
				var leadingGraces []GraceNote
				if len(pendingGraces) > 0 {
					if graceIsTail && lastNoteIdx >= 0 {
						notes = flushTailGraces(notes, pendingGraces, transpose, pm.verse)
					} else {
						for _, g := range pendingGraces {
							if g.Pitch == nil {
								continue
							}
							leadingGraces = append(leadingGraces, makeGraceNote(g, transpose, pm.verse))
						}
						// Compute durations and steal time from following note.
						if len(leadingGraces) > 0 {
							graceDurs, totalGrace := capGraceDurs(leadingGraces, blicks/2)
							for i := range leadingGraces {
								leadingGraces[i].Duration = graceDurs[i]
							}
							onset += totalGrace
							blicks -= totalGrace
						}
					}
					pendingGraces = nil
				}

				midi, detune := pitchToMIDI(value.Pitch)
				midi += transpose
				tieStart, tieStop := noteTieTypes(value)

				if tieStop {
					if idx, ok := pendingTies[midi]; ok {
						notes[idx].Duration += blicks
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

				lyric := extractLyric(value, pm.verse)

				note := Note{
					Onset:         onset,
					Duration:      blicks,
					Pitch:         midi,
					Detune:        detune,
					Lyric:         lyric,
					Articulations: noteArticulations(value),
					LeadingGraces: leadingGraces,
				}

				notes = append(notes, note)
				lastNoteIdx = len(notes) - 1

				if tieStart {
					pendingTies[midi] = lastNoteIdx
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
	}

	// Flush any remaining tail graces.
	if graceIsTail && len(pendingGraces) > 0 && lastNoteIdx >= 0 {
		notes = flushTailGraces(notes, pendingGraces, transpose, lastVerse)
	}

	return notes
}

// fillLyrics handles grace note lyric assignment and melismatic continuations.
func fillLyrics(notes []Note) {
	for i := range notes {
		n := &notes[i]

		// If no leading grace has a lyric, copy the main note's lyric to the first.
		if len(n.LeadingGraces) > 0 {
			anyGraceLyric := false
			for _, g := range n.LeadingGraces {
				if g.Lyric != "" {
					anyGraceLyric = true
					break
				}
			}
			if !anyGraceLyric && n.Lyric != "" {
				n.LeadingGraces[0].Lyric = n.Lyric
			}
		}

		// Fill empty lyrics with melismatic continuation.
		for j := range n.LeadingGraces {
			if n.LeadingGraces[j].Lyric == "" {
				n.LeadingGraces[j].Lyric = "-"
			}
		}
		for j := range n.TrailingGraces {
			if n.TrailingGraces[j].Lyric == "" {
				n.TrailingGraces[j].Lyric = "-"
			}
		}
		if n.Lyric == "" {
			n.Lyric = "-"
		}
	}
}

// buildDynamics collects dynamic events from directions in a part.
func buildDynamics(part *musicxml.Part, unrolled []playedMeasure, infos []measureInfo) []dynEvent {
	var events []dynEvent
	cursor := int64(0)

	for unrolledIdx, pm := range unrolled {
		measure := part.Measure[pm.measureIdx]

		if unrolledIdx < len(infos) {
			cursor = infos[unrolledIdx].startBlicks
		}

		divisions := 4
		if unrolledIdx < len(infos) {
			divisions = infos[unrolledIdx].divisions
		}

		for _, el := range measure.Element {
			switch value := el.Value.(type) {
			case *musicxml.Attributes:
				if value.Divisions != 0 {
					divisions = value.Divisions
				}
			case *musicxml.Direction:
				for _, dt := range value.DirectionType {
					for _, dyn := range dt.Dynamics {
						if lvl, ok := dynamicsToLevel(dyn); ok {
							events = append(events, dynEvent{
								position: cursor,
								kind:     dynLevel,
								loudness: lvl.loudness,
								tension:  lvl.tension,
							})
						}
					}
					if dt.Wedge != nil {
						num := dt.Wedge.Number
						if num == 0 {
							num = 1
						}
						switch dt.Wedge.Type {
						case "crescendo":
							events = append(events, dynEvent{
								position: cursor,
								kind:     dynCrescStart,
								number:   num,
							})
						case "diminuendo":
							events = append(events, dynEvent{
								position: cursor,
								kind:     dynDimStart,
								number:   num,
							})
						case "stop":
							events = append(events, dynEvent{
								position: cursor,
								kind:     dynWedgeStop,
								number:   num,
							})
						}
					}
					for _, w := range dt.Words {
						if isTextCresc(w.EnclosedText) {
							events = append(events, dynEvent{
								position: cursor,
								kind:     dynCrescStart,
								number:   -1,
							})
						} else if isTextDim(w.EnclosedText) {
							events = append(events, dynEvent{
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
				if dur > 0 && value.Chord == "" {
					cursor += durationToBlicks(dur, divisions)
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
	}

	return events
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

	// Pass 1: structure
	irScore := &Score{}
	var unrolled []playedMeasure
	var infos []measureInfo
	if len(score.Part) > 0 {
		unrolled, infos, irScore.Meters, irScore.Tempos = buildStructure(score.Part[0])
	}

	// Pass 2+3+4: per part
	for partIdx, part := range score.Part {
		partName := partNames[part.Id]
		if partName == "" {
			partName = fmt.Sprintf("Part %d", partIdx+1)
		}
		p := Part{Name: partName}
		p.Notes = buildNotes(part, unrolled, infos)
		fillLyrics(p.Notes)
		p.Dynamics = buildDynamics(part, unrolled, infos)
		irScore.Parts = append(irScore.Parts, p)
	}

	// Convert to SVP
	project := scoreToSVP(irScore)

	// Assign voices if requested.
	if *voiceFlag != "" {
		assignVoices(project.Tracks, strings.ToLower(*voiceFlag), *relaxedFlag, *panFlag)
		setNoteAttributes(project.Library)
	}

	// Convert lyrics to phonemes if requested.
	if *langFlag != "" {
		conv := phonemes.New(*langFlag)
		if conv == nil {
			fmt.Fprintf(os.Stderr, "unknown language: %q (options: estonian, karelian)\n", *langFlag)
			os.Exit(1)
		}
		applyPhonemes(project.Library, conv)
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

	fmt.Fprintf(os.Stderr, "wrote %s (%d parts, %d meters, %d tempos)\n", outputPath, len(project.Tracks), len(project.Time.Meters), len(project.Time.Tempos))
	for _, t := range project.Tracks {
		for _, g := range project.Library {
			if g.UUID == t.Groups[0].GroupID {
				fmt.Fprintf(os.Stderr, "  %s: %d notes\n", t.Name, len(g.Notes))
			}
		}
	}
}
