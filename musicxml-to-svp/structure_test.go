package main

import (
	"testing"
)

// TestBuildStructure_SimpleMetersTempos tests basic meter/tempo extraction.
func TestBuildStructure_SimpleMetersTempos(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>Soprano</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <direction><direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>120</per-minute></metronome></direction-type></direction>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
    <measure>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, infos, meters, tempos := buildStructure(score.Part[0])

	if len(unrolled) != 2 {
		t.Fatalf("expected 2 unrolled measures, got %d", len(unrolled))
	}
	if len(infos) != 2 {
		t.Fatalf("expected 2 measure infos, got %d", len(infos))
	}
	if infos[0].startBlicks != 0 {
		t.Errorf("measure 0 start: expected 0, got %d", infos[0].startBlicks)
	}
	// 4/4 measure = 4 quarter notes
	expectedMeasureDuration := int64(4 * blicksPerQuarter)
	if infos[1].startBlicks != expectedMeasureDuration {
		t.Errorf("measure 1 start: expected %d, got %d", expectedMeasureDuration, infos[1].startBlicks)
	}

	if len(meters) != 1 || meters[0].Numerator != 4 || meters[0].Denominator != 4 {
		t.Errorf("expected 4/4 meter, got %+v", meters)
	}
	if len(tempos) != 1 || tempos[0].BPM != 120 {
		t.Errorf("expected 120 BPM, got %+v", tempos)
	}
}

// TestBuildStructure_RepeatUnrolling tests that repeats are unrolled correctly.
func TestBuildStructure_RepeatUnrolling(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>1</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <barline><repeat direction="forward"/></barline>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>
    </measure>
    <measure>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>
      <barline><repeat direction="backward"/></barline>
    </measure>
    <measure>
      <note><pitch><step>E</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _, _ := buildStructure(score.Part[0])

	// Should be: m0, m1 (pass 1), m0, m1 (pass 2), m2
	if len(unrolled) != 5 {
		t.Fatalf("expected 5 unrolled measures, got %d", len(unrolled))
	}
	expectedIdxs := []int{0, 1, 0, 1, 2}
	expectedVerses := []int{1, 1, 2, 2, 1}
	for i, pm := range unrolled {
		if pm.measureIdx != expectedIdxs[i] {
			t.Errorf("unrolled[%d].measureIdx: expected %d, got %d", i, expectedIdxs[i], pm.measureIdx)
		}
		if pm.verse != expectedVerses[i] {
			t.Errorf("unrolled[%d].verse: expected %d, got %d", i, expectedVerses[i], pm.verse)
		}
	}
}

// TestBuildStructure_PickupMeasure tests that an anacrusis (pickup measure)
// doesn't cause timing drift for subsequent measures (bug #3).
func TestBuildStructure_PickupMeasure(t *testing.T) {
	// First measure is a pickup with just 1 quarter note in 4/4 time.
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note><pitch><step>G</step><octave>4</octave></pitch><duration>4</duration><type>quarter</type></note>
    </measure>
    <measure>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
    <measure>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	_, infos, _, _ := buildStructure(score.Part[0])

	if len(infos) != 3 {
		t.Fatalf("expected 3 measure infos, got %d", len(infos))
	}

	// Measure 0 (pickup): starts at 0
	if infos[0].startBlicks != 0 {
		t.Errorf("pickup measure start: expected 0, got %d", infos[0].startBlicks)
	}

	// Measure 1: should start at 1 quarter note (pickup duration), not 4 quarter notes
	expectedStart := int64(blicksPerQuarter) // 1 quarter note
	if infos[1].startBlicks != expectedStart {
		t.Errorf("measure 1 start: expected %d (1Q), got %d (diff = %d)",
			expectedStart, infos[1].startBlicks, infos[1].startBlicks-expectedStart)
	}

	// Measure 2: should start at 1Q + 4Q = 5Q
	expectedStart = int64(5 * blicksPerQuarter)
	if infos[2].startBlicks != expectedStart {
		t.Errorf("measure 2 start: expected %d (5Q), got %d (diff = %d)",
			expectedStart, infos[2].startBlicks, infos[2].startBlicks-expectedStart)
	}
}

