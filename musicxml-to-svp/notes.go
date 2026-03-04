package main

import (
	"strconv"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

func noteArticulations(note *musicxml.Note) Articulation {
	var a Articulation
	for _, n := range note.Notations {
		for _, art := range n.Articulations {
			if len(art.Staccato) > 0 {
				a |= ArticulationStaccato
			}
			if len(art.Staccatissimo) > 0 {
				a |= ArticulationStaccatissimo
			}
			if len(art.Tenuto) > 0 {
				a |= ArticulationTenuto
			}
			if len(art.Accent) > 0 {
				a |= ArticulationAccent
			}
			if len(art.StrongAccent) > 0 {
				a |= ArticulationStrongAccent
			}
		}
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

// pendingTieMap tracks pending ties per MIDI pitch using FIFO ordering.
// Multiple voices may have ties at the same pitch; FIFO ensures each
// tie-stop matches the earliest unresolved tie-start at that pitch.
type pendingTieMap map[int][]int // MIDI pitch -> FIFO list of note indices

func (m pendingTieMap) add(pitch, noteIdx int) {
	m[pitch] = append(m[pitch], noteIdx)
}

func (m pendingTieMap) find(pitch int) (int, bool) {
	if idxs := m[pitch]; len(idxs) > 0 {
		return idxs[0], true
	}
	return 0, false
}

func (m pendingTieMap) remove(pitch int) {
	if idxs := m[pitch]; len(idxs) > 1 {
		m[pitch] = idxs[1:]
	} else {
		delete(m, pitch)
	}
}

// buildNotes extracts notes from a part, resolving ties and attaching grace notes.
func buildNotes(part *musicxml.Part, unrolled []playedMeasure) []Note {
	transpose := 0
	var notes []Note
	pendingTies := pendingTieMap{}
	var prevOnset int64
	var pendingGraces []*musicxml.Note
	var graceIsTail bool
	var lastNoteIdx int = -1
	var lastVerse int
	pendingSlideIdx := -1 // index of note with slide-start

	walkPartElements(part, unrolled, func(cursor int64, divisions int, pm playedMeasure, value any) {
		lastVerse = pm.verse

		switch value := value.(type) {
		case *musicxml.Attributes:
			for _, tr := range value.Transpose {
				chromatic, _ := strconv.Atoi(tr.Chromatic)
				transpose = chromatic + tr.OctaveChange*12
			}

		case *musicxml.Note:
			if value.Grace != nil {
				if len(pendingGraces) == 0 {
					graceIsTail = lastNoteIdx >= 0
				}
				pendingGraces = append(pendingGraces, value)
				return
			}

			validateTimeModification(value, divisions)

			dur := parseDuration(value.Duration)
			if dur <= 0 {
				return
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
				}
				return
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
			slideStart, slideStop := noteSlideTypes(value)

			if tieStop {
				// Extend the tied note's duration. Any grace notes before
				// this tie-stop were already flushed as trailing graces to
				// the previous note (or discarded if no previous note exists).
				// The leading graces built above are intentionally unused
				// since they don't attach to a tie continuation.
				if idx, ok := pendingTies.find(midi); ok {
					notes[idx].Duration += durationToBlicks(dur, divisions)
					if !tieStart {
						pendingTies.remove(midi)
					}
				}
				if value.Chord == "" {
					prevOnset = cursor
				}
				return
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
				pendingTies.add(midi, lastNoteIdx)
			}

			if slideStop && pendingSlideIdx >= 0 {
				startNote := &notes[pendingSlideIdx]
				startNote.SlideDelta = (midi-startNote.Pitch)*100 + (detune - startNote.Detune)
				pendingSlideIdx = -1
			}
			if slideStart {
				pendingSlideIdx = lastNoteIdx
			}

			if value.Chord == "" {
				prevOnset = cursor
			}
		}
	})

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
