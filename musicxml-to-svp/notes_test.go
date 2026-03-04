package main

import (
	"testing"
)

// TestBuildNotes_TieResolution tests that tied notes are merged.
func TestBuildNotes_TieResolution(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="start"/>
        <lyric number="1"><text>la</text></lyric>
      </note>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="stop"/>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 1 {
		t.Fatalf("expected 1 merged note, got %d", len(notes))
	}
	// Half note + half note = whole note duration
	expectedDur := int64(blicksPerQuarter * 4)
	if notes[0].Duration != expectedDur {
		t.Errorf("tied note duration: expected %d, got %d", expectedDur, notes[0].Duration)
	}
	if notes[0].Lyric != "la" {
		t.Errorf("tied note lyric: expected %q, got %q", "la", notes[0].Lyric)
	}
}

// TestBuildNotes_Articulations tests articulation bitmask extraction.
func TestBuildNotes_Articulations(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <notations><articulations><staccato/></articulations></notations>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <notations><articulations><accent/></articulations></notations>
      </note>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <notations><articulations><tenuto/></articulations></notations>
      </note>
      <note>
        <pitch><step>F</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 4 {
		t.Fatalf("expected 4 notes, got %d", len(notes))
	}
	if notes[0].Articulations&ArticulationStaccato == 0 {
		t.Error("note 0: expected staccato")
	}
	if notes[1].Articulations&ArticulationAccent == 0 {
		t.Error("note 1: expected accent")
	}
	if notes[2].Articulations&ArticulationTenuto == 0 {
		t.Error("note 2: expected tenuto")
	}
	if notes[3].Articulations != 0 {
		t.Errorf("note 3: expected no articulations, got %d", notes[3].Articulations)
	}
}

// TestBuildNotes_GraceNoteAttachment tests leading grace note attachment.
func TestBuildNotes_GraceNoteAttachment(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <grace slash="yes"/>
        <pitch><step>D</step><octave>4</octave></pitch>
        <type>eighth</type>
      </note>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>16</duration><type>whole</type>
        <lyric number="1"><text>la</text></lyric>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 1 {
		t.Fatalf("expected 1 note (with leading grace), got %d", len(notes))
	}
	if len(notes[0].LeadingGraces) != 1 {
		t.Fatalf("expected 1 leading grace, got %d", len(notes[0].LeadingGraces))
	}
	g := notes[0].LeadingGraces[0]
	if g.Pitch != 62 { // D4
		t.Errorf("grace pitch: expected 62, got %d", g.Pitch)
	}
	if !g.Acciaccatura {
		t.Error("grace: expected acciaccatura (slash)")
	}
	if g.Duration == 0 {
		t.Error("grace: expected non-zero duration")
	}
	// Main note onset should be shifted by grace duration
	if notes[0].Onset != g.Duration {
		t.Errorf("main note onset: expected %d, got %d", g.Duration, notes[0].Onset)
	}
}

// TestBuildNotes_Transposition tests chromatic transposition.
func TestBuildNotes_Transposition(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>Bb Clarinet</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes>
        <divisions>4</divisions>
        <time><beats>4</beats><beat-type>4</beat-type></time>
        <transpose><chromatic>-2</chromatic></transpose>
      </attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>16</duration><type>whole</type>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	// C4 = MIDI 60, transposed down 2 semitones = 58
	if notes[0].Pitch != 58 {
		t.Errorf("transposed pitch: expected 58, got %d", notes[0].Pitch)
	}
}