// TestBuildStructure_MetronomeWithSoundElement tests that a Metronome mark
// is not skipped when the Direction has a Sound element for dynamics (bug #5).
func TestBuildStructure_MetronomeWithSoundElement(t *testing.T) {
	// First measure has tempo via Sound. Second measure has a Direction with
	// Sound (for dynamics, no tempo) AND a Metronome mark.
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <direction><sound tempo="120"/></direction>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
    <measure>
      <direction>
        <direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>90</per-minute></metronome></direction-type>
        <sound dynamics="80"/>
      </direction>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>16</duration><type>whole</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	_, _, _, tempos := buildStructure(score.Part[0])

	if len(tempos) < 2 {
		t.Fatalf("expected at least 2 tempo changes, got %d: %+v", len(tempos), tempos)
	}

	// Second tempo should be 90 BPM from the Metronome mark
	if tempos[1].BPM != 90 {
		t.Errorf("second tempo: expected 90 BPM, got %f", tempos[1].BPM)
	}
}

// TestBuildStructure_TupletMeasureTiming tests that a measure containing tuplets
// computes the correct duration (cursor advances properly).
func TestBuildStructure_TupletMeasureTiming(t *testing.T) {
	// Measure 1: triplet eighths filling one quarter + 3 regular quarters = 4/4
	// Measure 2: regular whole note
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>12</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>4</duration><type>eighth</type>
        <time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>4</duration><type>eighth</type>
        <time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>4</duration><type>eighth</type>
        <time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>F</step><octave>4</octave></pitch>
        <duration>12</duration><type>quarter</type>
      </note>
      <note>
        <pitch><step>G</step><octave>4</octave></pitch>
        <duration>12</duration><type>quarter</type>
      </note>
      <note>
        <pitch><step>A</step><octave>4</octave></pitch>
        <duration>12</duration><type>quarter</type>
      </note>
    </measure>
    <measure>
      <note>
        <pitch><step>C</step><octave>5</octave></pitch>
        <duration>48</duration><type>whole</type>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	_, infos, _, _ := buildStructure(score.Part[0])

	if len(infos) != 2 {
		t.Fatalf("expected 2 measure infos, got %d", len(infos))
	}

	// Measure 0 starts at 0
	if infos[0].startBlicks != 0 {
		t.Errorf("measure 0 start: expected 0, got %d", infos[0].startBlicks)
	}

	// Measure 1 should start at 4 quarter notes (full 4/4 measure)
	expectedStart := int64(4 * blicksPerQuarter)
	if infos[1].startBlicks != expectedStart {
		t.Errorf("measure 1 start: expected %d (4Q), got %d (diff = %d)",
			expectedStart, infos[1].startBlicks, infos[1].startBlicks-expectedStart)
	}
}

// TestBuildStructure_NestedRepeats tests that nested repeats are unrolled correctly.
func TestBuildStructure_NestedRepeats(t *testing.T) {
	// Outer repeat: measures 0-3 (2x)
	// Inner repeat: measures 1-2 (2x)
	// Expected: A, B, C, B, C, D, A, B, C, B, C, D
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>1</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <barline><repeat direction="forward"/></barline>
      <note><pitch><step>A</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>
    </measure>
    <measure>
      <barline><repeat direction="forward"/></barline>
      <note><pitch><step>B</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>
    </measure>
    <measure>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>
      <barline><repeat direction="backward"/></barline>
    </measure>
    <measure>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>
      <barline><repeat direction="backward"/></barline>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _, _ := buildStructure(score.Part[0])

	// Expected: A(0), B(1), C(2), B(1), C(2), D(3), A(0), B(1), C(2), B(1), C(2), D(3)
	expectedIdxs := []int{0, 1, 2, 1, 2, 3, 0, 1, 2, 1, 2, 3}
	if len(unrolled) != len(expectedIdxs) {
		t.Fatalf("expected %d unrolled measures, got %d", len(expectedIdxs), len(unrolled))
	}
	for i, pm := range unrolled {
		if pm.measureIdx != expectedIdxs[i] {
			t.Errorf("unrolled[%d].measureIdx: expected %d, got %d", i, expectedIdxs[i], pm.measureIdx)
		}
	}
}
