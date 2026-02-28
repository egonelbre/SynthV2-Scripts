package main

import (
	"testing"
)

// TestScoreToSVP_ArticulationShortening tests staccato/staccatissimo duration adjustment.
func TestScoreToSVP_ArticulationShortening(t *testing.T) {
	s := &Score{
		Meters: []MeterChange{{MeasureIndex: 0, Numerator: 4, Denominator: 4}},
		Tempos: []TempoChange{{Position: 0, BPM: 120}},
		Parts: []Part{{
			Name: "Test",
			Notes: []Note{
				{Onset: 0, Duration: blicksPerQuarter, Pitch: 60, Lyric: "a", Articulations: ArticulationStaccato},
				{Onset: blicksPerQuarter, Duration: blicksPerQuarter, Pitch: 62, Lyric: "b", Articulations: ArticulationStaccatissimo},
				{Onset: 2 * blicksPerQuarter, Duration: blicksPerQuarter, Pitch: 64, Lyric: "c", Articulations: ArticulationTenuto},
				{Onset: 3 * blicksPerQuarter, Duration: blicksPerQuarter, Pitch: 65, Lyric: "d"},
			},
		}},
	}

	proj := scoreToSVP(s)
	notes := proj.Library[0].Notes

	if len(notes) != 4 {
		t.Fatalf("expected 4 notes, got %d", len(notes))
	}

	// Staccato: 2/3 duration
	expectedStaccato := int64(blicksPerQuarter * 2 / 3)
	if notes[0].Duration != expectedStaccato {
		t.Errorf("staccato duration: expected %d, got %d", expectedStaccato, notes[0].Duration)
	}

	// Staccatissimo: 1/3 duration
	expectedStaccatissimo := int64(blicksPerQuarter / 3)
	if notes[1].Duration != expectedStaccatissimo {
		t.Errorf("staccatissimo duration: expected %d, got %d", expectedStaccatissimo, notes[1].Duration)
	}

	// Tenuto: full duration
	if notes[2].Duration != blicksPerQuarter {
		t.Errorf("tenuto duration: expected %d, got %d", blicksPerQuarter, notes[2].Duration)
	}

	// No articulation: full duration
	if notes[3].Duration != blicksPerQuarter {
		t.Errorf("plain duration: expected %d, got %d", blicksPerQuarter, notes[3].Duration)
	}
}

// TestScoreToSVP_TenutoOverridesStaccato tests that tenuto suppresses staccato shortening.
func TestScoreToSVP_TenutoOverridesStaccato(t *testing.T) {
	s := &Score{
		Parts: []Part{{
			Name: "Test",
			Notes: []Note{
				{Onset: 0, Duration: blicksPerQuarter, Pitch: 60, Lyric: "a",
					Articulations: ArticulationTenuto | ArticulationStaccato},
			},
		}},
	}

	proj := scoreToSVP(s)
	notes := proj.Library[0].Notes

	// Tenuto should prevent staccato shortening
	if notes[0].Duration != blicksPerQuarter {
		t.Errorf("tenuto+staccato duration: expected %d (full), got %d", blicksPerQuarter, notes[0].Duration)
	}
}

// TestScoreToSVP_GraceNoteEmission tests that grace notes are emitted correctly.
func TestScoreToSVP_GraceNoteEmission(t *testing.T) {
	graceDur := int64(blicksPerQuarter / 4) // sixteenth
	mainOnset := graceDur
	mainDur := int64(blicksPerQuarter) - graceDur

	s := &Score{
		Parts: []Part{{
			Name: "Test",
			Notes: []Note{
				{
					Onset:    mainOnset,
					Duration: mainDur,
					Pitch:    60,
					Lyric:    "la",
					LeadingGraces: []GraceNote{
						{Pitch: 62, Lyric: "la", Duration: graceDur},
					},
				},
			},
		}},
	}

	proj := scoreToSVP(s)
	notes := proj.Library[0].Notes

	if len(notes) != 2 {
		t.Fatalf("expected 2 SVP notes (grace + main), got %d", len(notes))
	}

	// Grace note should be emitted at mainOnset - graceDur = 0
	if notes[0].Onset != 0 {
		t.Errorf("grace onset: expected 0, got %d", notes[0].Onset)
	}
	if notes[0].Duration != graceDur {
		t.Errorf("grace duration: expected %d, got %d", graceDur, notes[0].Duration)
	}
	if notes[0].Pitch != 62 {
		t.Errorf("grace pitch: expected 62, got %d", notes[0].Pitch)
	}

	// Main note
	if notes[1].Onset != mainOnset {
		t.Errorf("main onset: expected %d, got %d", mainOnset, notes[1].Onset)
	}
}

