package main

import (
	"math"
	"testing"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/voice"
)

// TestComputePanning_Center tests that center scheme sets all known parts to 0.
func TestComputePanning_Center(t *testing.T) {
	infos := []voice.TrackInfo{
		{Name: "Soprano", Part: voice.Soprano, PartNum: 1},
		{Name: "Alto", Part: voice.Alto, PartNum: 1},
		{Name: "Piano", Part: voice.Unknown, PartNum: 1},
	}

	panning := computePanning(infos, "center")

	if pan, ok := panning[0]; !ok || pan != 0 {
		t.Errorf("soprano pan: expected 0, got %f (ok=%v)", pan, ok)
	}
	if pan, ok := panning[1]; !ok || pan != 0 {
		t.Errorf("alto pan: expected 0, got %f (ok=%v)", pan, ok)
	}
	if _, ok := panning[2]; ok {
		t.Error("unknown part should not get a pan value with center scheme")
	}
}

// TestComputePanning_Spread tests the spread panning scheme.
func TestComputePanning_Spread(t *testing.T) {
	infos := []voice.TrackInfo{
		{Name: "Soprano", Part: voice.Soprano, PartNum: 1},
		{Name: "Alto", Part: voice.Alto, PartNum: 1},
		{Name: "Bass", Part: voice.Bass, PartNum: 1},
	}

	panning := computePanning(infos, "spread")

	if pan := panning[0]; pan != -0.6 {
		t.Errorf("soprano pan: expected -0.6, got %f", pan)
	}
	if pan := panning[1]; pan != 0.6 {
		t.Errorf("alto pan: expected 0.6, got %f", pan)
	}
	if pan := panning[2]; pan != 0.3 {
		t.Errorf("bass pan: expected 0.3, got %f", pan)
	}
}

// TestComputePanning_Default tests the default choir seating panning.
func TestComputePanning_Default(t *testing.T) {
	infos := []voice.TrackInfo{
		{Name: "Soprano", Part: voice.Soprano, PartNum: 1},
		{Name: "Alto", Part: voice.Alto, PartNum: 1},
		{Name: "Tenor", Part: voice.Tenor, PartNum: 1},
		{Name: "Bass", Part: voice.Bass, PartNum: 1},
	}

	panning := computePanning(infos, "default")

	// With 4 tracks in seating order S, B, T, A, they should span -0.65 to +0.65.
	if len(panning) != 4 {
		t.Fatalf("expected 4 panning values, got %d", len(panning))
	}

	// Soprano should be leftmost.
	if panning[0] > panning[1] {
		t.Errorf("soprano (%f) should be left of alto (%f)", panning[0], panning[1])
	}

	// All values should be within [-0.65, 0.65].
	for i, pan := range panning {
		if math.Abs(pan) > 0.66 {
			t.Errorf("track %d pan %f exceeds range [-0.65, 0.65]", i, pan)
		}
	}
}

// TestComputePanning_SingleTrack tests that a single track is centered.
func TestComputePanning_SingleTrack(t *testing.T) {
	infos := []voice.TrackInfo{
		{Name: "Soprano", Part: voice.Soprano, PartNum: 1},
	}

	panning := computePanning(infos, "default")

	if pan, ok := panning[0]; !ok || pan != 0 {
		t.Errorf("single track pan: expected 0, got %f", pan)
	}
}

// TestComputePanning_DefaultSeatingOrder tests that the default seating order
// places S on the left and A on the right, following traditional choir seating.
func TestComputePanning_DefaultSeatingOrder(t *testing.T) {
	infos := []voice.TrackInfo{
		{Name: "Soprano", Part: voice.Soprano, PartNum: 1},
		{Name: "Alto", Part: voice.Alto, PartNum: 1},
	}

	panning := computePanning(infos, "default")

	// Soprano should be left of alto.
	if panning[0] >= panning[1] {
		t.Errorf("soprano (%f) should be left of alto (%f)", panning[0], panning[1])
	}
}
