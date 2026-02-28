package main

import (
	"testing"
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
