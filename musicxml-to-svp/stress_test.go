package main

import "testing"

func TestIsWordStart(t *testing.T) {
	tests := []struct {
		lyric, syllabic string
		want            bool
	}{
		{"la", "single", true},
		{"Põh", "begin", true},
		{"ja", "", true},
		{"tuul", "middle", false},
		{"ke", "end", false},
		{"-", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		if got := isWordStart(tt.lyric, tt.syllabic); got != tt.want {
			t.Errorf("isWordStart(%q, %q) = %v, want %v", tt.lyric, tt.syllabic, got, tt.want)
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

	on := scoreToSVP(&Score{StressWordStart: true, Parts: []Part{{Name: "Soprano", Notes: notes}}})
	loudness := on.Library[0].Parameters.Loudness.Points
	for _, pos := range []int64{0, 2 * q} {
		if got := curveValueAt(loudness, pos); got < stressLoudnessBump-0.01 {
			t.Errorf("loudness at word start %d = %v, want >= %v", pos, got, stressLoudnessBump)
		}
	}
	// A long note is only stressed at its start: the bump is gone well before
	// the note ends.
	long := scoreToSVP(&Score{StressWordStart: true, Parts: []Part{{Name: "Soprano", Notes: []Note{
		{Onset: 0, Duration: 8 * q, Pitch: 60, Lyric: "aa", WordStart: true},
	}}}})
	longLoudness := long.Library[0].Parameters.Loudness.Points
	if got := curveValueAt(longLoudness, q/2); got > 0.01 {
		t.Errorf("loudness half a beat into a long note = %v, want ~0", got)
	}

	// Mid-word syllable stays well below the stressed level (the cubic curve
	// ramping toward the next spike keeps it above the base).
	if got := curveValueAt(loudness, q); got > stressLoudnessBump/2 {
		t.Errorf("loudness at mid-word %d = %v, want < %v", q, got, stressLoudnessBump/2)
	}
}
