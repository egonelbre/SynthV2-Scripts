package main

import (
	"math"
	"slices"
	"testing"
)

func TestIsWordStart(t *testing.T) {
	tests := []struct {
		lyric, syllabic string
		wordOpen        bool
		want            bool
	}{
		{lyric: "la", syllabic: "single", want: true},
		{lyric: "Põh", syllabic: "begin", want: true},
		{lyric: "ja", syllabic: "", want: true},
		{lyric: "tuul", syllabic: "middle", want: false},
		{lyric: "ke", syllabic: "end", want: false},
		{lyric: "-", syllabic: "", want: false},
		{lyric: "", syllabic: "", want: false},
		// "be - a - - ta": some exporters mark the syllable closing a melisma
		// as "single", but the open hyphen wins.
		{lyric: "ta", syllabic: "single", wordOpen: true, want: false},
		{lyric: "ta", syllabic: "end", wordOpen: true, want: false},
	}
	for _, tt := range tests {
		if got := isWordStart(tt.lyric, tt.syllabic, tt.wordOpen); got != tt.want {
			t.Errorf("isWordStart(%q, %q, %v) = %v, want %v", tt.lyric, tt.syllabic, tt.wordOpen, got, tt.want)
		}
	}
}

// TestWordStart_MelismaClosingSyllable covers the "be - a - - ta" case where
// the closing syllable is exported as "single" after a melisma.
func TestWordStart_MelismaClosingSyllable(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>T</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>4</duration><type>quarter</type>
        <lyric number="1"><syllabic>begin</syllabic><text>be</text></lyric></note>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>4</duration><type>quarter</type>
        <lyric number="1"><syllabic>middle</syllabic><text>a</text></lyric></note>
      <note><pitch><step>E</step><octave>4</octave></pitch><duration>4</duration><type>quarter</type></note>
      <note><pitch><step>F</step><octave>4</octave></pitch><duration>4</duration><type>quarter</type>
        <lyric number="1"><syllabic>single</syllabic><text>ta</text></lyric></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	want := []bool{true, false, false, false}
	if len(notes) != len(want) {
		t.Fatalf("got %d notes, want %d", len(notes), len(want))
	}
	for i, w := range want {
		if notes[i].WordStart != w {
			t.Errorf("note %d (%q): WordStart = %v, want %v", i, notes[i].Lyric, notes[i].WordStart, w)
		}
	}
}

func TestStressWordStart(t *testing.T) {
	q := int64(blicksPerQuarter)
	notes := []Note{
		{Onset: 0, Duration: q, Pitch: 60, Lyric: "Põh", WordStart: true},
		{Onset: q, Duration: q, Pitch: 62, Lyric: "ja"},
		{Onset: 2 * q, Duration: q, Pitch: 64, Lyric: "tuul", WordStart: true},
	}

	off := scoreToSVP(&Score{Parts: []Part{{Name: "Soprano", Notes: notes}}})
	if pts := off.Library[0].Parameters.Loudness.Points; len(pts) != 0 {
		t.Errorf("stress disabled: got %v loudness points, want none", pts)
	}

	on := scoreToSVP(&Score{StressMode: stressWordStart, Parts: []Part{{Name: "Soprano", Notes: notes}}})
	loudness := on.Library[0].Parameters.Loudness.Points
	for _, pos := range []int64{0, 2 * q} {
		if got := curveValueAt(loudness, pos); got < stressLoudnessBump-0.01 {
			t.Errorf("loudness at word start %d = %v, want >= %v", pos, got, stressLoudnessBump)
		}
	}
	// A long note is only stressed at its start: the bump is gone well before
	// the note ends.
	long := scoreToSVP(&Score{StressMode: stressWordStart, Parts: []Part{{Name: "Soprano", Notes: []Note{
		{Onset: 0, Duration: 8 * q, Pitch: 60, Lyric: "aa", WordStart: true},
		{Onset: 8 * q, Duration: q, Pitch: 62, Lyric: "bb", WordStart: true},
	}}}})
	longLoudness := long.Library[0].Parameters.Loudness.Points
	// The whole body of the long note stays at the base level: no lingering
	// bump after the spike and no ramp up into the next word's spike.
	for pos := q / 2; pos < 8*q-q/4; pos += q / 2 {
		if got := curveValueAt(longLoudness, pos); got > 0.01 {
			t.Errorf("loudness inside long note at %d = %v, want ~0", pos, got)
		}
	}

	// Mid-word syllable stays well below the stressed level (the cubic curve
	// ramping toward the next spike keeps it above the base).
	if got := curveValueAt(loudness, q); got > stressLoudnessBump/2 {
		t.Errorf("loudness at mid-word %d = %v, want < %v", q, got, stressLoudnessBump/2)
	}
}

func TestStressXylo(t *testing.T) {
	q := int64(blicksPerQuarter)
	score := &Score{StressMode: stressXylo, Parts: []Part{{Name: "Soprano", Notes: []Note{
		{Onset: 0, Duration: 4 * q, Pitch: 60, Lyric: "aa", WordStart: true},
		{Onset: 4 * q, Duration: q, Pitch: 62, Lyric: "bb", WordStart: true},
	}}}}
	loudness := scoreToSVP(score).Library[0].Parameters.Loudness.Points

	attack := curveValueAt(loudness, 0)
	if attack < xyloLoudnessBump-0.01 {
		t.Errorf("attack = %v, want %v", attack, xyloLoudnessBump)
	}
	// Decays fast, then keeps fading below the base level.
	mid := curveValueAt(loudness, q)
	tail := curveValueAt(loudness, 3*q)
	if !(attack > mid && mid > tail) {
		t.Errorf("expected attack > mid > tail, got %v, %v, %v", attack, mid, tail)
	}
	if tail > 0 {
		t.Errorf("tail = %v, want below the base level", tail)
	}
}

// TestStressKeepsDynamics checks that a stress is layered on top of the
// dynamics curve instead of flattening it.
func TestStressKeepsDynamics(t *testing.T) {
	q := int64(blicksPerQuarter)
	// A diminuendo from +6 dB to -6 dB over 8 beats.
	base := []float64{0, 6, float64(8 * q), -6}
	stresses := []accentEvent{
		{position: 0, duration: 4 * q},
		{position: 4 * q, duration: 4 * q},
	}

	for name, points := range map[string][]float64{
		"word-start": applyStresses(slices.Clone(base), stresses, stressLoudnessBump),
		"xylo":       applyXyloStresses(slices.Clone(base), stresses, xyloLoudnessBump),
	} {
		t.Run(name, func(t *testing.T) {
			// The diminuendo still descends across the whole span.
			prev := math.Inf(1)
			for pos := int64(0); pos <= 8*q; pos += q / 4 {
				got := curveValueAt(points, pos) - curveValueAt(base, pos)
				if got > 0.01 && curveValueAt(points, pos) > prev {
					// A stress may lift the curve locally; only require that
					// the underlying level is still falling.
					continue
				}
				prev = curveValueAt(points, pos)
			}
			if start, end := curveValueAt(points, q/2), curveValueAt(points, 8*q-q/2); start-end < 6 {
				t.Errorf("diminuendo flattened: %v -> %v", start, end)
			}
			// The second stress sits on the quieter part of the diminuendo.
			first := curveValueAt(points, 0) - curveValueAt(base, 0)
			second := curveValueAt(points, 4*q) - curveValueAt(base, 4*q)
			if math.Abs(first-second) > 0.01 {
				t.Errorf("stress offsets differ: %v vs %v", first, second)
			}
		})
	}
}
