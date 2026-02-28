package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/egonelbre/synthv2-scripts/internal/phonemes"
	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

func main() {
	voiceFlag := flag.String("voice", "", "assign voices: choir1, choir2, choir3, or soloists")
	panFlag := flag.String("pan", "default", "panning scheme: default, spread, center")
	langFlag := flag.String("lang", "", "convert lyrics to phonemes: estonian, karelian")
	relaxedFlag := flag.Bool("relaxed", false, "enable relaxed consonant pronunciation")
	outputFlag := flag.String("o", "", "output file path (default: input with .svp extension)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: musicxml-to-svp [flags] <input.musicxml>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := flag.Arg(0)
	outputPath := *outputFlag
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = inputPath[:len(inputPath)-len(ext)] + ".svp"
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}

	var score musicxml.ScorePartwise
	if err := xml.Unmarshal(data, &score); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing MusicXML: %v\n", err)
		os.Exit(1)
	}

	// Build part name lookup
	partNames := map[string]string{}
	for _, sp := range score.PartList.GroupScorePart.ScorePart {
		name := ""
		if sp.PartName != nil {
			name = sp.PartName.EnclosedText
		}
		partNames[sp.Id] = name
	}

	// Pass 1: structure
	irScore := &Score{}
	var unrolled []playedMeasure
	var infos []measureInfo
	if len(score.Part) > 0 {
		unrolled, infos, irScore.Meters, irScore.Tempos = buildStructure(score.Part[0])
	}

	// Pass 2+3+4: per part
	for partIdx, part := range score.Part {
		partName := partNames[part.Id]
		if partName == "" {
			partName = fmt.Sprintf("Part %d", partIdx+1)
		}
		p := Part{Name: partName}
		p.Notes = buildNotes(part, unrolled, infos)
		fillLyrics(p.Notes)
		p.Dynamics = buildDynamics(part, unrolled, infos)
		irScore.Parts = append(irScore.Parts, p)
	}

	// Convert to SVP
	project := scoreToSVP(irScore)

	// Assign voices if requested.
	if *voiceFlag != "" {
		assignVoices(project.Tracks, strings.ToLower(*voiceFlag), *relaxedFlag, *panFlag)
		setNoteAttributes(project.Library)
	}

	// Convert lyrics to phonemes if requested.
	if *langFlag != "" {
		conv := phonemes.New(*langFlag)
		if conv == nil {
			fmt.Fprintf(os.Stderr, "unknown language: %q (options: estonian, karelian)\n", *langFlag)
			os.Exit(1)
		}
		applyPhonemes(project.Library, conv)
	}

	out, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %s (%d parts, %d meters, %d tempos)\n", outputPath, len(project.Tracks), len(project.Time.Meters), len(project.Time.Tempos))
	for _, t := range project.Tracks {
		for _, g := range project.Library {
			if g.UUID == t.Groups[0].GroupID {
				fmt.Fprintf(os.Stderr, "  %s: %d notes\n", t.Name, len(g.Notes))
			}
		}
	}
}
