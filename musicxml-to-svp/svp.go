package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// SVP output structs

type SVPProject struct {
	Version      int             `json:"version"`
	Time         SVPTime         `json:"time"`
	Library      []*SVPGroup     `json:"library"`
	Tracks       []*SVPTrack     `json:"tracks"`
	RenderConfig SVPRenderConfig `json:"renderConfig"`
	ProjectMixer SVPProjectMixer `json:"projectMixer"`
	UUID         string          `json:"uuid"`
}

type SVPTime struct {
	Meters []*SVPMeter `json:"meter"`
	Tempos []*SVPTempo `json:"tempo"`
}

type SVPMeter struct {
	Index       int `json:"index"`
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}

type SVPTempo struct {
	Position int64   `json:"position"`
	BPM      float64 `json:"bpm"`
}

type SVPGroup struct {
	Name       string        `json:"name"`
	UUID       string        `json:"uuid"`
	Notes      []*SVPNote    `json:"notes"`
	Parameters SVPParameters `json:"parameters"`
}

type SVPNote struct {
	Onset    int64    `json:"onset"`
	Duration int64    `json:"duration"`
	Lyrics   string   `json:"lyrics"`
	Phonemes string   `json:"phonemes"`
	Pitch    int      `json:"pitch"`
	Detune   int      `json:"detune"`
	Takes    SVPTakes `json:"attributes"`
}

type SVPTakes struct {
	EvenSyllableDuration  bool         `json:"evenSyllableDuration"`
	Muted                 bool         `json:"muted"`
	TF0Offset             float64      `json:"tF0Offset"`
	SystemPitchDelta      SVPParamMode `json:"systemPitchDelta"`
	CPhraseTailDispersion float64      `json:"cPhraseTailDispersion,omitempty"`
	CPitchDispersion      float64      `json:"cPitchDispersion,omitempty"`
	CTimeDispersion       float64      `json:"cTimeDispersion,omitempty"`
	DF0VbrMod             float64      `json:"dF0VbrMod,omitempty"`
	ExpValueX             float64      `json:"expValueX,omitempty"`
	ExpValueY             float64      `json:"expValueY,omitempty"`
	LanguageOverride      string       `json:"languageOverride,omitempty"`
	PhonesetOverride      string       `json:"phonesetOverride,omitempty"`
	Takes                 []SVPTake    `json:"dur,omitempty"`
}

type SVPParamMode struct {
	Mode string `json:"mode"`
}

type SVPTake struct {
	ID    int      `json:"id"`
	Liked bool     `json:"liked"`
	Seeds SVPSeeds `json:"seeds"`
}

type SVPSeeds struct {
	SingingSeed int `json:"singingSeed"`
	BackingSeed int `json:"backingSeed"`
}

type SVPParameters struct {
	PitchDelta   SVPParamCurve `json:"pitchDelta"`
	VibratoEnv   SVPParamCurve `json:"vibratoEnv"`
	Loudness     SVPParamCurve `json:"loudness"`
	Tension      SVPParamCurve `json:"tension"`
	Breathiness  SVPParamCurve `json:"breathiness"`
	Voicing      SVPParamCurve `json:"voicing"`
	Gender       SVPParamCurve `json:"gender"`
	ToneShift    SVPParamCurve `json:"toneShift"`
	MouthOpening SVPParamCurve `json:"mouthOpening"`
}

type SVPParamCurve struct {
	Mode   string    `json:"mode"`
	Points []float64 `json:"points"`
}

type SVPTrack struct {
	Name      string        `json:"name"`
	DispColor string        `json:"dispColor"`
	DispOrder int           `json:"dispOrder"`
	Mixer     SVPMixer      `json:"mixer"`
	MainGroup SVPGroupRef   `json:"mainGroup"`
	MainRef   SVPGroupRef   `json:"mainRef"`
	Groups    []SVPGroupRef `json:"groups"`
	UUID      string        `json:"uuid"`
}

type SVPGroupRef struct {
	GroupID        string       `json:"groupID"`
	BlickOffset    int64        `json:"blickOffset"`
	PitchOffset    int          `json:"pitchOffset"`
	IsInstrumental bool         `json:"isInstrumental"`
	Database       *SVPDatabase `json:"database,omitempty"`
	Voice          *SVPVoice    `json:"voice,omitempty"`
	UUID           string       `json:"uuid"`
}

type SVPDatabase struct {
	Name             string `json:"name"`
	Language         string `json:"language"`
	Phoneset         string `json:"phoneset"`
	LanguageOverride string `json:"languageOverride"`
	PhonesetOverride string `json:"phonesetOverride"`
	BackendType      string `json:"backendType"`
	Version          string `json:"version"`
}