// TestBuildNotes_LeadingGraceNoCursorDrift tests that leading grace notes
// don't cause cursor drift for subsequent notes (bug #1).
func TestBuildNotes_LeadingGraceNoCursorDrift(t *testing.T) {
	// Two quarter notes, first has a leading grace. The second note should
	// start at exactly 1 quarter note, not shifted by the grace duration.
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <grace slash="yes"/>
        <pitch><step>D</step><octave>4</octave></pitch>
        <type>16th</type>
      </note>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <lyric number="1"><text>do</text></lyric>
      </note>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <lyric number="1"><text>re</text></lyric>
      </note>
      <note>
        <pitch><step>F</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <lyric number="1"><text>mi</text></lyric>
      </note>
      <note>
        <pitch><step>G</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <lyric number="1"><text>fa</text></lyric>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 4 {
		t.Fatalf("expected 4 notes, got %d", len(notes))
	}

	// Second note (E4) should start at exactly 1 quarter note
	expectedOnset := int64(blicksPerQuarter)
	if notes[1].Onset != expectedOnset {
		t.Errorf("note 1 onset: expected %d, got %d (drift = %d)",
			expectedOnset, notes[1].Onset, notes[1].Onset-expectedOnset)
	}

	// Third note (F4) at 2 quarter notes
	expectedOnset = int64(2 * blicksPerQuarter)
	if notes[2].Onset != expectedOnset {
		t.Errorf("note 2 onset: expected %d, got %d (drift = %d)",
			expectedOnset, notes[2].Onset, notes[2].Onset-expectedOnset)
	}

	// Fourth note (G4) at 3 quarter notes
	expectedOnset = int64(3 * blicksPerQuarter)
	if notes[3].Onset != expectedOnset {
		t.Errorf("note 3 onset: expected %d, got %d (drift = %d)",
			expectedOnset, notes[3].Onset, notes[3].Onset-expectedOnset)
	}
}

// TestBuildNotes_TiedNotesSamePitchTwoVoices tests that ties at the same
// pitch in different voices don't interfere (bug #2).
func TestBuildNotes_TiedNotesSamePitchTwoVoices(t *testing.T) {
	// Two voices both have tied C4 across two measures.
	// Voice 1: C4 half tied to C4 half
	// Voice 2: C4 half tied to C4 half (via backup)
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="start"/>
        <lyric number="1"><text>voice1</text></lyric>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
      </note>
      <backup><duration>16</duration></backup>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="start"/>
        <lyric number="1"><text>voice2</text></lyric>
      </note>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
      </note>
    </measure>
    <measure>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="stop"/>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
      </note>
      <backup><duration>16</duration></backup>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="stop"/>
      </note>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	// Should have 6 notes: voice1-C4(tied), voice1-D4, voice2-C4(tied), voice2-E4, voice1-D4, voice2-E4
	// The two C4 tied notes should each be half+half = whole note duration
	expectedTiedDur := int64(blicksPerQuarter * 4) // half + half

	var tiedC4Notes []Note
	for _, n := range notes {
		if n.Pitch == 60 { // C4
			tiedC4Notes = append(tiedC4Notes, n)
		}
	}

	if len(tiedC4Notes) != 2 {
		t.Fatalf("expected 2 tied C4 notes (one per voice), got %d", len(tiedC4Notes))
	}

	for i, n := range tiedC4Notes {
		if n.Duration != expectedTiedDur {
			t.Errorf("tied C4 note %d: expected duration %d, got %d", i, expectedTiedDur, n.Duration)
		}
	}
}

// TestBuildNotes_SimpleTriplets tests 3 eighth-note triplets filling one quarter note beat.
func TestBuildNotes_SimpleTriplets(t *testing.T) {
	// divisions=12, each triplet eighth has duration=4 (normal eighth=6, scaled by 2/3)
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
        <lyric number="1"><text>la</text></lyric>
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
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 4 {
		t.Fatalf("expected 4 notes, got %d", len(notes))
	}

	// Each triplet eighth = Q/3 blicks
	tripletDur := int64(blicksPerQuarter / 3)
	for i := range 3 {
		expectedOnset := int64(i) * tripletDur
		if notes[i].Onset != expectedOnset {
			t.Errorf("triplet note %d onset: expected %d, got %d", i, expectedOnset, notes[i].Onset)
		}
		if notes[i].Duration != tripletDur {
			t.Errorf("triplet note %d duration: expected %d, got %d", i, tripletDur, notes[i].Duration)
		}
	}

	// Fourth note (regular quarter) starts at 1 quarter note
	expectedOnset := int64(blicksPerQuarter)
	if notes[3].Onset != expectedOnset {
		t.Errorf("quarter note onset: expected %d, got %d", expectedOnset, notes[3].Onset)
	}
	if notes[3].Duration != int64(blicksPerQuarter) {
		t.Errorf("quarter note duration: expected %d, got %d", int64(blicksPerQuarter), notes[3].Duration)
	}
}

