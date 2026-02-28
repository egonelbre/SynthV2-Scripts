package main

import (
	"testing"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

// TestBuildDynamics tests dynamics event extraction.
func TestBuildDynamics_LevelsAndWedges(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <direction><direction-type><dynamics><mf/></dynamics></direction-type></direction>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>8</duration><type>half</type></note>
      <direction><direction-type><wedge type="crescendo"/></direction-type></direction>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>8</duration><type>half</type></note>
    </measure>
    <measure>
      <direction><direction-type><wedge type="stop"/></direction-type></direction>
      <direction><direction-type><dynamics><f/></dynamics></direction-type></direction>
      <note><pitch><step>E</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	events := buildDynamics(score.Part[0], unrolled)

	if len(events) < 3 {
		t.Fatalf("expected at least 3 dynamic events, got %d", len(events))
	}

	// First event: mf level
	if events[0].kind != dynLevel {
		t.Errorf("event 0: expected dynLevel, got %d", events[0].kind)
	}
	if events[0].loudness != 3 { // mf = 3 dB
		t.Errorf("event 0 loudness: expected 3, got %f", events[0].loudness)
	}

	// Second event: crescendo start
	if events[1].kind != dynCrescStart {
		t.Errorf("event 1: expected dynCrescStart, got %d", events[1].kind)
	}

	// Third event: wedge stop
	if events[2].kind != dynWedgeStop {
		t.Errorf("event 2: expected dynWedgeStop, got %d", events[2].kind)
	}
}

// TestCurveValueAt_CubicInterpolation tests that curveValueAt uses cubic
// interpolation rather than linear, matching SVP's "cubic" curve mode.
func TestCurveValueAt_CubicInterpolation(t *testing.T) {
	// Three-point curve: 0→0, 100→10, 200→0 (a peak).
	// At the midpoint (50), linear would give 5.0.
	// Cubic (Catmull-Rom) should differ from linear due to curvature.
	points := []float64{0, 0, 100, 10, 200, 0}

	val := curveValueAt(points, 50)

	// With Catmull-Rom and clamped boundaries, the curve overshoots slightly
	// compared to linear. The exact value depends on tangent computation,
	// but it should NOT be exactly 5.0 (which would indicate linear).
	if val == 5.0 {
		t.Error("curveValueAt returned exactly 5.0, suggesting linear interpolation instead of cubic")
	}

	// Should be in a reasonable range around 5.
	if val < 3.0 || val > 8.0 {
		t.Errorf("curveValueAt at midpoint: expected ~5-6, got %f", val)
	}

	// Endpoints should return exact values.
	if v := curveValueAt(points, 0); v != 0.0 {
		t.Errorf("curveValueAt at start: expected 0, got %f", v)
	}
	if v := curveValueAt(points, 100); v != 10.0 {
		t.Errorf("curveValueAt at middle point: expected 10, got %f", v)
	}
	if v := curveValueAt(points, 200); v != 0.0 {
		t.Errorf("curveValueAt at end: expected 0, got %f", v)
	}

	// Before/after curve should clamp.
	if v := curveValueAt(points, -50); v != 0.0 {
		t.Errorf("curveValueAt before curve: expected 0, got %f", v)
	}
	if v := curveValueAt(points, 300); v != 0.0 {
		t.Errorf("curveValueAt after curve: expected 0, got %f", v)
	}
}

// TestCurveValueAt_TwoPoints tests that curveValueAt works correctly with
// only two points (degenerates to linear since no neighboring points exist).
func TestCurveValueAt_TwoPoints(t *testing.T) {
	points := []float64{0, 0, 100, 10}
	val := curveValueAt(points, 50)
	if val != 5.0 {
		t.Errorf("curveValueAt with 2 points at midpoint: expected 5.0, got %f", val)
	}
}

// TestBuildCurve_WedgeTargetWithClusteredEvents tests that a dynLevel after
// a wedge stop is correctly identified as the ramp target even when other
// events cluster at the same position.
func TestBuildCurve_WedgeTargetWithClusteredEvents(t *testing.T) {
	// Events: mf at 0, cresc at Q, wedge-stop at 2Q, another-wedge-stop at 2Q, f at 2Q.
	// The f at 2Q should be detected as the wedge target (no step transition).
	events := []dynEvent{
		{position: 0, kind: dynLevel, loudness: 3},
		{position: blicksPerQuarter, kind: dynCrescStart, number: 1},
		{position: 2 * blicksPerQuarter, kind: dynWedgeStop, number: 2}, // unrelated stop
		{position: 2 * blicksPerQuarter, kind: dynWedgeStop, number: 1}, // matching stop
		{position: 2 * blicksPerQuarter, kind: dynLevel, loudness: 6},
	}

	points := buildCurve(events, func(e dynEvent) float64 { return e.loudness }, 6)

	// The curve should NOT have a step transition (hold-then-jump) before the f at 2Q.
	// Count points at the step transition position (2Q - stepTransitionBlicks).
	transitionPos := float64(2*blicksPerQuarter) - float64(stepTransitionBlicks)
	hasStepTransition := false
	for k := 0; k < len(points); k += 2 {
		if points[k] == transitionPos {
			hasStepTransition = true
		}
	}
	if hasStepTransition {
		t.Errorf("unexpected step transition at %f; f at 2Q should be detected as wedge target", transitionPos)
	}
}

// TestBuildDynamics_DifferentDivisions tests that dynamics positions are
// correct when a part has different divisions than the first part.
func TestBuildDynamics_DifferentDivisions(t *testing.T) {
	// Part 1 has divisions=4, Part 2 has divisions=8.
	// Dynamics at the half-note boundary should be at the same blick position
	// regardless of divisions.
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list>
    <score-part id="P1"><part-name>S</part-name></score-part>
    <score-part id="P2"><part-name>A</part-name></score-part>
  </part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
  </part>
  <part id="P2">
    <measure>
      <attributes><divisions>8</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>16</duration><type>half</type></note>
      <direction><direction-type><dynamics><f/></dynamics></direction-type></direction>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>16</duration><type>half</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])

	// Build dynamics for part 2 (which has divisions=8)
	events := buildDynamics(score.Part[1], unrolled)

	if len(events) != 1 {
		t.Fatalf("expected 1 dynamic event, got %d", len(events))
	}

	// The f dynamic is after a half note in divisions=8, so duration=16 means 2 quarter notes.
	expectedPos := int64(2 * blicksPerQuarter)
	if events[0].position != expectedPos {
		t.Errorf("dynamics position: expected %d (2Q), got %d (diff = %d)",
			expectedPos, events[0].position, events[0].position-expectedPos)
	}
}

