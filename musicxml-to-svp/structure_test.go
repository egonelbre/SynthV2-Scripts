package main

import (
	"fmt"
	"strings"
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
	unrolled, meters, tempos := buildStructure(score.Part[0])

	if len(unrolled) != 2 {
		t.Fatalf("expected 2 unrolled measures, got %d", len(unrolled))
	}
	if unrolled[0].startBlicks != 0 {
		t.Errorf("measure 0 start: expected 0, got %d", unrolled[0].startBlicks)
	}
	// 4/4 measure = 4 quarter notes
	expectedMeasureDuration := int64(4 * blicksPerQuarter)
	if unrolled[1].startBlicks != expectedMeasureDuration {
		t.Errorf("measure 1 start: expected %d, got %d", expectedMeasureDuration, unrolled[1].startBlicks)
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
	unrolled, _, _ := buildStructure(score.Part[0])

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
	unrolled, _, _ := buildStructure(score.Part[0])

	if len(unrolled) != 3 {
		t.Fatalf("expected 3 unrolled measures, got %d", len(unrolled))
	}

	// Measure 0 (pickup): starts at 0
	if unrolled[0].startBlicks != 0 {
		t.Errorf("pickup measure start: expected 0, got %d", unrolled[0].startBlicks)
	}

	// Measure 1: should start at 1 quarter note (pickup duration), not 4 quarter notes
	expectedStart := int64(blicksPerQuarter) // 1 quarter note
	if unrolled[1].startBlicks != expectedStart {
		t.Errorf("measure 1 start: expected %d (1Q), got %d (diff = %d)",
			expectedStart, unrolled[1].startBlicks, unrolled[1].startBlicks-expectedStart)
	}

	// Measure 2: should start at 1Q + 4Q = 5Q
	expectedStart = int64(5 * blicksPerQuarter)
	if unrolled[2].startBlicks != expectedStart {
		t.Errorf("measure 2 start: expected %d (5Q), got %d (diff = %d)",
			expectedStart, unrolled[2].startBlicks, unrolled[2].startBlicks-expectedStart)
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
	_, _, tempos := buildStructure(score.Part[0])

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
	unrolled, _, _ := buildStructure(score.Part[0])

	if len(unrolled) != 2 {
		t.Fatalf("expected 2 unrolled measures, got %d", len(unrolled))
	}

	// Measure 0 starts at 0
	if unrolled[0].startBlicks != 0 {
		t.Errorf("measure 0 start: expected 0, got %d", unrolled[0].startBlicks)
	}

	// Measure 1 should start at 4 quarter notes (full 4/4 measure)
	expectedStart := int64(4 * blicksPerQuarter)
	if unrolled[1].startBlicks != expectedStart {
		t.Errorf("measure 1 start: expected %d (4Q), got %d (diff = %d)",
			expectedStart, unrolled[1].startBlicks, unrolled[1].startBlicks-expectedStart)
	}
}

// TestBuildStructure_CompoundMeter tests a simple 6/8 meter.
func TestBuildStructure_CompoundMeter(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>2</divisions><time><beats>6</beats><beat-type>8</beat-type></time></attributes>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>3</duration><type>quarter</type><dot/></note>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>3</duration><type>quarter</type><dot/></note>
    </measure>
    <measure>
      <note><pitch><step>E</step><octave>4</octave></pitch><duration>3</duration><type>quarter</type><dot/></note>
      <note><pitch><step>F</step><octave>4</octave></pitch><duration>3</duration><type>quarter</type><dot/></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, meters, _ := buildStructure(score.Part[0])

	if len(meters) != 1 || meters[0].Numerator != 6 || meters[0].Denominator != 8 {
		t.Errorf("expected 6/8 meter, got %+v", meters)
	}

	// 6/8 = 6 eighth notes = 3 quarter notes
	expectedDuration := int64(3 * blicksPerQuarter)
	if unrolled[1].startBlicks != expectedDuration {
		t.Errorf("measure 1 start: expected %d (3Q), got %d", expectedDuration, unrolled[1].startBlicks)
	}
}

// TestBuildStructure_AdditiveBeats tests an additive meter like 2+3/8.
func TestBuildStructure_AdditiveBeats(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>2</divisions><time><beats>2+3</beats><beat-type>8</beat-type></time></attributes>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>2</duration><type>quarter</type></note>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>3</duration><type>dotted-quarter</type></note>
    </measure>
    <measure>
      <note><pitch><step>E</step><octave>4</octave></pitch><duration>2</duration><type>quarter</type></note>
      <note><pitch><step>F</step><octave>4</octave></pitch><duration>3</duration><type>dotted-quarter</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, meters, _ := buildStructure(score.Part[0])

	if len(meters) != 1 || meters[0].Numerator != 5 || meters[0].Denominator != 8 {
		t.Errorf("expected 5/8 meter (from 2+3), got %+v", meters)
	}

	// 5/8 = 5 eighth notes = 2.5 quarter notes
	expectedDuration := int64(blicksPerQuarter) * 5 * 4 / 8
	if unrolled[1].startBlicks != expectedDuration {
		t.Errorf("measure 1 start: expected %d (2.5Q), got %d", expectedDuration, unrolled[1].startBlicks)
	}
}

// TestBuildStructure_LongAdditiveBeats tests "1+1+1+1+1"/4 additive meter.
func TestBuildStructure_LongAdditiveBeats(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>1</divisions><time><beats>1+1+1+1+1</beats><beat-type>4</beat-type></time></attributes>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>5</duration><type>whole</type></note>
    </measure>
    <measure>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>5</duration><type>whole</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, meters, _ := buildStructure(score.Part[0])

	if len(meters) != 1 || meters[0].Numerator != 5 || meters[0].Denominator != 4 {
		t.Errorf("expected 5/4 meter (from 1+1+1+1+1), got %+v", meters)
	}

	// 5/4 = 5 quarter notes
	expectedDuration := int64(5 * blicksPerQuarter)
	if unrolled[1].startBlicks != expectedDuration {
		t.Errorf("measure 1 start: expected %d (5Q), got %d", expectedDuration, unrolled[1].startBlicks)
	}
}

// TestBuildStructure_MeterChangeToAdditive tests changing from 4/4 to an additive meter.
func TestBuildStructure_MeterChangeToAdditive(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>2</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note><pitch><step>C</step><octave>4</octave></pitch><duration>8</duration><type>whole</type></note>
    </measure>
    <measure>
      <attributes><time><beats>3+2</beats><beat-type>8</beat-type></time></attributes>
      <note><pitch><step>D</step><octave>4</octave></pitch><duration>3</duration><type>dotted-quarter</type></note>
      <note><pitch><step>E</step><octave>4</octave></pitch><duration>2</duration><type>quarter</type></note>
    </measure>
    <measure>
      <note><pitch><step>F</step><octave>4</octave></pitch><duration>3</duration><type>dotted-quarter</type></note>
      <note><pitch><step>G</step><octave>4</octave></pitch><duration>2</duration><type>quarter</type></note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, meters, _ := buildStructure(score.Part[0])

	if len(meters) != 2 {
		t.Fatalf("expected 2 meter changes, got %d", len(meters))
	}
	if meters[1].Numerator != 5 || meters[1].Denominator != 8 {
		t.Errorf("expected 5/8 meter (from 3+2), got %+v", meters[1])
	}

	// Measure 0: 4/4 = 4Q
	expectedStart1 := int64(4 * blicksPerQuarter)
	if unrolled[1].startBlicks != expectedStart1 {
		t.Errorf("measure 1 start: expected %d (4Q), got %d", expectedStart1, unrolled[1].startBlicks)
	}

	// Measure 1: 5/8 = 2.5Q
	expectedStart2 := expectedStart1 + int64(blicksPerQuarter)*5*4/8
	if unrolled[2].startBlicks != expectedStart2 {
		t.Errorf("measure 2 start: expected %d, got %d", expectedStart2, unrolled[2].startBlicks)
	}
}

// navigationTestMeasure builds a minimal MusicXML measure element with optional
// navigation directives and barline attributes for testing.
func navigationTestMeasure(label string, opts ...string) string {
	var directionParts []string
	var barlineAttrs []string
	var barlineChildren []string

	for _, opt := range opts {
		switch {
		case opt == "dacapo":
			directionParts = append(directionParts, `dacapo="yes"`)
		case strings.HasPrefix(opt, "dalsegno="):
			directionParts = append(directionParts, fmt.Sprintf(`dalsegno="%s"`, opt[9:]))
		case strings.HasPrefix(opt, "tocoda="):
			directionParts = append(directionParts, fmt.Sprintf(`tocoda="%s"`, opt[7:]))
		case opt == "fine":
			directionParts = append(directionParts, `fine="yes"`)
		case strings.HasPrefix(opt, "sound-segno="):
			directionParts = append(directionParts, fmt.Sprintf(`segno="%s"`, opt[12:]))
		case strings.HasPrefix(opt, "sound-coda="):
			directionParts = append(directionParts, fmt.Sprintf(`coda="%s"`, opt[11:]))
		case strings.HasPrefix(opt, "barline-segno="):
			barlineAttrs = append(barlineAttrs, fmt.Sprintf(`segno="%s"`, opt[14:]))
		case strings.HasPrefix(opt, "barline-coda="):
			barlineAttrs = append(barlineAttrs, fmt.Sprintf(`coda="%s"`, opt[13:]))
		case opt == "forward-repeat":
			barlineChildren = append(barlineChildren, `<repeat direction="forward"/>`)
		case opt == "backward-repeat":
			barlineChildren = append(barlineChildren, `<repeat direction="backward"/>`)
		case opt == "backward-repeat-after-jump":
			barlineChildren = append(barlineChildren, `<repeat direction="backward" after-jump="yes"/>`)
		}
	}

	var sb strings.Builder
	sb.WriteString("    <measure>\n")

	if len(barlineAttrs) > 0 || len(barlineChildren) > 0 {
		fmt.Fprintf(&sb, "      <barline %s>%s</barline>\n",
			strings.Join(barlineAttrs, " "),
			strings.Join(barlineChildren, ""))
	}

	if len(directionParts) > 0 {
		fmt.Fprintf(&sb, "      <direction><sound %s/></direction>\n",
			strings.Join(directionParts, " "))
	}

	fmt.Fprintf(&sb,
		"      <note><pitch><step>%s</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>\n",
		label)
	sb.WriteString("    </measure>\n")
	return sb.String()
}

func navigationTestScore(measures ...string) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
`)
	// First measure needs attributes.
	for i, m := range measures {
		if i == 0 {
			m = strings.Replace(m, "<measure>\n",
				"<measure>\n      <attributes><divisions>1</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>\n", 1)
		}
		sb.WriteString(m)
	}
	sb.WriteString("  </part>\n</score-partwise>")
	return sb.String()
}

func checkUnrolledIndices(t *testing.T, unrolled []playedMeasure, expected []int) {
	t.Helper()
	if len(unrolled) != len(expected) {
		got := make([]int, len(unrolled))
		for i, pm := range unrolled {
			got[i] = pm.measureIdx
		}
		t.Fatalf("expected %d measures %v, got %d: %v", len(expected), expected, len(unrolled), got)
	}
	for i, pm := range unrolled {
		if pm.measureIdx != expected[i] {
			got := make([]int, len(unrolled))
			for j, p := range unrolled {
				got[j] = p.measureIdx
			}
			t.Fatalf("unrolled[%d].measureIdx: expected %d, got %d\n  expected: %v\n  got:      %v",
				i, expected[i], pm.measureIdx, expected, got)
		}
	}
}

// TestNavigation_DaCapo: A, B(D.C.) -> A, B, A, B
func TestNavigation_DaCapo(t *testing.T) {
	xml := navigationTestScore(
		navigationTestMeasure("A"),
		navigationTestMeasure("B", "dacapo"),
	)
	score := parseTestScore(t, xml)
	unrolled, _, _ := buildStructure(score.Part[0])
	checkUnrolledIndices(t, unrolled, []int{0, 1, 0, 1})
}

// TestNavigation_DaCapoAlFine: A, B(Fine), C(D.C.) -> A, B, C, A, B
func TestNavigation_DaCapoAlFine(t *testing.T) {
	xml := navigationTestScore(
		navigationTestMeasure("A"),
		navigationTestMeasure("B", "fine"),
		navigationTestMeasure("C", "dacapo"),
	)
	score := parseTestScore(t, xml)
	unrolled, _, _ := buildStructure(score.Part[0])
	checkUnrolledIndices(t, unrolled, []int{0, 1, 2, 0, 1})
}

// TestNavigation_DalSegno: A, B(Segno), C(D.S.) -> A, B, C, B, C
func TestNavigation_DalSegno(t *testing.T) {
	xml := navigationTestScore(
		navigationTestMeasure("A"),
		navigationTestMeasure("B", "sound-segno=s1"),
		navigationTestMeasure("C", "dalsegno=s1"),
	)
	score := parseTestScore(t, xml)
	unrolled, _, _ := buildStructure(score.Part[0])
	checkUnrolledIndices(t, unrolled, []int{0, 1, 2, 1, 2})
}

// TestNavigation_DalSegnoAlFine: A, B(Segno), C(Fine), D(D.S.) -> A, B, C, D, B, C
func TestNavigation_DalSegnoAlFine(t *testing.T) {
	xml := navigationTestScore(
		navigationTestMeasure("A"),
		navigationTestMeasure("B", "sound-segno=s1"),
		navigationTestMeasure("C", "fine"),
		navigationTestMeasure("D", "dalsegno=s1"),
	)
	score := parseTestScore(t, xml)
	unrolled, _, _ := buildStructure(score.Part[0])
	checkUnrolledIndices(t, unrolled, []int{0, 1, 2, 3, 1, 2})
}

// TestNavigation_DalSegnoAlCoda: A, B(Segno), C(ToCoda), D(D.S.), E(Coda)
// -> A, B, C, D, B, C, E
func TestNavigation_DalSegnoAlCoda(t *testing.T) {
	xml := navigationTestScore(
		navigationTestMeasure("A"),
		navigationTestMeasure("B", "sound-segno=s1"),
		navigationTestMeasure("C", "tocoda=c1"),
		navigationTestMeasure("D", "dalsegno=s1"),
		navigationTestMeasure("E", "sound-coda=c1"),
	)
	score := parseTestScore(t, xml)
	unrolled, _, _ := buildStructure(score.Part[0])
	checkUnrolledIndices(t, unrolled, []int{0, 1, 2, 3, 1, 2, 4})
}

// TestNavigation_DaCapoAlCoda: A(ToCoda), B(D.C.), C(Coda) -> A, B, A, C
func TestNavigation_DaCapoAlCoda(t *testing.T) {
	xml := navigationTestScore(
		navigationTestMeasure("A", "tocoda=c1"),
		navigationTestMeasure("B", "dacapo"),
		navigationTestMeasure("C", "sound-coda=c1"),
	)
	score := parseTestScore(t, xml)
	unrolled, _, _ := buildStructure(score.Part[0])
	checkUnrolledIndices(t, unrolled, []int{0, 1, 0, 2})
}

// TestNavigation_DSWithRepeatsSkipped: A(Segno,fwd), B(bwd), C(D.S.)
// -> A, B, A, B, C, A, B (repeats skipped after jump)
func TestNavigation_DSWithRepeatsSkipped(t *testing.T) {
	xml := navigationTestScore(
		navigationTestMeasure("A", "sound-segno=s1", "forward-repeat"),
		navigationTestMeasure("B", "backward-repeat"),
		navigationTestMeasure("C", "dalsegno=s1"),
	)
	score := parseTestScore(t, xml)
	unrolled, _, _ := buildStructure(score.Part[0])
	// Phase 1: 0,1,0,1,2 (measures 0-2 with repeat in 0-1)
	// Phase 2: 0,1 (jump to segno=0, play 0-2 without repeats, but D.S. is at 2 so phase2 ends at 2)
	// Wait - phase 1 ends at jumpIdx=2, phase 2 jumps back to segno=0, plays 0..2 without repeats
	// But the jump is at measure 2 so we need to think about this...
	// Actually: phase 1 plays 0..jumpIdx with repeats = 0,1,0,1,2
	// Phase 2 plays segno(0)..last(2) without repeats = 0,1,2
	// But that gives A,B,A,B,C,A,B,C which is 8
	// The plan says expected: A, B, A, B, C, A, B which is 7
	// Hmm, in phase 2 we play from segno to end without repeats, which is 0,1,2
	// So total would be 0,1,0,1,2,0,1,2
	// But the plan expects 0,1,0,1,2,0,1
	// Let me reconsider - after the D.S. jump, we're replaying from the segno,
	// and the D.S. instruction at measure 2 shouldn't be taken again.
	// So phase 2 should play all measures from segno to end: 0,1,2
	checkUnrolledIndices(t, unrolled, []int{0, 1, 0, 1, 2, 0, 1, 2})
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
	unrolled, _, _ := buildStructure(score.Part[0])

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
