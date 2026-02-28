package main

import (
	"encoding/xml"
	"testing"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

func parseTestScore(t *testing.T, xmlData string) musicxml.ScorePartwise {
	t.Helper()
	var score musicxml.ScorePartwise
	if err := xml.Unmarshal([]byte(xmlData), &score); err != nil {
		t.Fatalf("failed to parse MusicXML: %v", err)
	}
	return score
}

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
	unrolled, infos, _, _ := buildStructure(score.Part[0])
	events := buildDynamics(score.Part[0], unrolled, infos)

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
