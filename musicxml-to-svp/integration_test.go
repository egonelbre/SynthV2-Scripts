package main

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

// TestIntegration_MelismaticMusicXML runs the full pipeline on a real MusicXML
// file and verifies key properties of the output.
func TestIntegration_MelismaticMusicXML(t *testing.T) {
	data, err := os.ReadFile("testdata/Melismatic.musicxml")
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	var score musicxml.ScorePartwise
	if err := xml.Unmarshal(data, &score); err != nil {
		t.Fatalf("failed to parse MusicXML: %v", err)
	}

	if len(score.Part) != 1 {
		t.Fatalf("expected 1 part, got %d", len(score.Part))
	}

	// Pass 1: structure
	unrolled, meters, tempos := buildStructure(score.Part[0])

	// File has 4 measures with a 3x repeat and volta endings.
	// Should have more than 4 unrolled measures.
	if len(unrolled) < 4 {
		t.Errorf("expected at least 4 unrolled measures, got %d", len(unrolled))
	}

	// Should detect 4/4 time.
	if len(meters) != 1 {
		t.Fatalf("expected 1 meter, got %d", len(meters))
	}
	if meters[0].Numerator != 4 || meters[0].Denominator != 4 {
		t.Errorf("expected 4/4 meter, got %d/%d", meters[0].Numerator, meters[0].Denominator)
	}

	// No explicit tempo in this file.
	if len(tempos) != 0 {
		t.Errorf("expected 0 tempos (none specified), got %d", len(tempos))
	}

	// Pass 2: notes
	notes := buildNotes(score.Part[0], unrolled)
	if len(notes) == 0 {
		t.Fatal("expected at least 1 note")
	}

	// All notes should have non-negative onset and positive duration.
	for i, n := range notes {
		if n.Onset < 0 {
			t.Errorf("note %d: negative onset %d", i, n.Onset)
		}
		if n.Duration <= 0 {
			t.Errorf("note %d: non-positive duration %d", i, n.Duration)
		}
		if n.Pitch < 0 || n.Pitch > 127 {
			t.Errorf("note %d: MIDI pitch %d out of range [0,127]", i, n.Pitch)
		}
	}

	// Pass 3: lyrics
	fillLyrics(notes)

	// Every note should have a non-empty lyric after fillLyrics.
	for i, n := range notes {
		if n.Lyric == "" {
			t.Errorf("note %d: empty lyric after fillLyrics", i)
		}
	}

	// Pass 4: dynamics
	dynamics := buildDynamics(score.Part[0], unrolled)
	// Dynamics may or may not exist in this file; just verify no panic.
	_ = dynamics

	// Pass 5: convert to SVP
	irScore := &Score{
		Meters: meters,
		Tempos: tempos,
		Parts: []Part{{
			Name:     "Soprano",
			Notes:    notes,
			Dynamics: dynamics,
		}},
	}
	project := scoreToSVP(irScore)

	if len(project.Tracks) != 1 {
		t.Errorf("expected 1 track, got %d", len(project.Tracks))
	}
	if len(project.Library) < 1 {
		t.Fatal("expected at least 1 library group")
	}
	if len(project.Library[0].Notes) == 0 {
		t.Error("expected SVP notes in first library group")
	}

	// Verify SVP defaults.
	if project.Version != 196 {
		t.Errorf("expected version 196, got %d", project.Version)
	}
	if len(project.Time.Meters) == 0 {
		t.Error("expected at least 1 SVP meter")
	}
	if len(project.Time.Tempos) == 0 {
		t.Error("expected at least 1 SVP tempo (default 120)")
	}
}