// TestInsertCurvePoints_NoDuplicates tests that insertCurvePoints replaces
// existing points at the same positions rather than creating duplicates.
func TestInsertCurvePoints_NoDuplicates(t *testing.T) {
	// Existing curve: points at positions 0, 100, 200
	points := []float64{0, 1.0, 100, 2.0, 200, 3.0}

	// Insert spike at position 100 (overlapping) to 150 (new)
	result := insertCurvePoints(points, 100, 5.0, 150, 2.5)

	// Check no duplicate positions
	seen := map[float64]bool{}
	for i := 0; i < len(result); i += 2 {
		pos := result[i]
		if seen[pos] {
			t.Errorf("duplicate position %f in curve points: %v", pos, result)
		}
		seen[pos] = true
	}

	// Should have positions: 0, 100(replaced), 150(new), 200
	if len(result) != 8 {
		t.Errorf("expected 8 values (4 points), got %d: %v", len(result), result)
	}

	// The value at position 100 should be the new value (5.0), not the old (2.0)
	for i := 0; i < len(result); i += 2 {
		if result[i] == 100 && result[i+1] != 5.0 {
			t.Errorf("value at position 100: expected 5.0, got %f", result[i+1])
		}
	}
}

// TestInsertCurvePoints_NoOverlap tests insertion with no overlapping positions.
func TestInsertCurvePoints_NoOverlap(t *testing.T) {
	points := []float64{0, 1.0, 200, 3.0}
	result := insertCurvePoints(points, 50, 2.0, 100, 2.5)

	// Should have: 0, 50, 100, 200
	if len(result) != 8 {
		t.Errorf("expected 8 values (4 points), got %d: %v", len(result), result)
	}
}

// TestDynamicsToLevel_Niente tests that niente (<n/>) is matched correctly
// and that unrelated n-prefixed elements don't match.
func TestDynamicsToLevel_Niente(t *testing.T) {
	tests := []struct {
		name     string
		innerXML string
		wantOK   bool
		wantLoud float64
	}{
		{"n self-closing", "<n/>", true, -12},
		{"n with closing tag", "<n></n>", true, -12},
		{"unrelated n-prefix", "<notreal/>", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &musicxml.Dynamics{InnerXML: tt.innerXML}
			lvl, ok := dynamicsToLevel(d)
			if ok != tt.wantOK {
				t.Errorf("dynamicsToLevel(%q): ok = %v, want %v", tt.innerXML, ok, tt.wantOK)
			}
			if ok && lvl.loudness != tt.wantLoud {
				t.Errorf("dynamicsToLevel(%q): loudness = %f, want %f", tt.innerXML, lvl.loudness, tt.wantLoud)
			}
		})
	}
}
