package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

const blicksPerQuarter = 705600000

var stepToSemitone = map[musicxml.Step]int{
	"C": 0, "D": 2, "E": 4, "F": 5, "G": 7, "A": 9, "B": 11,
}

func pitchToMIDI(p *musicxml.Pitch) (midi int, detune int) {
	semitone, ok := stepToSemitone[p.Step]
	if !ok {
		fmt.Fprintf(os.Stderr, "warning: unknown pitch step %q, defaulting to C\n", p.Step)
	}
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

// noteTypeToQuarters returns the duration in quarter notes for a note type name.
func noteTypeToQuarters(noteType string) float64 {
	switch noteType {
	case "whole":
		return 4.0
	case "half":
		return 2.0
	case "quarter":
		return 1.0
	case "eighth":
		return 0.5
	case "16th":
		return 0.25
	case "32nd":
		return 0.125
	default:
		return 1.0
	}
}

func beatUnitToQuarters(beatUnit musicxml.NoteTypeValue, hasDot bool) float64 {
	q := noteTypeToQuarters(string(beatUnit))
	if hasDot {
		q *= 1.5
	}
	return q
}

func parseBeats(s string) int {
	total := 0
	for _, part := range strings.Split(s, "+") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: invalid beats value %q: %v\n", part, err)
			continue
		}
		total += v
	}
	return total
}

func parseDuration(s string) int {
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: invalid duration %q: %v\n", s, err)
		return 0
	}
	return v
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
		return bestText
	}
	// No verse filtering: return first non-empty lyric.
	for _, lyric := range note.Lyric {
		if len(lyric.Text) > 0 && lyric.Text[0].EnclosedText != "" {
			return lyric.Text[0].EnclosedText
		}
	}
	return ""
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

// validateTimeModification checks that a note's duration is consistent with its
// time-modification element. Logs a warning to stderr on mismatch but does not
// fail, since the duration field takes precedence.
func validateTimeModification(note *musicxml.Note, divisions int) {
	if note.TimeModification == nil {
		return
	}
	tm := note.TimeModification
	if tm.ActualNotes <= 0 || tm.NormalNotes <= 0 {
		return
	}
	if note.Type == nil {
		return
	}

	normalQuarters := beatUnitToQuarters(musicxml.NoteTypeValue(note.Type.EnclosedText), len(note.Dot) > 0)
	expectedDuration := normalQuarters * float64(divisions) * float64(tm.NormalNotes) / float64(tm.ActualNotes)

	dur := parseDuration(note.Duration)
	if dur <= 0 {
		return
	}

	diff := math.Abs(float64(dur) - expectedDuration)
	if diff > 0.5 {
		fmt.Fprintf(os.Stderr, "warning: time-modification mismatch: note type %q with %d/%d tuplet expects duration %.0f, got %d (divisions=%d)\n",
			note.Type.EnclosedText, tm.ActualNotes, tm.NormalNotes, expectedDuration, dur, divisions)
	}
}

// walkPartElements iterates through unrolled measures of a part, tracking
// cursor position and divisions. It calls fn for each element with the current
// cursor position and divisions. Cursor advancement for Note, Backup, and
// Forward elements is handled automatically after the callback.
func walkPartElements(
	part *musicxml.Part,
	unrolled []playedMeasure,
	fn func(cursor int64, divisions int, pm playedMeasure, value any),
) {
	divisions := 4
	cursor := int64(0)

	for _, pm := range unrolled {
		measure := part.Measure[pm.measureIdx]

		cursor = pm.startBlicks
		divisions = pm.divisions

		for _, el := range measure.Element {
			switch value := el.Value.(type) {
			case *musicxml.Attributes:
				if value.Divisions != 0 {
					divisions = value.Divisions
				}
				fn(cursor, divisions, pm, value)
			case *musicxml.Direction:
				fn(cursor, divisions, pm, value)
			case *musicxml.Note:
				fn(cursor, divisions, pm, value)
				if value.Grace == nil {
					dur := parseDuration(value.Duration)
					if dur > 0 && value.Chord == "" {
						cursor += durationToBlicks(dur, divisions)
					}
				}
			case *musicxml.Backup:
				fn(cursor, divisions, pm, value)
				dur := parseDuration(value.Duration)
				if dur > 0 {
					cursor -= durationToBlicks(dur, divisions)
				}
			case *musicxml.Forward:
				fn(cursor, divisions, pm, value)
				dur := parseDuration(value.Duration)
				if dur > 0 {
					cursor += durationToBlicks(dur, divisions)
				}
			}
		}
	}
}

