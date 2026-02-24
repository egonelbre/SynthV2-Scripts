// Command svp-voice assigns voice databases to tracks in SynthV .svp project files
// based on track names (e.g., "Soprano 1", "Tenor 2", "Bass").
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// VoicePart represents a choir voice part.
type VoicePart string

const (
	Soprano      VoicePart = "S"
	MezzoSoprano VoicePart = "MS"
	Alto         VoicePart = "A"
	Tenor        VoicePart = "T"
	Baritone     VoicePart = "Bar"
	Bass         VoicePart = "B"
	Unknown      VoicePart = ""
)

func (p VoicePart) String() string {
	switch p {
	case Soprano:
		return "Soprano"
	case MezzoSoprano:
		return "Mezzo-Soprano"
	case Alto:
		return "Alto"
	case Tenor:
		return "Tenor"
	case Baritone:
		return "Baritone"
	case Bass:
		return "Bass"
	default:
		return "Unknown"
	}
}

// ChoirInfo holds the database configuration for a choir voice collection.
type ChoirInfo struct {
	Name     string
	Language string
	Phoneset string
	// Parts available in this choir. Choir #2 has Mezzo-Soprano and Baritone
	// instead of Alto and Bass.
	Parts []VoicePart
}

var choirs = []ChoirInfo{
	{
		Name:     "Choir Voices #1",
		Language: "english",
		Phoneset: "arpabet",
		Parts:    []VoicePart{Soprano, Alto, Tenor, Bass},
	},
	{
		Name:     "Choir Voices #2",
		Language: "mandarin",
		Phoneset: "xsampa",
		Parts:    []VoicePart{Soprano, MezzoSoprano, Tenor, Baritone},
	},
	{
		Name:     "Choir Voices #3",
		Language: "japanese",
		Phoneset: "romaji",
		Parts:    []VoicePart{Soprano, Alto, Tenor, Bass},
	},
}

// SoloistDB holds a soloist voice database configuration.
type SoloistDB struct {
	Name     string
	Language string
	Phoneset string
	// Preferred track numbers (1-indexed) within the voice part, tried in
	// order. E.g. []int{2, 1} means: prefer "Bass 2", then "Bass 1".
	// A single unnumbered track (e.g. just "Bass") counts as track 1.
	Preferred []int
}

// Ranges from the installed voice databases:
//
//	SOLARIA II:  A#3-A4  (soprano, EN)
//	Sheena 2:    A#3-A#4 (soprano, EN/JA)
//	Mai 2:       B3-B4   (soprano, JA)
//	Ayame 2:     B3-A#4  (soprano, JA)
//	Felicia 2:   A#3-B4  (mezzo/soprano, EN)
//	Natalie 2:   G3-A4   (mezzo/alto, EN)
//	Eri 2:       G#3-E4  (alto, JA)
//	Wei Shu 2:   D3-D4   (alto, ZH/YUE)
//	Kevin 2:     E3-F4   (tenor, EN)
//	Hayden 2:    G3-F#4  (tenor, EN)
//	Ninezero 2:  G#3-G#4 (tenor, EN)
//	Liam:        D3-E4   (baritone, EN)
//	ASTERIAN II: F2-G#3  (bass, EN)
//	Ritchy 2:    G#2-A#3 (bass, EN)
var soloists = map[VoicePart][]SoloistDB{
	Soprano: {
		{Name: "SOLARIA II", Language: "english", Phoneset: "arpabet"},
		{Name: "Sheena 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Mai 2", Language: "japanese", Phoneset: "romaji"},
		{Name: "Ayame 2", Language: "japanese", Phoneset: "romaji"},
		{Name: "Felicia 2", Language: "english", Phoneset: "arpabet"},
	},
	MezzoSoprano: {
		{Name: "Felicia 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Natalie 2", Language: "english", Phoneset: "arpabet"},
		{Name: "SOLARIA II", Language: "english", Phoneset: "arpabet"},
		{Name: "Sheena 2", Language: "english", Phoneset: "arpabet"},
	},
	Alto: {
		{Name: "Natalie 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Felicia 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Eri 2", Language: "japanese", Phoneset: "romaji"},
		{Name: "Wei Shu 2", Language: "cantonese", Phoneset: "xsampa"},
	},
	Tenor: {
		{Name: "Kevin 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Hayden 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Ninezero 2", Language: "english", Phoneset: "arpabet"},
	},
	Baritone: {
		{Name: "Liam", Language: "english", Phoneset: "arpabet"},
		{Name: "Kevin 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Hayden 2", Language: "english", Phoneset: "arpabet"},
	},
	Bass: {
		{Name: "ASTERIAN II", Language: "english", Phoneset: "arpabet", Preferred: []int{3, 2, 1}},
		{Name: "Ritchy 2", Language: "english", Phoneset: "arpabet"},
		{Name: "Liam", Language: "english", Phoneset: "arpabet"},
	},
}

