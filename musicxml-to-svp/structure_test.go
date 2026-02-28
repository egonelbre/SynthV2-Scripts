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
