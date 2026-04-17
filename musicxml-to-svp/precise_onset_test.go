package main

import "testing"

// findControlsAt returns the pitch controls anchored at pos for the given pitch.
func findControlsAt(controls []SVPPitchControl, pos int64, pitch float64) []SVPPitchControl {
	var out []SVPPitchControl
	for _, c := range controls {
		if c.Pos == pos && c.Pitch == pitch {
			out = append(out, c)
		}
	}
	return out
}

// TestAddPreciseOnsetControls_NonTouchingNotes verifies that two notes with a
// real gap between them get an onset control before each note and a phrase-end
// control after each note, all placed one epsilon outside the note edge.
func TestAddPreciseOnsetControls_NonTouchingNotes(t *testing.T) {
	gap := 8 * preciseOnsetEpsilon // ≥ 4ε → fits offset+onset controls without overlap
	noteDur := int64(blicksPerQuarter)
	library := []*SVPGroup{{
		Notes: []*SVPNote{
			{Onset: 0, Duration: noteDur, Pitch: 60},
			{Onset: noteDur + gap, Duration: noteDur, Pitch: 62},
		},
	}}

	addPreciseOnsetControls(library)
	pcs := library[0].PitchControls

	// Both notes are at phrase boundaries on both sides → 4 controls total.
	if len(pcs) != 4 {
		t.Fatalf("expected 4 pitch controls, got %d", len(pcs))
	}

	wantOnsetPoints := []float64{-2 * float64(preciseOnsetEpsilon), 0, -float64(preciseOnsetEpsilon), 0}
	wantOffsetPoints := []float64{float64(preciseOnsetEpsilon), 0, 2 * float64(preciseOnsetEpsilon), 0}

	c := findControlsAt(pcs, 0, 60)
	if len(c) != 1 || !floatSliceEq(c[0].Points, wantOnsetPoints) {
		t.Errorf("note 0 onset: want pos=0 points=%v, got %+v", wantOnsetPoints, c)
	}
	c = findControlsAt(pcs, noteDur, 60)
	if len(c) != 1 || !floatSliceEq(c[0].Points, wantOffsetPoints) {
		t.Errorf("note 0 offset: want pos=%d points=%v, got %+v", noteDur, wantOffsetPoints, c)
	}
	c = findControlsAt(pcs, noteDur+gap, 62)
	if len(c) != 1 || !floatSliceEq(c[0].Points, wantOnsetPoints) {
		t.Errorf("note 1 onset: want points=%v, got %+v", wantOnsetPoints, c)
	}
	c = findControlsAt(pcs, 2*noteDur+gap, 62)
	if len(c) != 1 || !floatSliceEq(c[0].Points, wantOffsetPoints) {
		t.Errorf("note 1 offset: want points=%v, got %+v", wantOffsetPoints, c)
	}
}

// TestAddPreciseOnsetControls_TouchingNotes verifies that consecutive notes in
// the same phrase are skipped on the touching side: no offset control on the
// first note, no onset control on the second.
func TestAddPreciseOnsetControls_TouchingNotes(t *testing.T) {
	noteDur := int64(blicksPerQuarter)
	gap := 3 * preciseOnsetEpsilon // < 4ε → controls would overlap, so skip both
	library := []*SVPGroup{{
		Notes: []*SVPNote{
			{Onset: 0, Duration: noteDur, Pitch: 60},
			{Onset: noteDur + gap, Duration: noteDur, Pitch: 62},
		},
	}}

	addPreciseOnsetControls(library)
	pcs := library[0].PitchControls

	// First note: onset only (no offset because next is touching).
	// Second note: offset only (no onset because prev is touching).
	if len(pcs) != 2 {
		t.Fatalf("expected 2 pitch controls, got %d: %+v", len(pcs), pcs)
	}
	if c := findControlsAt(pcs, 0, 60); len(c) != 1 {
		t.Errorf("expected onset control on first note, got %d", len(c))
	}
	if c := findControlsAt(pcs, noteDur, 60); len(c) != 0 {
		t.Errorf("expected no offset control on first note (phrase continuation), got %d", len(c))
	}
	if c := findControlsAt(pcs, noteDur+gap, 62); len(c) != 0 {
		t.Errorf("expected no onset control on second note (phrase continuation), got %d", len(c))
	}
	if c := findControlsAt(pcs, 2*noteDur+gap, 62); len(c) != 1 {
		t.Errorf("expected offset control on second note, got %d", len(c))
	}
}

// TestAddPreciseOnsetControls_AppliesPerGroup verifies that controls are added
// independently to each library group.
func TestAddPreciseOnsetControls_AppliesPerGroup(t *testing.T) {
	library := []*SVPGroup{
		{Notes: []*SVPNote{{Onset: 0, Duration: blicksPerQuarter, Pitch: 60}}},
		{Notes: []*SVPNote{{Onset: 0, Duration: blicksPerQuarter, Pitch: 64}}},
	}

	addPreciseOnsetControls(library)

	// Single isolated note: onset + offset = 2 controls per group.
	if len(library[0].PitchControls) != 2 {
		t.Errorf("group 0: expected 2 controls, got %d", len(library[0].PitchControls))
	}
	if len(library[1].PitchControls) != 2 {
		t.Errorf("group 1: expected 2 controls, got %d", len(library[1].PitchControls))
	}
	for _, pc := range library[0].PitchControls {
		if pc.Pitch != 60 {
			t.Errorf("group 0 control got wrong pitch %v", pc.Pitch)
		}
	}
	for _, pc := range library[1].PitchControls {
		if pc.Pitch != 64 {
			t.Errorf("group 1 control got wrong pitch %v", pc.Pitch)
		}
	}
}

func floatSliceEq(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
