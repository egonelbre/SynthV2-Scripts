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
	unrolled, infos, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled, infos)

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
	unrolled, infos, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled, infos)

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
	unrolled, infos, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled, infos)

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
	unrolled, infos, _, _ := buildStructure(score.Part[0])
	notes := buildNotes(score.Part[0], unrolled, infos)

	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	// C4 = MIDI 60, transposed down 2 semitones = 58
	if notes[0].Pitch != 58 {
		t.Errorf("transposed pitch: expected 58, got %d", notes[0].Pitch)
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
