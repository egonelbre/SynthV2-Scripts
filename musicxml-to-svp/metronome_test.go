package main

import (
	"math"
	"testing"
)

func TestMetronomeBeatsPickup(t *testing.T) {
	q := int64(blicksPerQuarter)
	// One-beat pickup in 3/4, then a full measure.
	unrolled := []playedMeasure{{startBlicks: 0}, {startBlicks: q}}
	meters := []MeterChange{{MeasureIndex: 0, Numerator: 3, Denominator: 4}}
	beats := metronomeBeats(unrolled, meters, 4*q)
	want := []metronomeBeat{{0, false}, {q, true}, {2 * q, false}, {3 * q, false}}
	if len(beats) != len(want) {
		t.Fatalf("got %v, want %v", beats, want)
	}
	for i := range want {
		if beats[i] != want[i] {
			t.Errorf("beat %d: got %v, want %v", i, beats[i], want[i])
		}
	}
}

func TestBlicksToSeconds(t *testing.T) {
	q := int64(blicksPerQuarter)
	tempos := []TempoChange{{0, 120}, {2 * q, 60}}
	// Two quarters at 120 = 1s, then one quarter at 60 = 1s.
	if got := blicksToSeconds(tempos, 3*q); math.Abs(got-2) > 1e-9 {
		t.Errorf("got %v, want 2", got)
	}
}
