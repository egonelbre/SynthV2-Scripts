package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/phonemes"
	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/voice"
)

func assignVoices(tracks []*SVPTrack, voiceArg string, relaxed bool, panScheme string) {
	relaxedStr := "false"
	if relaxed {
		relaxedStr = "true"
	}

	// Parse track names into voice parts.
	infos := make([]voice.TrackInfo, len(tracks))
	for i, t := range tracks {
		infos[i] = voice.TrackInfo{
			Name:    t.Name,
			Part:    voice.ParseVoicePart(t.Name),
			PartNum: voice.ParsePartNum(t.Name),
			IsSolo:  voice.ParseIsSolo(t.Name),
		}
	}

	switch voiceArg {
	case "choir1", "1":
		applyChoirWithSoloists(tracks, infos, voice.Choirs[0], relaxedStr)
	case "choir2", "2":
		applyChoirWithSoloists(tracks, infos, voice.Choirs[1], relaxedStr)
	case "choir3", "3":
		applyChoirWithSoloists(tracks, infos, voice.Choirs[2], relaxedStr)
	case "soloists", "solo", "4":
		applySoloistsToTracks(tracks, infos, relaxedStr)
	default:
		fmt.Fprintf(os.Stderr, "unknown voice source: %q (options: choir1, choir2, choir3, soloists)\n", voiceArg)
		os.Exit(1)
	}

	// Apply panning.
	panning := computePanning(infos, panScheme)
	for i, t := range tracks {
		if pan, ok := panning[i]; ok {
			t.Mixer.Pan = pan
		}
	}
}

// computePanning returns pan values keyed by track index.
func computePanning(infos []voice.TrackInfo, scheme string) map[int]float64 {
	result := map[int]float64{}

	switch scheme {
	case "center":
		for i, info := range infos {
			if info.Part != voice.Unknown {
				result[i] = 0
			}
		}

	case "spread":
		spreadPan := map[voice.VoicePart]float64{
			voice.Soprano:      -0.6,
			voice.MezzoSoprano: -0.3,
			voice.Alto:         0.6,
			voice.Tenor:        -0.3,
			voice.Baritone:     0.3,
			voice.Bass:         0.3,
		}
		for i, info := range infos {
			if pan, ok := spreadPan[info.Part]; ok {
				result[i] = pan
			}
		}

	default: // "default"
		// Choir seating order left-to-right: S, MS, B, Bar, T, A.
		// Tracks are placed as slots in this order, evenly spaced across
		// the stereo field. Within each part section, sub-numbering follows
		// the convention: ascending for left/center-left parts, descending
		// for center-right/right parts.
		seatingOrder := []voice.VoicePart{
			voice.Soprano, voice.MezzoSoprano, voice.Bass, voice.Baritone, voice.Tenor, voice.Alto,
		}
		ascending := map[voice.VoicePart]bool{
			voice.Soprano: true, voice.MezzoSoprano: true,
			voice.Bass: false, voice.Baritone: false,
			voice.Tenor: true, voice.Alto: false,
		}

		type slot struct {
			trackIdx int
		}
		var slots []slot

		for _, part := range seatingOrder {
			var partTracks []struct{ idx, num int }
			for i, info := range infos {
				if info.Part == part {
					partTracks = append(partTracks, struct{ idx, num int }{i, info.PartNum})
				}
			}
			sort.Slice(partTracks, func(a, b int) bool {
				if ascending[part] {
					return partTracks[a].num < partTracks[b].num
				}
				return partTracks[a].num > partTracks[b].num
			})
			for _, pt := range partTracks {
				slots = append(slots, slot{trackIdx: pt.idx})
			}
		}

		n := len(slots)
		if n == 1 {
			result[slots[0].trackIdx] = 0
		}
		for i, s := range slots {
			if n <= 1 {
				break
			}
			t := float64(i) / float64(n-1)
			result[s.trackIdx] = -0.65 + t*1.3
		}
	}

	return result
}

// applyChoirWithSoloists splits tracks into solo and non-solo groups,
// applying soloist voices to solo tracks and choir voices to the rest.
func applyChoirWithSoloists(tracks []*SVPTrack, infos []voice.TrackInfo, choir voice.ChoirInfo, relaxed string) {
	var soloTracks []*SVPTrack
	var soloInfos []voice.TrackInfo
	var choirTracks []*SVPTrack
	var choirInfos []voice.TrackInfo

	for i, info := range infos {
		if info.IsSolo {
			soloTracks = append(soloTracks, tracks[i])
			soloInfos = append(soloInfos, info)
		} else {
			choirTracks = append(choirTracks, tracks[i])
			choirInfos = append(choirInfos, info)
		}
	}

	if len(choirTracks) > 0 {
		applyChoirToTracks(choirTracks, choirInfos, choir, relaxed)
	}
	if len(soloTracks) > 0 {
		applySoloistsToTracks(soloTracks, soloInfos, relaxed)
	}
}

