// Command svp-voice assigns voice databases to tracks in SynthV .svp project files
// based on track names (e.g., "Soprano 1", "Tenor 2", "Bass").
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/egonelbre/synthv2-scripts/internal/voice"
)

type trackInfo struct {
	Name    string
	Part    voice.VoicePart
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
		part := voice.ParseVoicePart(name)

		db := ""
		if mainRef, ok := track["mainRef"].(map[string]any); ok {
			if database, ok := mainRef["database"].(map[string]any); ok {
				db = fmt.Sprintf("%v", database["name"])
			}
		}

		infos = append(infos, trackInfo{Name: name, Part: part, PartNum: voice.ParsePartNum(name), DB: db})
	}

	// Display tracks
	fmt.Println("Tracks:")
	for i, info := range infos {
		partStr := info.Part.String()
		if info.Part == voice.Unknown {
			partStr = "???"
		}
		fmt.Printf("  %d. %-20s  [%s]  (current: %s)\n", i+1, info.Name, partStr, info.DB)
	}
	fmt.Println()

	// Check for unknown parts
	for _, info := range infos {
		if info.Part == voice.Unknown {
			fmt.Fprintf(os.Stderr, "Warning: cannot determine voice part for track %q\n", info.Name)
		}
	}

	switch voiceArg {
	case "choir1", "1":
		applyChoir(tracks, infos, voice.Choirs[0])
	case "choir2", "2":
		applyChoir(tracks, infos, voice.Choirs[1])
	case "choir3", "3":
		applyChoir(tracks, infos, voice.Choirs[2])
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

func applyChoir(tracks []any, infos []trackInfo, choir voice.ChoirInfo) {
	fmt.Printf("\nApplying %s to all tracks...\n", choir.Name)

	for i, t := range tracks {
		track := t.(map[string]any)
		info := infos[i]

		if info.Part == voice.Unknown {
			fmt.Printf("  Skipping %q (unknown voice part)\n", info.Name)
			continue
		}

		choirPart := voice.MapPartToChoir(info.Part, choir)
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
		v := mainRef["voice"].(map[string]any)
		v["vocalModeInherited"] = true
		v["vocalModePreset"] = ""
		if _, ok := v["vocalModeParams"]; !ok {
			v["vocalModeParams"] = map[string]any{}
		}
		// Set choir separation, preserve existing if present.
		if _, ok := v["choirSeatingSeparation"]; !ok {
			v["choirSeatingSeparation"] = 0.7
		}
		// Remove soloist fields that don't belong on choir mainRef.
		delete(v, "choirNumStems")
		// choirPartName is not set on mainRef (Soprano is default).
		delete(v, "choirPartName")

		// Update groups
		groups, _ := track["groups"].([]any)
		for _, g := range groups {
			group := g.(map[string]any)
			gVoice := group["voice"].(map[string]any)

			// Set choir part name (omit for Soprano as it's the default).
			if choirPart != voice.Soprano {
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
	// Convert trackInfo to voice.TrackInfo for the shared assignment logic.
	vInfos := make([]voice.TrackInfo, len(infos))
	for i, info := range infos {
		vInfos[i] = voice.TrackInfo{Name: info.Name, Part: info.Part, PartNum: info.PartNum}
	}
	assignments := voice.AssignSoloists(vInfos)

	fmt.Println("\nAssigning soloists to tracks...")
	fmt.Println()

	for i, t := range tracks {
		track := t.(map[string]any)
		info := infos[i]

		soloist, ok := assignments[i]
		if !ok {
			if info.Part == voice.Unknown {
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
		v := mainRef["voice"].(map[string]any)
		v["vocalModeInherited"] = true
		v["vocalModePreset"] = ""
		if _, ok := v["vocalModeParams"]; !ok {
			v["vocalModeParams"] = map[string]any{}
		}
		delete(v, "choirSeatingSeparation")
		delete(v, "choirNumStems")
		delete(v, "choirPartName")

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
		if info.Part != voice.Tenor && info.Part != voice.Baritone && info.Part != voice.Bass {
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
