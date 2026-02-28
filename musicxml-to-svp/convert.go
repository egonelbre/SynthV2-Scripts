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

func parseBeats(s string) int {
	total := 0
	for _, part := range strings.Split(s, "+") {
		v, _ := strconv.Atoi(strings.TrimSpace(part))
		total += v
	}
	return total
}

func parseDuration(s string) int {
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
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
