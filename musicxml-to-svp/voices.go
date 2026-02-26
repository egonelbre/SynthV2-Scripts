package main

import (
	"fmt"
	"os"

	"github.com/egonelbre/synthv2-scripts/internal/phonemes"
	"github.com/egonelbre/synthv2-scripts/internal/voice"
)

func assignVoices(tracks []*SVPTrack, voiceArg string) {
	// Parse track names into voice parts.
	infos := make([]voice.TrackInfo, len(tracks))
	for i, t := range tracks {
		infos[i] = voice.TrackInfo{
			Name:    t.Name,
			Part:    voice.ParseVoicePart(t.Name),
			PartNum: voice.ParsePartNum(t.Name),
		}
	}

	switch voiceArg {
	case "choir1", "1":
		applyChoirToTracks(tracks, infos, voice.Choirs[0])
	case "choir2", "2":
		applyChoirToTracks(tracks, infos, voice.Choirs[1])
	case "choir3", "3":
		applyChoirToTracks(tracks, infos, voice.Choirs[2])
	case "soloists", "solo", "4":
		applySoloistsToTracks(tracks, infos)
	default:
		fmt.Fprintf(os.Stderr, "unknown voice source: %q (options: choir1, choir2, choir3, soloists)\n", voiceArg)
		os.Exit(1)
	}
}

func applyChoirToTracks(tracks []*SVPTrack, infos []voice.TrackInfo, choir voice.ChoirInfo) {
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
				VocalModeInherited:     true,
				VocalModeParams:        map[string]float64{},
				ChoirNumStems:          4,
				ChoirSeatingSeparation: 0.7,
				ChoirPartName:          partName,
			}
		}
	}
}

func applySoloistsToTracks(tracks []*SVPTrack, infos []voice.TrackInfo) {
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
			VocalModeInherited: true,
			VocalModeParams:    map[string]float64{},
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
