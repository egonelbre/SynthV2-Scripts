package main

import (
	"strconv"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

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
func buildNotes(part *musicxml.Part, unrolled []playedMeasure, infos []measureInfo) []Note {
	divisions := 4
	transpose := 0
	cursor := int64(0)
	var notes []Note
	pendingTies := pendingTieMap{}
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

				validateTimeModification(value, divisions)

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

				// Save original blicks for cursor advancement (before grace adjustment).
				cursorBlicks := blicks

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
					if idx, ok := pendingTies.find(midi); ok {
						notes[idx].Duration += cursorBlicks
						if !tieStart {
							pendingTies.remove(midi)
						}
					}
					if value.Chord == "" {
						prevOnset = cursor
						cursor += cursorBlicks
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
					pendingTies.add(midi, lastNoteIdx)
				}

				if value.Chord == "" {
					prevOnset = cursor
					cursor += cursorBlicks
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
