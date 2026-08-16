package main

import "testing"

func TestCoverAmount(t *testing.T) {
	const soprano = 77 // F5
	tests := []struct {
		pitch int
		want  float64
	}{
		{74, 0},                               // D5, below the passaggio
		{76, 0},                               // E5, still below
		{77, coverMaxGender / coverSemitones}, // F5, just entering
		{83, coverMaxGender},                  // B5, full cover
		{88, coverMaxGender},                  // clamped
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
	points := coverGender(part)
	if len(points) == 0 {
		t.Fatal("expected gender points")
	}
	if got := curveValueAt(points, q/2); got != 0 {
		t.Errorf("gender under the low note = %v, want 0", got)
	}
	if got := curveValueAt(points, q+q/2); got > coverMaxGender+0.01 {
		t.Errorf("gender under the high note = %v, want %v", got, coverMaxGender)
	}
	if got := curveValueAt(points, 3*q); got != 0 {
		t.Errorf("gender after the notes = %v, want 0", got)
	}

	// Unknown voice part is left alone.
	if pts := coverGender(Part{Name: "Piano", Notes: part.Notes}); pts != nil {
		t.Errorf("unknown part: got %v, want nil", pts)
	}
}

func TestCoverPhonemes(t *testing.T) {
	note := func(pitch int, phonemes, lang string) *SVPNote {
		n := &SVPNote{Pitch: pitch, Phonemes: phonemes}
		n.Takes.LanguageOverride = lang
		return n
	}
	// Soprano passaggio is F5 (77), so substitution starts at G#5 (80).
	group := &SVPGroup{Name: "Soprano", Notes: []*SVPNote{
		note(72, "b e a t a", "spanish"), // low, untouched
		note(80, "b e a t a", "spanish"), // high, a -> o
		note(80, "f aa r", "english"),    // high, aa -> ao
		note(80, "b e a", ""),            // unknown language, untouched
	}}
	// A part we can't identify is skipped entirely.
	other := &SVPGroup{Name: "Piano", Notes: []*SVPNote{note(80, "b e a", "spanish")}}

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