type SVPVoice struct {
	RelaxedPronunciation   string             `json:"relaxedPronunciation,omitempty"`
	VocalModeInherited     bool               `json:"vocalModeInherited"`
	VocalModePreset        string             `json:"vocalModePreset"`
	VocalModeParams        map[string]float64 `json:"vocalModeParams"`
	ChoirSeatingSeparation float64            `json:"choirSeatingSeparation,omitempty"`
	ChoirNumStems          int                `json:"choirNumStems,omitempty"`
	ChoirPartName          string             `json:"choirPartName,omitempty"`
}

type SVPMixer struct {
	GainDecibel float64 `json:"gainDecibel"`
	Pan         float64 `json:"pan"`
	Mute        bool    `json:"mute"`
	Solo        bool    `json:"solo"`
	Display     bool    `json:"display"`
}

type SVPRenderConfig struct {
	Destination      string `json:"destination"`
	Filename         string `json:"filename"`
	NumChannels      int    `json:"numChannels"`
	AspirationFormat string `json:"aspirationFormat"`
	BitDepth         int    `json:"bitDepth"`
	SampleRate       int    `json:"sampleRate"`
	ExportMixDown    bool   `json:"exportMixDown"`
}

type SVPProjectMixer struct {
	GainDecibel      float64 `json:"gainDecibel"`
	Pan              float64 `json:"pan"`
	Mute             bool    `json:"mute"`
	Solo             bool    `json:"solo"`
	LinkRoomSettings bool    `json:"linkRoomSettings"`
}

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func newShortID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func newEmptyParamCurve() SVPParamCurve {
	return SVPParamCurve{Mode: "cubic", Points: []float64{}}
}

var dispColors = []string{
	"#6699cc", "#cc6699", "#99cc66", "#cc9966",
	"#9966cc", "#66cc99", "#cc6666", "#6666cc",
}

func newSVPNote(onset, duration int64, pitch, detune int, lyric string) *SVPNote {
	return &SVPNote{
		Onset:    onset,
		Duration: duration,
		Lyrics:   lyric,
		Phonemes: "",
		Pitch:    pitch,
		Detune:   detune,
		Takes: SVPTakes{
			EvenSyllableDuration: true,
			SystemPitchDelta:     SVPParamMode{Mode: "cubic"},
			Takes: []SVPTake{{
				ID:    0,
				Liked: false,
				Seeds: SVPSeeds{},
			}},
		},
	}
}

// graceNoteDurationFromType computes the duration of a grace note from its notated type.
func graceNoteDurationFromType(notatedType string, acciaccatura bool) int64 {
	dur := int64(blicksPerQuarter / 4) // default: sixteenth
	switch notatedType {
	case "whole":
		dur = blicksPerQuarter * 4
	case "half":
		dur = blicksPerQuarter * 2
	case "quarter":
		dur = blicksPerQuarter
	case "eighth":
		dur = blicksPerQuarter / 2
	case "16th":
		dur = blicksPerQuarter / 4
	case "32nd":
		dur = blicksPerQuarter / 8
	}
	if acciaccatura {
		dur /= 2
	}
	return dur
}

// capGraceDurs limits total grace note duration and scales proportionally.
func capGraceDurs(graces []GraceNote, maxTotal int64) (graceDurs []int64, totalGrace int64) {
	graceDurs = make([]int64, len(graces))
	for i, g := range graces {
		graceDurs[i] = graceNoteDurationFromType(g.NotatedType, g.Acciaccatura)
		totalGrace += graceDurs[i]
	}
	if totalGrace > maxTotal {
		scale := float64(maxTotal) / float64(totalGrace)
		totalGrace = 0
		for i := range graceDurs {
			graceDurs[i] = int64(float64(graceDurs[i]) * scale)
			totalGrace += graceDurs[i]
		}
		// Distribute rounding remainder to the last grace note.
		if remainder := maxTotal - totalGrace; remainder > 0 {
			graceDurs[len(graceDurs)-1] += remainder
			totalGrace = maxTotal
		}
	}
	return
}

// emitGraces creates SVPNote objects for grace notes at the given onset.
func emitGraces(graces []GraceNote, onset int64) []*SVPNote {
	var out []*SVPNote
	for _, g := range graces {
		out = append(out, newSVPNote(onset, g.Duration, g.Pitch, g.Detune, g.Lyric))
		onset += g.Duration
	}
	return out
}