// TestScoreToSVP_TrailingGraceEmission tests trailing grace note emission.
func TestScoreToSVP_TrailingGraceEmission(t *testing.T) {
	graceDur := int64(blicksPerQuarter / 4)
	mainDur := int64(blicksPerQuarter) - graceDur

	s := &Score{
		Parts: []Part{{
			Name: "Test",
			Notes: []Note{
				{
					Onset:    0,
					Duration: mainDur,
					Pitch:    60,
					Lyric:    "la",
					TrailingGraces: []GraceNote{
						{Pitch: 62, Lyric: "-", Duration: graceDur},
					},
				},
			},
		}},
	}

	proj := scoreToSVP(s)
	notes := proj.Library[0].Notes

	if len(notes) != 2 {
		t.Fatalf("expected 2 SVP notes (main + trailing grace), got %d", len(notes))
	}

	// Main note
	if notes[0].Onset != 0 {
		t.Errorf("main onset: expected 0, got %d", notes[0].Onset)
	}
	if notes[0].Duration != mainDur {
		t.Errorf("main duration: expected %d, got %d", mainDur, notes[0].Duration)
	}

	// Trailing grace starts right after main note
	if notes[1].Onset != mainDur {
		t.Errorf("trailing grace onset: expected %d, got %d", mainDur, notes[1].Onset)
	}
}

// TestScoreToSVP_TrailingGraceWithStaccato tests that trailing grace notes
// are placed at the pre-staccato endpoint, not the shortened one (bug #4).
func TestScoreToSVP_TrailingGraceWithStaccato(t *testing.T) {
	graceDur := int64(blicksPerQuarter / 4)
	mainDur := int64(blicksPerQuarter) - graceDur // pre-staccato duration

	s := &Score{
		Parts: []Part{{
			Name: "Test",
			Notes: []Note{
				{
					Onset:         0,
					Duration:      mainDur,
					Pitch:         60,
					Lyric:         "la",
					Articulations: ArticulationStaccato,
					TrailingGraces: []GraceNote{
						{Pitch: 62, Lyric: "-", Duration: graceDur},
					},
				},
			},
		}},
	}

	proj := scoreToSVP(s)
	notes := proj.Library[0].Notes

	if len(notes) != 2 {
		t.Fatalf("expected 2 SVP notes (main + trailing grace), got %d", len(notes))
	}

	// Main note should have staccato-shortened duration
	expectedMainDur := mainDur * 2 / 3
	if notes[0].Duration != expectedMainDur {
		t.Errorf("main duration: expected %d (staccato-shortened), got %d", expectedMainDur, notes[0].Duration)
	}

	// Trailing grace should start at the original (pre-staccato) endpoint
	expectedTrailOnset := mainDur
	if notes[1].Onset != expectedTrailOnset {
		t.Errorf("trailing grace onset: expected %d (pre-staccato end), got %d", expectedTrailOnset, notes[1].Onset)
	}
}

// TestScoreToSVP_AccentSpikes tests that accents generate accent events.
func TestScoreToSVP_AccentSpikes(t *testing.T) {
	s := &Score{
		Parts: []Part{{
			Name: "Test",
			Notes: []Note{
				{Onset: 0, Duration: blicksPerQuarter, Pitch: 60, Lyric: "a",
					Articulations: ArticulationAccent},
			},
			Dynamics: []dynEvent{
				{position: 0, kind: dynLevel, loudness: 3, tension: 0},
			},
		}},
	}

	proj := scoreToSVP(s)
	params := proj.Library[0].Parameters

	// Loudness curve should have accent spike points
	if len(params.Loudness.Points) < 4 {
		t.Errorf("expected loudness curve with accent spike, got %d points", len(params.Loudness.Points))
	}
}

// TestCapGraceDurs_RoundingRemainder tests that grace note duration scaling
// doesn't lose blicks to float64→int64 truncation (bug #6).
func TestCapGraceDurs_RoundingRemainder(t *testing.T) {
	// 3 equal grace notes capped to 100 blicks.
	// Without fix: 33+33+33=99. With fix: 33+33+34=100.
	graces := []GraceNote{
		{NotatedType: "quarter"},
		{NotatedType: "quarter"},
		{NotatedType: "quarter"},
	}
	maxTotal := int64(100)

	graceDurs, totalGrace := capGraceDurs(graces, maxTotal)

	if totalGrace != maxTotal {
		t.Errorf("totalGrace: expected %d, got %d (lost %d blicks)",
			maxTotal, totalGrace, maxTotal-totalGrace)
	}

	// Verify individual durations sum to total
	var sum int64
	for _, d := range graceDurs {
		sum += d
	}
	if sum != maxTotal {
		t.Errorf("sum of graceDurs: expected %d, got %d", maxTotal, sum)
	}
}
