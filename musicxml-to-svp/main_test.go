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
