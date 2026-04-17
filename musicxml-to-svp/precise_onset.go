package main

import "sort"

// preciseOnsetEpsilon is one 128th note. Each precision control is one epsilon
// wide and placed one epsilon away from the note edge:
//
//	[1/64 control][1/64 rest] | note onset ... note offset | [1/64 rest][1/64 control]
const preciseOnsetEpsilon = int64(blicksPerQuarter / 32)

// addPreciseOnsetControls appends pitch-control curves just before each phrase
// onset and just after each phrase offset so the SynthV engine doesn't slide
// into or out of the note. Two adjacent notes need at least 4 epsilons of rest
// between them to fit the prev-note offset control (rest + control) and the
// next-note onset control (control + rest); anything less is treated as a
// phrase continuation and both-side controls are skipped.
func addPreciseOnsetControls(library []*SVPGroup) {
	for _, g := range library {
		notes := append([]*SVPNote(nil), g.Notes...)
		sort.Slice(notes, func(i, j int) bool { return notes[i].Onset < notes[j].Onset })

		touching := func(a, b *SVPNote) bool {
			return a.Onset+a.Duration+4*preciseOnsetEpsilon > b.Onset
		}

		for i, n := range notes {
			var prev, next *SVPNote
			if i > 0 {
				prev = notes[i-1]
			}
			if i+1 < len(notes) {
				next = notes[i+1]
			}

			pitch := float64(n.Pitch) + float64(n.Detune)/100.0

			// Onset control: 1/64 wide, ending 1/64 before the note onset.
			if prev == nil || !touching(prev, n) {
				g.PitchControls = append(g.PitchControls, SVPPitchControl{
					Pos:    n.Onset,
					Pitch:  pitch,
					ID:     newShortID(),
					Type:   "curve",
					Points: []float64{-2 * float64(preciseOnsetEpsilon), 0, -float64(preciseOnsetEpsilon), 0},
				})
			}

			// Offset control: 1/64 wide, starting 1/64 after the note offset.
			if next == nil || !touching(n, next) {
				g.PitchControls = append(g.PitchControls, SVPPitchControl{
					Pos:    n.Onset + n.Duration,
					Pitch:  pitch,
					ID:     newShortID(),
					Type:   "curve",
					Points: []float64{float64(preciseOnsetEpsilon), 0, 2 * float64(preciseOnsetEpsilon), 0},
				})
			}
		}
	}
}