// TestBuildNotes_Quintuplets tests 5 notes in the time of 4.
func TestBuildNotes_Quintuplets(t *testing.T) {
	// divisions=20, normal sixteenth=5, quintuplet sixteenth = 5*4/5 = 4
	// 5 quintuplet sixteenths fill one quarter note
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>20</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>4</duration><type>16th</type>
        <time-modification><actual-notes>5</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>4</duration><type>16th</type>
        <time-modification><actual-notes>5</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>4</duration><type>16th</type>
        <time-modification><actual-notes>5</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>F</step><octave>4</octave></pitch>
        <duration>4</duration><type>16th</type>
        <time-modification><actual-notes>5</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>G</step><octave>4</octave></pitch>
        <duration>4</duration><type>16th</type>
        <time-modification><actual-notes>5</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 5 {
		t.Fatalf("expected 5 notes, got %d", len(notes))
	}

	// Each quintuplet = Q/5 blicks
	quintDur := int64(blicksPerQuarter / 5)
	for i := range 5 {
		expectedOnset := int64(i) * quintDur
		if notes[i].Onset != expectedOnset {
			t.Errorf("quintuplet note %d onset: expected %d, got %d", i, expectedOnset, notes[i].Onset)
		}
		if notes[i].Duration != quintDur {
			t.Errorf("quintuplet note %d duration: expected %d, got %d", i, quintDur, notes[i].Duration)
		}
	}
}

// TestBuildNotes_NestedTuplets tests a tuplet inside a tuplet (3-in-2 of 3-in-2).
func TestBuildNotes_NestedTuplets(t *testing.T) {
	// Outer triplet: 3 in the time of 2 quarter notes (each = 2Q/3)
	// Inner triplet on first outer note: 3 in the time of 2/3 of a quarter
	// divisions=18 (quarter = 18)
	// Outer note duration: 18 * 2/3 = 12
	// Inner note duration: 12 * 2/3 = 8
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>18</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>quarter</type>
        <time-modification><actual-notes>9</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>8</duration><type>quarter</type>
        <time-modification><actual-notes>9</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>8</duration><type>quarter</type>
        <time-modification><actual-notes>9</actual-notes><normal-notes>4</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>F</step><octave>4</octave></pitch>
        <duration>12</duration><type>quarter</type>
        <time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>
      </note>
      <note>
        <pitch><step>G</step><octave>4</octave></pitch>
        <duration>12</duration><type>quarter</type>
        <time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 5 {
		t.Fatalf("expected 5 notes, got %d", len(notes))
	}

	// Inner nested notes: each duration = 8/18 * Q = 4Q/9
	innerDur := int64(blicksPerQuarter) * 8 / 18
	for i := range 3 {
		expectedOnset := int64(i) * innerDur
		if notes[i].Onset != expectedOnset {
			t.Errorf("inner note %d onset: expected %d, got %d", i, expectedOnset, notes[i].Onset)
		}
		if notes[i].Duration != innerDur {
			t.Errorf("inner note %d duration: expected %d, got %d", i, innerDur, notes[i].Duration)
		}
	}

	// Outer triplet notes 4,5: each duration = 12/18 * Q = 2Q/3
	outerDur := int64(blicksPerQuarter) * 12 / 18
	// Note index 3 onset = 3 * innerDur
	expectedOnset3 := 3 * innerDur
	if notes[3].Onset != expectedOnset3 {
		t.Errorf("outer note 3 onset: expected %d, got %d", expectedOnset3, notes[3].Onset)
	}
	if notes[3].Duration != outerDur {
		t.Errorf("outer note 3 duration: expected %d, got %d", outerDur, notes[3].Duration)
	}
}

