package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// SVP output structs

type SVPProject struct {
	Version      int              `json:"version"`
	Time         SVPTime          `json:"time"`
	Library      []*SVPGroup      `json:"library"`
	Tracks       []*SVPTrack      `json:"tracks"`
	RenderConfig SVPRenderConfig  `json:"renderConfig"`
	ProjectMixer SVPProjectMixer  `json:"projectMixer"`
	UUID         string           `json:"uuid"`
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
	Onset            int64    `json:"onset"`
	Duration         int64    `json:"duration"`
	Lyrics           string   `json:"lyrics"`
	Phonemes         string   `json:"phonemes"`
	Pitch            int      `json:"pitch"`
	Detune           int      `json:"detune"`
	Takes            SVPTakes `json:"attributes"`
}

type SVPTakes struct {
	EvenSyllableDuration    bool         `json:"evenSyllableDuration"`
	Muted                   bool         `json:"muted"`
	TF0Offset               float64      `json:"tF0Offset"`
	SystemPitchDelta        SVPParamMode `json:"systemPitchDelta"`
	CPhraseTailDispersion   float64      `json:"cPhraseTailDispersion,omitempty"`
	CPitchDispersion        float64      `json:"cPitchDispersion,omitempty"`
	CTimeDispersion         float64      `json:"cTimeDispersion,omitempty"`
	DF0VbrMod               float64      `json:"dF0VbrMod,omitempty"`
	ExpValueX               float64      `json:"expValueX,omitempty"`
	ExpValueY               float64      `json:"expValueY,omitempty"`
	LanguageOverride        string       `json:"languageOverride,omitempty"`
	PhonesetOverride        string       `json:"phonesetOverride,omitempty"`
	Takes                   []SVPTake    `json:"dur,omitempty"`
}

type SVPParamMode struct {
	Mode string `json:"mode"`
}

type SVPTake struct {
	ID    int       `json:"id"`
	Liked bool      `json:"liked"`
	Seeds SVPSeeds  `json:"seeds"`
}

type SVPSeeds struct {
	SingingSeed  int `json:"singingSeed"`
	BackingSeed  int `json:"backingSeed"`
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
	Name      string       `json:"name"`
	DispColor string       `json:"dispColor"`
	DispOrder int          `json:"dispOrder"`
	Mixer     SVPMixer     `json:"mixer"`
	MainGroup SVPGroupRef  `json:"mainGroup"`
	MainRef   SVPGroupRef  `json:"mainRef"`
	Groups    []SVPGroupRef `json:"groups"`
	UUID      string       `json:"uuid"`
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
	RelaxedPronunciation   string            `json:"relaxedPronunciation,omitempty"`
	VocalModeInherited     bool              `json:"vocalModeInherited"`
	VocalModePreset        string            `json:"vocalModePreset"`
	VocalModeParams        map[string]float64 `json:"vocalModeParams"`
	ChoirSeatingSeparation float64           `json:"choirSeatingSeparation,omitempty"`
	ChoirNumStems          int               `json:"choirNumStems,omitempty"`
	ChoirPartName          string            `json:"choirPartName,omitempty"`
}

type SVPMixer struct {
	GainDecibel    float64 `json:"gainDecibel"`
	Pan            float64 `json:"pan"`
	Mute           bool    `json:"mute"`
	Solo           bool    `json:"solo"`
	Display        bool    `json:"display"`
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

func newEmptyParameters() SVPParameters {
	return SVPParameters{
		PitchDelta:  newEmptyParamCurve(),
		VibratoEnv:  newEmptyParamCurve(),
		Loudness:    newEmptyParamCurve(),
		Tension:     newEmptyParamCurve(),
		Breathiness: newEmptyParamCurve(),
		Voicing:     newEmptyParamCurve(),
		Gender:       newEmptyParamCurve(),
		ToneShift:    newEmptyParamCurve(),
		MouthOpening: newEmptyParamCurve(),
	}
}
