package main

import (
	"math"
	"slices"
	"testing"
)

func TestCoverAmount(t *testing.T) {
	const soprano = 77 // F5
	span := float64(coverLeadIn + coverSemitones)
	tests := []struct {
		pitch int
		want  float64
	}{
		{72, 0},                         // C5, below the lead-in
		{74, coverMaxGender * 2 / span}, // D5, easing in
		{77, coverMaxGender * coverLeadIn / span}, // F5, partly covered at the passaggio
		{84, coverMaxGender},                      // C6, full cover
		{88, coverMaxGender},                      // clamped
	}
	for _, tt := range tests {
		if got := coverAmount(tt.pitch, soprano); got != tt.want {
			t.Errorf("coverAmount(%d) = %v, want %v", tt.pitch, got, tt.want)
		}
	}
}

func TestCoverGender(t *testing.T) {
	q := int64(blicksPerQuarter)
	part := Part{Name: "Soprano 1", Notes: []Note{
		{Onset: 0, Duration: q, Pitch: 72}, // C5, uncovered
		{Onset: q, Duration: q, Pitch: 84}, // C6, full cover
	}}
	base := []float64{0, defaultGender}
	points := coverGender(part, base)
	if len(points) == 0 {
		t.Fatal("expected gender points")
	}
	// The cubic curve bulges a little ahead of the covered note.
	if got := curveValueAt(points, q/2); math.Abs(got-defaultGender) > 0.05 {
		t.Errorf("gender under the low note = %v, want ~%v", got, defaultGender)
	}
	if want := defaultGender + coverMaxGender; curveValueAt(points, q+q/2) > want+0.01 {
		t.Errorf("gender under the high note = %v, want %v", curveValueAt(points, q+q/2), want)
	}
	if got := curveValueAt(points, 3*q); math.Abs(got-defaultGender) > 0.01 {
		t.Errorf("gender after the notes = %v, want %v", got, defaultGender)
	}

	// Unknown voice part keeps the curve it was given.
	if pts := coverGender(Part{Name: "Piano", Notes: part.Notes}, base); !slices.Equal(pts, base) {
		t.Errorf("unknown part: got %v, want %v", pts, base)
	}
}

func TestCoverPhonemes(t *testing.T) {
	note := func(pitch int, phonemes, lang string) *SVPNote {
		n := &SVPNote{Pitch: pitch, Phonemes: phonemes}
		n.Takes.LanguageOverride = lang
		return n
	}
	// Soprano passaggio is F5 (77), so substitution starts at B5 (83).
	group := &SVPGroup{Name: "Soprano", Notes: []*SVPNote{
		note(72, "b e a t a", "spanish"), // low, untouched
		note(84, "b e a t a", "spanish"), // high, a -> o
		note(84, "f aa r", "english"),    // high, aa -> ao
		note(84, "b e a", ""),            // unknown language, untouched
	}}
	// A part we can't identify is skipped entirely.
	other := &SVPGroup{Name: "Piano", Notes: []*SVPNote{note(84, "b e a", "spanish")}}

	coverPhonemes([]*SVPGroup{group, other})

	want := []string{"b e a t a", "b e o t o", "f ao r", "b e a"}
	for i, w := range want {
		if group.Notes[i].Phonemes != w {
			t.Errorf("note %d: got %q, want %q", i, group.Notes[i].Phonemes, w)
		}
	}
	if other.Notes[0].Phonemes != "b e a" {
		t.Errorf("unknown part: got %q, want unchanged", other.Notes[0].Phonemes)
	}
}