type trackInfo struct {
	Name    string
	Part    VoicePart
	PartNum int // Track number within part (1 for "Bass", 1 for "Bass 1", 2 for "Bass 2", etc.)
	DB      string
}

func main() {
	octaveFix := flag.Bool("octave-fix", false, "transpose tenor and bass voices an octave higher (+12 semitones)")
	output := flag.String("o", "", "output file path (default: overwrite input)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: svp-voice [flags] <file.svp> <choir1|choir2|choir3|soloists>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	path := flag.Arg(0)
	voiceArg := strings.ToLower(flag.Arg(1))
	outPath := path
	if *output != "" {
		outPath = *output
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var project map[string]any
	if err := json.Unmarshal(data, &project); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing SVP file: %v\n", err)
		os.Exit(1)
	}

	tracks, ok := project["tracks"].([]any)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: no tracks found in project\n")
		os.Exit(1)
	}

	var infos []trackInfo
	for _, t := range tracks {
		track := t.(map[string]any)
		name := track["name"].(string)
		part := parseVoicePart(name)

		db := ""
		if mainRef, ok := track["mainRef"].(map[string]any); ok {
			if database, ok := mainRef["database"].(map[string]any); ok {
				db = fmt.Sprintf("%v", database["name"])
			}
		}

		infos = append(infos, trackInfo{Name: name, Part: part, PartNum: parsePartNum(name), DB: db})
	}

	// Display tracks
	fmt.Println("Tracks:")
	for i, info := range infos {
		partStr := info.Part.String()
		if info.Part == Unknown {
			partStr = "???"
		}
		fmt.Printf("  %d. %-20s  [%s]  (current: %s)\n", i+1, info.Name, partStr, info.DB)
	}
	fmt.Println()

	// Check for unknown parts
	for _, info := range infos {
		if info.Part == Unknown {
			fmt.Fprintf(os.Stderr, "Warning: cannot determine voice part for track %q\n", info.Name)
		}
	}

	switch voiceArg {
	case "choir1", "1":
		applyChoir(tracks, infos, choirs[0])
	case "choir2", "2":
		applyChoir(tracks, infos, choirs[1])
	case "choir3", "3":
		applyChoir(tracks, infos, choirs[2])
	case "soloists", "solo", "4":
		applySoloists(tracks, infos)
	default:
		fmt.Fprintf(os.Stderr, "Unknown voice source: %q\n", voiceArg)
		fmt.Fprintf(os.Stderr, "Options: choir1, choir2, choir3, soloists\n")
		os.Exit(1)
	}

	// Transpose tenor/bass an octave up if requested.
	if *octaveFix {
		applyOctaveFix(project, infos)
	}

	// Set note attributes on all notes in the project.
	setNoteAttributes(project)

	// Write backup when overwriting the input file.
	if outPath == path {
		backupPath := path + ".bak"
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing backup: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Backup saved to %s\n", backupPath)
	}

	// Write modified file.
	out, err := json.MarshalIndent(project, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Written %s\n", outPath)
}

func parseVoicePart(name string) VoicePart {
	lower := strings.ToLower(strings.TrimSpace(name))

	// Check for mezzo-soprano first (multi-word).
	if strings.HasPrefix(lower, "mezzo") {
		return MezzoSoprano
	}

	// Take the first word.
	first := strings.Fields(lower)[0]

	switch first {
	case "soprano", "sop":
		return Soprano
	case "alto":
		return Alto
	case "tenor":
		return Tenor
	case "baritone", "bariton":
		return Baritone
	case "bass":
		return Bass
	default:
		return Unknown
	}
}

// parsePartNum extracts the track number from a name like "Bass 2" → 2, "Bass" → 1.
func parsePartNum(name string) int {
	fields := strings.Fields(strings.TrimSpace(name))
	if len(fields) >= 2 {
		// Try to parse the second field as a number (handles "Bass 2", "Soprano 1a" → 1).
		numStr := fields[1]
		// Strip trailing non-digit characters (e.g. "1a" → "1").
		for len(numStr) > 0 && (numStr[len(numStr)-1] < '0' || numStr[len(numStr)-1] > '9') {
			numStr = numStr[:len(numStr)-1]
		}
		if n, err := strconv.Atoi(numStr); err == nil && n > 0 {
			return n
		}
	}
	return 1
}

// mapPartToChoir maps a detected voice part to the closest choir part available.
func mapPartToChoir(part VoicePart, choir ChoirInfo) VoicePart {
	// Check if the part is directly available.
	for _, p := range choir.Parts {
		if p == part {
			return part
		}
	}

	// Fallback mappings.
	switch part {
	case MezzoSoprano:
		// Choir #1/#3 don't have mezzo, use alto.
		return Alto
	case Alto:
		// Choir #2 doesn't have alto, use mezzo.
		return MezzoSoprano
	case Baritone:
		// Choir #1/#3 don't have baritone, use bass.
		return Bass
	case Bass:
		// Choir #2 doesn't have bass, use baritone.
		return Baritone
	default:
		return Soprano
	}
}

func applyChoir(tracks []any, infos []trackInfo, choir ChoirInfo) {
	fmt.Printf("\nApplying %s to all tracks...\n", choir.Name)

	for i, t := range tracks {
		track := t.(map[string]any)
		info := infos[i]

		if info.Part == Unknown {
			fmt.Printf("  Skipping %q (unknown voice part)\n", info.Name)
			continue
		}

		choirPart := mapPartToChoir(info.Part, choir)
		fmt.Printf("  %s → %s (%s)\n", info.Name, choir.Name, choirPart.String())

		// Update mainRef.database
		mainRef := track["mainRef"].(map[string]any)
		database := mainRef["database"].(map[string]any)
		database["name"] = choir.Name
		database["language"] = choir.Language
		database["phoneset"] = choir.Phoneset
		database["backendType"] = "SVR3"
		database["version"] = "202"
		database["languageOverride"] = ""
		database["phonesetOverride"] = ""

		// Update mainRef.voice for choir
		voice := mainRef["voice"].(map[string]any)
		voice["vocalModeInherited"] = true
		voice["vocalModePreset"] = ""
		if _, ok := voice["vocalModeParams"]; !ok {
			voice["vocalModeParams"] = map[string]any{}
		}
		// Set choir separation, preserve existing if present.
		if _, ok := voice["choirSeatingSeparation"]; !ok {
			voice["choirSeatingSeparation"] = 0.7
		}
		// Remove soloist fields that don't belong on choir mainRef.
		delete(voice, "choirNumStems")
		// choirPartName is not set on mainRef (Soprano is default).
		delete(voice, "choirPartName")

		// Update groups
		groups, _ := track["groups"].([]any)
		for _, g := range groups {
			group := g.(map[string]any)
			gVoice := group["voice"].(map[string]any)

			// Set choir part name (omit for Soprano as it's the default).
			if choirPart != Soprano {
				gVoice["choirPartName"] = string(choirPart)
			} else {
				delete(gVoice, "choirPartName")
			}

			// Set/preserve choir stems.
			if _, ok := gVoice["choirNumStems"]; !ok {
				gVoice["choirNumStems"] = float64(4)
			}
			if _, ok := gVoice["choirSeatingSeparation"]; !ok {
				gVoice["choirSeatingSeparation"] = 0.7
			}
		}
	}
}

func applySoloists(tracks []any, infos []trackInfo) {
	// Build assignment map: track index → soloist.
	assignments := make(map[int]SoloistDB)

	// For each voice part, resolve assignments.
	// Group track indices by part.
	partTracks := map[VoicePart][]int{}
	for i, info := range infos {
		if info.Part != Unknown {
			partTracks[info.Part] = append(partTracks[info.Part], i)
		}
	}

	for part, trackIdxs := range partTracks {
		pool := soloists[part]
		if len(pool) == 0 {
			continue
		}

		assigned := map[int]bool{}   // track indices already assigned
		usedVoice := map[int]bool{} // pool indices already used

		// First pass: assign voices that have Preferred set.
		for vi, voice := range pool {
			if len(voice.Preferred) == 0 {
				continue
			}
			for _, prefNum := range voice.Preferred {
				// Find a track with this part number that isn't assigned yet.
				for _, ti := range trackIdxs {
					if !assigned[ti] && infos[ti].PartNum == prefNum {
						assignments[ti] = voice
						assigned[ti] = true
						usedVoice[vi] = true
						break
					}
				}
				if usedVoice[vi] {
					break
				}
			}
		}

		// Second pass: assign remaining tracks from unused voices in pool order.
		vi := 0
		for _, ti := range trackIdxs {
			if assigned[ti] {
				continue
			}
			// Skip used voices.
			for vi < len(pool) && usedVoice[vi] {
				vi++
			}
			if vi >= len(pool) {
				// Wrap around if we run out.
				vi = 0
			}
			assignments[ti] = pool[vi]
			usedVoice[vi] = true
			vi++
		}
	}

	fmt.Println("\nAssigning soloists to tracks...")
	fmt.Println()

	for i, t := range tracks {
		track := t.(map[string]any)
		info := infos[i]

		soloist, ok := assignments[i]
		if !ok {
			if info.Part == Unknown {
				fmt.Printf("  Skipping %q (unknown voice part)\n", info.Name)
			} else {
				fmt.Printf("  Skipping %q (no soloists for %s)\n", info.Name, info.Part.String())
			}
			continue
		}

		fmt.Printf("  %s → %s\n", info.Name, soloist.Name)

		// Update mainRef.database
		mainRef := track["mainRef"].(map[string]any)
		database := mainRef["database"].(map[string]any)
		database["name"] = soloist.Name
		database["language"] = soloist.Language
		database["phoneset"] = soloist.Phoneset
		database["backendType"] = "SVR3"
		database["version"] = "201"
		database["languageOverride"] = ""
		database["phonesetOverride"] = ""

		// Update mainRef.voice for soloist (remove choir fields).
		voice := mainRef["voice"].(map[string]any)
		voice["vocalModeInherited"] = true
		voice["vocalModePreset"] = ""
		if _, ok := voice["vocalModeParams"]; !ok {
			voice["vocalModeParams"] = map[string]any{}
		}
		delete(voice, "choirSeatingSeparation")
		delete(voice, "choirNumStems")
		delete(voice, "choirPartName")

		// Update groups - remove choir fields.
		groups, _ := track["groups"].([]any)
		for _, g := range groups {
			group := g.(map[string]any)
			gVoice := group["voice"].(map[string]any)
			delete(gVoice, "choirPartName")
			delete(gVoice, "choirNumStems")
			delete(gVoice, "choirSeatingSeparation")
		}
	}
}

// applyOctaveFix transposes notes in tenor, baritone, and bass tracks up by 12 semitones.
func applyOctaveFix(project map[string]any, infos []trackInfo) {
	// Build set of library group UUIDs that belong to tenor/baritone/bass tracks.
	transposeGroups := map[string]bool{}

	tracks, _ := project["tracks"].([]any)
	for i, t := range tracks {
		info := infos[i]
		if info.Part != Tenor && info.Part != Baritone && info.Part != Bass {
			continue
		}

		track := t.(map[string]any)

		// Transpose mainGroup notes.
		if mainGroup, ok := track["mainGroup"].(map[string]any); ok {
			transposeNotes(mainGroup)
		}

		// Collect library group IDs referenced by this track.
		groups, _ := track["groups"].([]any)
		for _, g := range groups {
			group := g.(map[string]any)
			if groupID, ok := group["groupID"].(string); ok {
				transposeGroups[groupID] = true
			}
		}
	}

	// Transpose notes in matching library groups.
	library, _ := project["library"].([]any)
	var count int
	for _, lg := range library {
		group, ok := lg.(map[string]any)
		if !ok {
			continue
		}
		uuid, _ := group["uuid"].(string)
		if !transposeGroups[uuid] {
			continue
		}
		count += transposeNotes(group)
	}

	fmt.Printf("\nOctave fix: transposed %d notes +12 semitones in tenor/baritone/bass tracks.\n", count)
}

// transposeNotes shifts all note pitches in a note group up by 12 semitones.
// Returns the number of notes transposed.
func transposeNotes(noteGroup map[string]any) int {
	notes, _ := noteGroup["notes"].([]any)
	for _, n := range notes {
		note := n.(map[string]any)
		if pitch, ok := note["pitch"].(float64); ok {
			note["pitch"] = pitch + 12
		}
	}
	return len(notes)
}

// setNoteAttributes updates attributes on every note in the project.
func setNoteAttributes(project map[string]any) {
	attrs := map[string]any{
		"cPhraseTailDispersion": 0.2,
		"cPitchDispersion":      0.1,
		"cTimeDispersion":       0.2,
		"dF0VbrMod":             0.05,
		"expValueX":             -0.9,
		"expValueY":             -0.9,
	}

	var count int

	updateNotes := func(noteGroup map[string]any) {
		notes, _ := noteGroup["notes"].([]any)
		for _, n := range notes {
			note := n.(map[string]any)
			noteAttrs, ok := note["attributes"].(map[string]any)
			if !ok {
				noteAttrs = map[string]any{}
				note["attributes"] = noteAttrs
			}
			for k, v := range attrs {
				noteAttrs[k] = v
			}
			count++
		}
	}

	// Update notes in mainGroup of each track.
	tracks, _ := project["tracks"].([]any)
	for _, t := range tracks {
		track := t.(map[string]any)
		if mainGroup, ok := track["mainGroup"].(map[string]any); ok {
			updateNotes(mainGroup)
		}
	}

	// Update notes in library groups.
	library, _ := project["library"].([]any)
	for _, lg := range library {
		if group, ok := lg.(map[string]any); ok {
			updateNotes(group)
		}
	}

	fmt.Printf("\nUpdated attributes on %d notes.\n", count)
}
