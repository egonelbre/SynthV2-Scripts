package main

import (
	"fmt"
	"slices"
)

// splitVoices splits parts with overlapping notes into one part per voice.
// Overlapping notes must belong to different voices, so a greedy assignment
// places each note on the first voice that is free at its onset.
func splitVoices(parts []Part) []Part {
	var result []Part
	for _, p := range parts {
		result = append(result, splitPart(p)...)
	}
	return result
}

func splitPart(p Part) []Part {
	notes := slices.Clone(p.Notes)
	slices.SortStableFunc(notes, func(a, b Note) int {
		if a.Onset != b.Onset {
			return int(min(max(a.Onset-b.Onset, -1), 1))
		}
		// Same onset: higher pitch first, so voice 1 is the top voice.
		return b.Pitch - a.Pitch
	})

	var voices [][]Note
	var ends []int64 // ends[v] = end of the last note assigned to voice v
	for _, n := range notes {
		v := slices.IndexFunc(ends, func(end int64) bool { return n.Onset >= end })
		if v < 0 {
			v = len(ends)
			ends = append(ends, 0)
			voices = append(voices, nil)
		}
		ends[v] = n.Onset + n.Duration
		voices[v] = append(voices[v], n)
	}

	if len(voices) <= 1 {
		return []Part{p}
	}

	out := make([]Part, len(voices))
	for i, notes := range voices {
		out[i] = Part{
			Name:  fmt.Sprintf("%s (Voice %d)", p.Name, i+1),
			Notes: notes,
			// ponytail: dynamics are duplicated to every voice, splitting them
			// per voice needs staff/voice info we don't carry in the IR.
			Dynamics: p.Dynamics,
		}
	}
	return out
}