// TestBuildNotes_TupletWithTie tests triplet notes tied across a barline.
func TestBuildNotes_TupletWithTie(t *testing.T) {
	// divisions=12, triplet eighth = duration 4
	// Last triplet tied to first note of next measure (quarter = 12)
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
        <lyric number="1"><text>la</text></lyric>
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
        <tie type="start"/>
      </note>
      <note>
        <pitch><step>F</step><octave>4</octave></pitch>
        <duration>36</duration><type>dotted-half</type>
      </note>
    </measure>
    <measure>
      <note>
        <pitch><step>E</step><octave>4</octave></pitch>
        <duration>12</duration><type>quarter</type>
        <tie type="stop"/>
      </note>
      <note>
        <pitch><step>G</step><octave>4</octave></pitch>
        <duration>36</duration><type>dotted-half</type>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	// Find the tied E4 note
	var tiedNote *Note
	for i := range notes {
		if notes[i].Pitch == 64 { // E4
			tiedNote = &notes[i]
			break
		}
	}
	if tiedNote == nil {
		t.Fatal("expected tied E4 note")
	}

	// Triplet eighth (4/12 Q) + quarter (12/12 Q) = 4/12 + 12/12 = 16/12 Q
	tripletDur := int64(blicksPerQuarter) * 4 / 12
	quarterDur := int64(blicksPerQuarter)
	expectedDur := tripletDur + quarterDur
	if tiedNote.Duration != expectedDur {
		t.Errorf("tied tuplet note duration: expected %d, got %d", expectedDur, tiedNote.Duration)
	}
}

// TestBuildNotes_TupletWithRest tests a triplet group containing a rest.
func TestBuildNotes_TupletWithRest(t *testing.T) {
	// divisions=12, triplet eighth = duration 4
	// Pattern: note, rest, note (triplet group), then a quarter note
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
        <rest/>
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
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	// Should have 3 pitched notes (rest is skipped): C4, E4, F4
	if len(notes) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(notes))
	}

	tripletDur := int64(blicksPerQuarter) / 3

	// C4 at onset 0
	if notes[0].Onset != 0 {
		t.Errorf("note 0 onset: expected 0, got %d", notes[0].Onset)
	}
	if notes[0].Duration != tripletDur {
		t.Errorf("note 0 duration: expected %d, got %d", tripletDur, notes[0].Duration)
	}

	// E4 at onset 2*tripletDur (after rest)
	expectedOnset := 2 * tripletDur
	if notes[1].Onset != expectedOnset {
		t.Errorf("note 1 onset: expected %d, got %d", expectedOnset, notes[1].Onset)
	}

	// F4 at onset blicksPerQuarter
	expectedOnset = int64(blicksPerQuarter)
	if notes[2].Onset != expectedOnset {
		t.Errorf("note 2 onset: expected %d, got %d", expectedOnset, notes[2].Onset)
	}
}

// TestBuildNotes_CompoundMeter68 tests notes in a 6/8 measure.
func TestBuildNotes_CompoundMeter68(t *testing.T) {
	// 6/8 measure with dotted-quarter + dotted-quarter grouping
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>2</divisions><time><beats>6</beats><beat-type>8</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>3</duration><type>quarter</type><dot/>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>3</duration><type>quarter</type><dot/>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}

	// Dotted quarter = 1.5 quarter notes
	dottedQuarterBlicks := int64(blicksPerQuarter) * 3 / 2
	if notes[0].Onset != 0 {
		t.Errorf("note 0 onset: expected 0, got %d", notes[0].Onset)
	}
	if notes[0].Duration != dottedQuarterBlicks {
		t.Errorf("note 0 duration: expected %d, got %d", dottedQuarterBlicks, notes[0].Duration)
	}
	if notes[1].Onset != dottedQuarterBlicks {
		t.Errorf("note 1 onset: expected %d, got %d", dottedQuarterBlicks, notes[1].Onset)
	}
	if notes[1].Duration != dottedQuarterBlicks {
		t.Errorf("note 1 duration: expected %d, got %d", dottedQuarterBlicks, notes[1].Duration)
	}
}