// scoreToSVP converts a Score to an SVPProject.
func scoreToSVP(score *Score) *SVPProject {
	svpMeters := make([]*SVPMeter, len(score.Meters))
	for i, m := range score.Meters {
		svpMeters[i] = &SVPMeter{Index: m.MeasureIndex, Numerator: m.Numerator, Denominator: m.Denominator}
	}
	svpTempos := make([]*SVPTempo, len(score.Tempos))
	for i, t := range score.Tempos {
		svpTempos[i] = &SVPTempo{Position: t.Position, BPM: t.BPM}
	}
	if len(svpMeters) == 0 {
		svpMeters = append(svpMeters, &SVPMeter{Index: 0, Numerator: 4, Denominator: 4})
	}
	if len(svpTempos) == 0 {
		svpTempos = append(svpTempos, &SVPTempo{Position: 0, BPM: 120})
	}

	var library []*SVPGroup
	var tracks []*SVPTrack

	for partIdx, part := range score.Parts {
		var svpNotes []*SVPNote
		var accents []accentEvent

		for _, n := range part.Notes {
			onset := n.Onset
			duration := n.Duration

			// Emit leading grace notes (durations already computed in buildNotes).
			if len(n.LeadingGraces) > 0 {
				var totalGrace int64
				for _, g := range n.LeadingGraces {
					totalGrace += g.Duration
				}
				graceOnset := onset - totalGrace
				svpNotes = append(svpNotes, emitGraces(n.LeadingGraces, graceOnset)...)
			}

			// Apply articulation adjustments.
			hasTenuto := n.Articulations&ArticulationTenuto != 0
			if !hasTenuto {
				if n.Articulations&ArticulationStaccatissimo != 0 {
					duration = duration / 3
				} else if n.Articulations&ArticulationStaccato != 0 {
					duration = duration * 2 / 3
				}
			}

			// Collect accent events (use note duration for spike scaling).
			if n.Articulations&ArticulationStrongAccent != 0 {
				accents = append(accents, accentEvent{
					position: onset,
					duration: n.Duration,
					strong:   true,
				})
			} else if n.Articulations&ArticulationAccent != 0 {
				accents = append(accents, accentEvent{
					position: onset,
					duration: n.Duration,
				})
			}

			svpNotes = append(svpNotes, newSVPNote(onset, duration, n.Pitch, n.Detune, n.Lyric))

			// Emit trailing grace notes at the pre-staccato endpoint
			// (durations already computed in buildNotes).
			if len(n.TrailingGraces) > 0 {
				trailOnset := onset + n.Duration
				svpNotes = append(svpNotes, emitGraces(n.TrailingGraces, trailOnset)...)
			}
		}

		// Build curves from dynamics and accents.
		params := newEmptyParameters()
		if len(part.Dynamics) > 0 {
			params.Loudness.Points = buildCurve(part.Dynamics, func(e dynEvent) float64 { return e.loudness }, 6)
			params.Tension.Points = buildCurve(part.Dynamics, func(e dynEvent) float64 { return e.tension }, 0.15)
		}
		if len(accents) > 0 {
			params.Loudness.Points = applyAccents(params.Loudness.Points, accents, 1.5, 3)
			params.Tension.Points = applyAccents(params.Tension.Points, accents, 0.15, 0.3)
		}

		group := &SVPGroup{
			Name:       part.Name,
			UUID:       newUUID(),
			Notes:      svpNotes,
			Parameters: params,
		}
		library = append(library, group)

		mainGroup := &SVPGroup{
			Name:       "main",
			UUID:       newUUID(),
			Notes:      []*SVPNote{},
			Parameters: newEmptyParameters(),
		}

		color := dispColors[partIdx%len(dispColors)]

		track := &SVPTrack{
			Name:      part.Name,
			DispColor: color,
			DispOrder: partIdx,
			UUID:      newUUID(),
			Mixer: SVPMixer{
				GainDecibel: 0,
				Pan:         0,
				Display:     true,
			},
			MainGroup: SVPGroupRef{
				GroupID: mainGroup.UUID,
				UUID:    newUUID(),
			},
			MainRef: SVPGroupRef{
				GroupID: mainGroup.UUID,
				UUID:    newUUID(),
			},
			Groups: []SVPGroupRef{{
				GroupID: group.UUID,
				UUID:    newUUID(),
			}},
		}
		tracks = append(tracks, track)
		library = append(library, mainGroup)
	}

	return &SVPProject{
		Version: 196,
		UUID:    newUUID(),
		Time: SVPTime{
			Meters: svpMeters,
			Tempos: svpTempos,
		},
		Library: library,
		Tracks:  tracks,
		RenderConfig: SVPRenderConfig{
			Destination:      "",
			Filename:         "",
			NumChannels:      1,
			AspirationFormat: "noAspiration",
			BitDepth:         16,
			SampleRate:       44100,
			ExportMixDown:    true,
		},
		ProjectMixer: SVPProjectMixer{
			LinkRoomSettings: true,
		},
	}
}

func newEmptyParameters() SVPParameters {
	return SVPParameters{
		PitchDelta:   newEmptyParamCurve(),
		VibratoEnv:   newEmptyParamCurve(),
		Loudness:     newEmptyParamCurve(),
		Tension:      newEmptyParamCurve(),
		Breathiness:  newEmptyParamCurve(),
		Voicing:      newEmptyParamCurve(),
		Gender:       newEmptyParamCurve(),
		ToneShift:    newEmptyParamCurve(),
		MouthOpening: newEmptyParamCurve(),
	}
}
