package main

import "testing"

func TestSplitPart(t *testing.T) {
	n := func(onset, dur int64, pitch int) Note {
		return Note{Onset: onset, Duration: dur, Pitch: pitch}
	}

	t.Run("no overlap", func(t *testing.T) {
		p := Part{Name: "Soprano", Notes: []Note{n(0, 10, 60), n(10, 10, 62)}}
		got := splitPart(p)
		if len(got) != 1 || got[0].Name != "Soprano" {
			t.Fatalf("got %v, want single unchanged part", got)
		}
	})

	t.Run("two voices", func(t *testing.T) {
		// Voice 1: two half notes; Voice 2: overlapping lower notes.
		p := Part{Name: "Alto", Notes: []Note{
			n(0, 20, 60), n(5, 20, 48),
			n(20, 20, 62), n(25, 20, 50),
		}}
		got := splitPart(p)
		if len(got) != 2 {
			t.Fatalf("got %d parts, want 2", len(got))
		}
		if got[0].Name != "Alto (Voice 1)" || got[1].Name != "Alto (Voice 2)" {
			t.Fatalf("names: %q, %q", got[0].Name, got[1].Name)
		}
		for i, want := range [][]int{{60, 62}, {48, 50}} {
			if len(got[i].Notes) != len(want) {
				t.Fatalf("voice %d: got %d notes, want %d", i+1, len(got[i].Notes), len(want))
			}
			for j, pitch := range want {
				if got[i].Notes[j].Pitch != pitch {
					t.Errorf("voice %d note %d: got pitch %d, want %d", i+1, j, got[i].Notes[j].Pitch, pitch)
				}
			}
		}
	})

	t.Run("chord splits by pitch", func(t *testing.T) {
		p := Part{Name: "Bass", Notes: []Note{n(0, 10, 48), n(0, 10, 52), n(0, 10, 55)}}
		got := splitPart(p)
		if len(got) != 3 {
			t.Fatalf("got %d parts, want 3", len(got))
		}
		for i, pitch := range []int{55, 52, 48} {
			if got[i].Notes[0].Pitch != pitch {
				t.Errorf("voice %d: got pitch %d, want %d", i+1, got[i].Notes[0].Pitch, pitch)
			}
		}
	})
}