// TestFillLyrics_Melismatic tests that empty lyrics become "-".
func TestFillLyrics_Melismatic(t *testing.T) {
	notes := []Note{
		{Lyric: "la"},
		{Lyric: ""},
		{Lyric: ""},
		{Lyric: "ti"},
	}
	fillLyrics(notes)

	expected := []string{"la", "-", "-", "ti"}
	for i, n := range notes {
		if n.Lyric != expected[i] {
			t.Errorf("note %d lyric: expected %q, got %q", i, expected[i], n.Lyric)
		}
	}
}

// TestFillLyrics_GraceLyricTransfer tests lyric transfer to first grace note.
func TestFillLyrics_GraceLyricTransfer(t *testing.T) {
	notes := []Note{
		{
			Lyric: "la",
			LeadingGraces: []GraceNote{
				{Lyric: "", Pitch: 62},
				{Lyric: "", Pitch: 64},
			},
		},
	}
	fillLyrics(notes)

	if notes[0].LeadingGraces[0].Lyric != "la" {
		t.Errorf("first grace lyric: expected %q, got %q", "la", notes[0].LeadingGraces[0].Lyric)
	}
	if notes[0].LeadingGraces[1].Lyric != "-" {
		t.Errorf("second grace lyric: expected %q, got %q", "-", notes[0].LeadingGraces[1].Lyric)
	}
	// Main note keeps its lyric
	if notes[0].Lyric != "la" {
		t.Errorf("main note lyric: expected %q, got %q", "la", notes[0].Lyric)
	}
}

// TestFillLyrics_GracesWithOwnLyrics tests that grace lyrics are preserved.
func TestFillLyrics_GracesWithOwnLyrics(t *testing.T) {
	notes := []Note{
		{
			Lyric: "main",
			LeadingGraces: []GraceNote{
				{Lyric: "grace1", Pitch: 62},
			},
		},
	}
	fillLyrics(notes)

	if notes[0].LeadingGraces[0].Lyric != "grace1" {
		t.Errorf("grace lyric: expected %q, got %q", "grace1", notes[0].LeadingGraces[0].Lyric)
	}
}

// TestFillLyrics_TrailingGracesWithOwnLyrics tests that trailing grace notes
// with existing lyrics are preserved and not overwritten with "-".
func TestFillLyrics_TrailingGracesWithOwnLyrics(t *testing.T) {
	notes := []Note{
		{
			Lyric: "main",
			TrailingGraces: []GraceNote{
				{Lyric: "trail1", Pitch: 62},
				{Lyric: "", Pitch: 64},
			},
		},
	}
	fillLyrics(notes)

	if notes[0].TrailingGraces[0].Lyric != "trail1" {
		t.Errorf("trailing grace 0 lyric: expected %q, got %q", "trail1", notes[0].TrailingGraces[0].Lyric)
	}
	if notes[0].TrailingGraces[1].Lyric != "-" {
		t.Errorf("trailing grace 1 lyric: expected %q, got %q", "-", notes[0].TrailingGraces[1].Lyric)
	}
}

// TestBuildNotes_TieStopWithGraceUsesFullDuration tests that a tie-stop note
// with leading grace notes adds the full notated duration, not the grace-reduced one.
func TestBuildNotes_TieStopWithGraceUsesFullDuration(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="start"/>
        <lyric number="1"><text>la</text></lyric>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>8</duration><type>half</type>
      </note>
    </measure>
    <measure>
      <note>
        <grace slash="yes"/>
        <pitch><step>E</step><octave>4</octave></pitch>
        <type>eighth</type>
      </note>
      <note>
        <pitch><step>C</step><octave>4</octave></pitch>
        <duration>16</duration><type>whole</type>
        <tie type="stop"/>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	// Should have 2 notes: the tied C4 and the D4.
	// The grace note before the tie-stop should not reduce the tied duration.
	var tiedNote *Note
	for i := range notes {
		if notes[i].Pitch == 60 { // C4
			tiedNote = &notes[i]
			break
		}
	}
	if tiedNote == nil {
		t.Fatal("expected tied C4 note")
	}

	// Half note (8 divs) + whole note (16 divs) = 6 quarter notes
	expectedDur := int64(blicksPerQuarter * 6)
	if tiedNote.Duration != expectedDur {
		t.Errorf("tied note duration: expected %d (6Q), got %d (diff = %d)",
			expectedDur, tiedNote.Duration, tiedNote.Duration-expectedDur)
	}
}