func applyChoirToTracks(tracks []*SVPTrack, infos []voice.TrackInfo, choir voice.ChoirInfo, relaxed string) {
	db := &SVPDatabase{
		Name:     choir.Name,
		Language: choir.Language,
		Phoneset: choir.Phoneset,
		BackendType: "SVR3",
		Version:     "202",
	}

	for i, t := range tracks {
		info := infos[i]
		if info.Part == voice.Unknown {
			fmt.Fprintf(os.Stderr, "  skipping %q (unknown voice part)\n", info.Name)
			continue
		}

		choirPart := voice.MapPartToChoir(info.Part, choir)
		fmt.Fprintf(os.Stderr, "  %s -> %s (%s)\n", info.Name, choir.Name, choirPart.String())

		// Set database and voice on MainRef.
		t.MainRef.Database = db
		t.MainRef.Voice = &SVPVoice{
			RelaxedPronunciation:   relaxed,
			VocalModeInherited:     true,
			VocalModeParams:        map[string]float64{},
			ChoirSeatingSeparation: 0.7,
		}

		// Set choir part on group refs.
		for j := range t.Groups {
			partName := ""
			if choirPart != voice.Soprano {
				partName = string(choirPart)
			}
			t.Groups[j].Voice = &SVPVoice{
				RelaxedPronunciation:   relaxed,
				VocalModeInherited:     true,
				VocalModeParams:        map[string]float64{},
				ChoirNumStems:          4,
				ChoirSeatingSeparation: 0.7,
				ChoirPartName:          partName,
			}
		}
	}
}

func applySoloistsToTracks(tracks []*SVPTrack, infos []voice.TrackInfo, relaxed string) {
	assignments := voice.AssignSoloists(infos)

	for i, t := range tracks {
		info := infos[i]
		soloist, ok := assignments[i]
		if !ok {
			if info.Part == voice.Unknown {
				fmt.Fprintf(os.Stderr, "  skipping %q (unknown voice part)\n", info.Name)
			}
			continue
		}

		fmt.Fprintf(os.Stderr, "  %s -> %s\n", info.Name, soloist.Name)

		t.MainRef.Database = &SVPDatabase{
			Name:        soloist.Name,
			Language:    soloist.Language,
			Phoneset:    soloist.Phoneset,
			BackendType: "SVR3",
			Version:     "201",
		}
		t.MainRef.Voice = &SVPVoice{
			RelaxedPronunciation: relaxed,
			VocalModeInherited:   true,
			VocalModeParams:      map[string]float64{},
		}
	}
}

func applyPhonemes(library []*SVPGroup, conv *phonemes.Converter) {
	for _, group := range library {
		for _, note := range group.Notes {
			result := conv.Convert(note.Lyrics)
			if result.Phonemes != "" {
				note.Phonemes = result.Phonemes
				note.Takes.LanguageOverride = result.Language
				note.Takes.PhonesetOverride = result.Phoneset
			}
		}
	}
}

func parseLyricReplacement(v string) lyricReplacement {
	// Format: "[language] phonemes" or just "phonemes"
	if after, ok := strings.CutPrefix(v, "["); ok {
		lang, phoneme, _ := strings.Cut(after, "]")
		lang = strings.TrimSpace(lang)
		phoneme = strings.TrimSpace(phoneme)
		return lyricReplacement{
			Phonemes: phoneme,
			Language: lang,
			Phoneset: phonesetForLanguage(lang),
		}
	}
	return lyricReplacement{Phonemes: v}
}

func phonesetForLanguage(lang string) string {
	switch lang {
	case "japanese":
		return "romaji"
	case "english":
		return "arpabet"
	default:
		return "xsampa"
	}
}

type lyricReplacement struct {
	Phonemes string
	Language string
	Phoneset string
}

func applyLyricReplacements(library []*SVPGroup, replacements map[string]lyricReplacement) {
	for _, group := range library {
		for _, note := range group.Notes {
			if r, ok := replacements[note.Lyrics]; ok {
				note.Phonemes = r.Phonemes
				if r.Language != "" {
					note.Takes.LanguageOverride = r.Language
					note.Takes.PhonesetOverride = r.Phoneset
				}
			}
		}
	}
}

func setNoteAttributes(library []*SVPGroup) {
	for _, group := range library {
		for _, note := range group.Notes {
			note.Takes.CPhraseTailDispersion = 0.2
			note.Takes.CPitchDispersion = 0.1
			note.Takes.CTimeDispersion = 0.2
			note.Takes.DF0VbrMod = 0.05
			note.Takes.ExpValueX = -0.9
			note.Takes.ExpValueY = -0.9
		}
	}
}
