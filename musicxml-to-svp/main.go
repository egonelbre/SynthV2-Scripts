package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/phonemes"
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

	data = ensureUTF8(data)

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

	// Pass 1: structure (derived from first part)
	irScore := &Score{}
	var unrolled []playedMeasure
	if len(score.Part) > 0 {
		unrolled, irScore.Meters, irScore.Tempos = buildStructure(score.Part[0])

		// Warn if other parts have more measures than the first part.
		firstMeasures := len(score.Part[0].Measure)
		for _, part := range score.Part[1:] {
			if len(part.Measure) > firstMeasures {
				fmt.Fprintf(os.Stderr, "warning: part %q has %d measures but structure is derived from first part with %d measures; extra measures will be ignored\n",
					partNames[part.Id], len(part.Measure), firstMeasures)
			}
		}
	}

	// Pass 2+3+4: per part
	for partIdx, part := range score.Part {
		partName := partNames[part.Id]
		if partName == "" {
			partName = fmt.Sprintf("Part %d", partIdx+1)
		}
		p := Part{Name: partName}
		p.Notes = buildNotes(part, unrolled)
		fillLyrics(p.Notes)
		p.Dynamics = buildDynamics(part, unrolled)
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

// ensureUTF8 converts UTF-16 (LE or BE) encoded data to UTF-8.
// If the data is already UTF-8 (or has no BOM), it is returned as-is.
func ensureUTF8(data []byte) []byte {
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		// UTF-16 Little Endian BOM
		return utf16ToUTF8(data[2:], binary.LittleEndian)
	}
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		// UTF-16 Big Endian BOM
		return utf16ToUTF8(data[2:], binary.BigEndian)
	}
	// Strip UTF-8 BOM if present.
	return bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
}

func utf16ToUTF8(data []byte, order binary.ByteOrder) []byte {
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}
	u16 := make([]uint16, len(data)/2)
	for i := range u16 {
		u16[i] = order.Uint16(data[2*i:])
	}
	runes := utf16.Decode(u16)
	var buf bytes.Buffer
	buf.Grow(len(data))
	for _, r := range runes {
		buf.WriteRune(r)
	}
	// Replace the encoding declaration so xml.Unmarshal is happy.
	result := buf.Bytes()
	result = bytes.Replace(result, []byte("encoding='UTF-16'"), []byte("encoding='UTF-8'"), 1)
	result = bytes.Replace(result, []byte(`encoding="UTF-16"`), []byte(`encoding="UTF-8"`), 1)
	return result
}
