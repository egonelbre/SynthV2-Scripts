// Package voice provides shared voice part types, choir/soloist databases,
// and helper functions for SynthV voice assignment.
package voice

import (
	"slices"
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

// Choirs contains the available choir voice databases.
var Choirs = []ChoirInfo{
	{
		Name:     "Choir Voices #1",
		Language: "english",
		Phoneset: "arpabet",
		Parts:    []VoicePart{Soprano, Alto, Tenor, Baritone, Bass},
	},
	{
		Name:     "Choir Voices #2",
		Language: "mandarin",
		Phoneset: "xsampa",
		Parts:    []VoicePart{Soprano, MezzoSoprano, Tenor, Baritone, Bass},
	},
	{
		Name:     "Choir Voices #3",
		Language: "japanese",
		Phoneset: "romaji",
		Parts:    []VoicePart{Soprano, Alto, Tenor, Baritone, Bass},
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

// Soloists maps voice parts to available soloist databases.
var Soloists = map[VoicePart][]SoloistDB{
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

// ParseIsSolo returns true if the name starts with the word "solo".
func ParseIsSolo(name string) bool {
	fields := strings.Fields(strings.TrimSpace(name))
	return len(fields) > 0 && strings.EqualFold(fields[0], "solo")
}

// stripSoloPrefix removes a leading "solo" word from the name, if present.
func stripSoloPrefix(name string) string {
	fields := strings.Fields(strings.TrimSpace(name))
	if len(fields) > 1 && strings.EqualFold(fields[0], "solo") {
		return strings.Join(fields[1:], " ")
	}
	return name
}

// ParseVoicePart extracts a VoicePart from a track name like "Soprano 1", "Bass", "Solo Alto", etc.
func ParseVoicePart(name string) VoicePart {
	lower := strings.ToLower(strings.TrimSpace(stripSoloPrefix(name)))

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

// ParsePartNum extracts the track number from a name like "Bass 2" -> 2, "Bass" -> 1, "Solo Bass 2" -> 2.
func ParsePartNum(name string) int {
	fields := strings.Fields(strings.TrimSpace(stripSoloPrefix(name)))
	if len(fields) >= 2 {
		// Try to parse the second field as a number (handles "Bass 2", "Soprano 1a" -> 1).
		numStr := fields[1]
		// Strip trailing non-digit characters (e.g. "1a" -> "1").
		for len(numStr) > 0 && (numStr[len(numStr)-1] < '0' || numStr[len(numStr)-1] > '9') {
			numStr = numStr[:len(numStr)-1]
		}
		if n, err := strconv.Atoi(numStr); err == nil && n > 0 {
			return n
		}
	}
	return 1
}

// MapPartToChoir maps a detected voice part to the closest choir part available.
func MapPartToChoir(part VoicePart, choir ChoirInfo) VoicePart {
	// Check if the part is directly available.
	if slices.Contains(choir.Parts, part) {
		return part
	}

	// Fallback mappings.
	switch part {
	case MezzoSoprano:
		return Alto
	case Alto:
		return MezzoSoprano
	case Baritone:
		return Bass
	case Bass:
		return Baritone
	default:
		return Soprano
	}
}

// TrackInfo holds parsed voice part information for a track.
type TrackInfo struct {
	Name    string
	Part    VoicePart
	PartNum int
	IsSolo  bool
}

// AssignSoloists resolves soloist assignments for a set of tracks.
// Returns a map from track index to the assigned SoloistDB.
func AssignSoloists(infos []TrackInfo) map[int]SoloistDB {
	assignments := make(map[int]SoloistDB)

	// Group track indices by part.
	partTracks := map[VoicePart][]int{}
	for i, info := range infos {
		if info.Part != Unknown {
			partTracks[info.Part] = append(partTracks[info.Part], i)
		}
	}

	for part, trackIdxs := range partTracks {
		pool := Soloists[part]
		if len(pool) == 0 {
			continue
		}

		assigned := map[int]bool{}  // track indices already assigned
		usedVoice := map[int]bool{} // pool indices already used

		// First pass: assign voices that have Preferred set.
		for vi, v := range pool {
			if len(v.Preferred) == 0 {
				continue
			}
			for _, prefNum := range v.Preferred {
				for _, ti := range trackIdxs {
					if !assigned[ti] && infos[ti].PartNum == prefNum {
						assignments[ti] = v
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
			for vi < len(pool) && usedVoice[vi] {
				vi++
			}
			if vi >= len(pool) {
				vi = 0
			}
			assignments[ti] = pool[vi]
			usedVoice[vi] = true
			vi++
		}
	}

	return assignments
}