// TestBuildNotes_SlideExtraction tests that slide start/stop pairs compute SlideDelta.
func TestBuildNotes_SlideExtraction(t *testing.T) {
	// A4 (MIDI 69) slides to G4 (MIDI 67) = -200 cents
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>A</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <notations><slide type="start" number="1"/></notations>
        <lyric number="1"><text>la</text></lyric>
      </note>
      <note>
        <pitch><step>G</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <notations><slide type="stop" number="1"/></notations>
        <lyric number="1"><text>ti</text></lyric>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
	if notes[0].SlideDelta != -200 {
		t.Errorf("note 0 SlideDelta: expected -200, got %d", notes[0].SlideDelta)
	}
	if notes[1].SlideDelta != 0 {
		t.Errorf("note 1 SlideDelta: expected 0, got %d", notes[1].SlideDelta)
	}
}

// TestBuildNotes_SlideAcrossMeasures tests slide start in one measure, stop in the next.
func TestBuildNotes_SlideAcrossMeasures(t *testing.T) {
	// C5 (MIDI 72) slides to E5 (MIDI 76) = +400 cents
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>5</octave></pitch>
        <duration>16</duration><type>whole</type>
        <notations><slide type="start" number="1"/></notations>
        <lyric number="1"><text>la</text></lyric>
      </note>
    </measure>
    <measure>
      <note>
        <pitch><step>E</step><octave>5</octave></pitch>
        <duration>16</duration><type>whole</type>
        <notations><slide type="stop" number="1"/></notations>
        <lyric number="1"><text>ti</text></lyric>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
	if notes[0].SlideDelta != 400 {
		t.Errorf("note 0 SlideDelta: expected 400, got %d", notes[0].SlideDelta)
	}
	if notes[1].SlideDelta != 0 {
		t.Errorf("note 1 SlideDelta: expected 0, got %d", notes[1].SlideDelta)
	}
}

// TestBuildNotes_SlideOnTieStop tests slide starting on a tie-stop note.
func TestBuildNotes_SlideOnTieStop(t *testing.T) {
	// C5 half tied to C5 quarter (with slide start) → D4 quarter (slide stop)
	// The tied C5 should have SlideDelta = (62-72)*100 = -1000 cents
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<score-partwise>
  <part-list><score-part id="P1"><part-name>S</part-name></score-part></part-list>
  <part id="P1">
    <measure>
      <attributes><divisions>4</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>
      <note>
        <pitch><step>C</step><octave>5</octave></pitch>
        <duration>8</duration><type>half</type>
        <tie type="start"/>
        <notations><tied type="start"/></notations>
      </note>
      <note>
        <pitch><step>C</step><octave>5</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <tie type="stop"/>
        <notations>
          <tied type="stop"/>
          <slide type="start" number="1"/>
        </notations>
      </note>
      <note>
        <pitch><step>D</step><octave>4</octave></pitch>
        <duration>4</duration><type>quarter</type>
        <notations><slide type="stop" number="1"/></notations>
      </note>
    </measure>
  </part>
</score-partwise>`

	score := parseTestScore(t, xmlData)
	unrolled, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled)

	if len(notes) != 2 {
		t.Fatalf("expected 2 notes (tied C5 + D4), got %d", len(notes))
	}
	// Tied C5 (half + quarter = 3Q) should slide to D4: (62-72)*100 = -1000
	if notes[0].SlideDelta != -1000 {
		t.Errorf("tied note SlideDelta: expected -1000, got %d", notes[0].SlideDelta)
	}
	expectedDur := int64(blicksPerQuarter * 3)
	if notes[0].Duration != expectedDur {
		t.Errorf("tied note duration: expected %d, got %d", expectedDur, notes[0].Duration)
	}
}
