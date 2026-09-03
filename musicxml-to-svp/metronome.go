package main

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
)

const metronomeSampleRate = 44100

// SVPAudio links an audio file into an instrumental track.
type SVPAudio struct {
	Filename        string    `json:"filename"`
	Duration        float64   `json:"duration"`
	BPM             float64   `json:"bpm"`
	AlternativeBPMs []float64 `json:"alternativeBPMs"`
	BeatLocations   []float64 `json:"beatLocations"`
}

// metronomeBeat is one click: its position in blicks and whether it starts a measure.
type metronomeBeat struct {
	pos      int64
	downbeat bool
}

// metronomeBeats lists every beat of the played measures. Measures are laid out
// by their actual content so pickups and short final measures stay aligned.
func metronomeBeats(unrolled []playedMeasure, meters []MeterChange, songEnd int64) []metronomeBeat {
	var beats []metronomeBeat
	num, den := 4, 4
	mi := 0
	for i, pm := range unrolled {
		for mi < len(meters) && meters[mi].MeasureIndex <= i {
			num, den = meters[mi].Numerator, meters[mi].Denominator
			mi++
		}
		end := songEnd
		if i+1 < len(unrolled) {
			end = unrolled[i+1].startBlicks
		}
		beat := int64(4 * blicksPerQuarter / den)
		full := int64(num) * beat
		// A pickup is shorter than the meter: count its beats back from
		// the next downbeat so none of them is accented.
		start := pm.startBlicks
		if i == 0 && end-start < full {
			start = end - full
		}
		for k, pos := 0, start; pos < end; k, pos = k+1, pos+beat {
			if pos >= pm.startBlicks {
				beats = append(beats, metronomeBeat{pos: pos, downbeat: k == 0})
			}
		}
	}
	return beats
}

// blicksToSeconds converts a position through the tempo map.
func blicksToSeconds(tempos []TempoChange, pos int64) float64 {
	bpm := 120.0
	if len(tempos) > 0 {
		bpm = tempos[0].BPM
	}
	secs := 0.0
	cur := int64(0)
	for _, t := range tempos {
		if t.Position >= pos {
			break
		}
		if t.Position > cur {
			secs += float64(t.Position-cur) / blicksPerQuarter * 60 / bpm
			cur = t.Position
		}
		bpm = t.BPM
	}
	return secs + float64(pos-cur)/blicksPerQuarter*60/bpm
}

// renderClicks writes a mono 16-bit WAV with a short sine click at each beat.
func renderClicks(beats []metronomeBeat, tempos []TempoChange, songEnd int64) ([]byte, []float64) {
	duration := blicksToSeconds(tempos, songEnd) + 0.5
	samples := make([]float64, int(duration*metronomeSampleRate))
	var locations []float64
	for _, b := range beats {
		t := blicksToSeconds(tempos, b.pos)
		locations = append(locations, t)
		freq, gain := 880.0, 0.6
		if b.downbeat {
			freq, gain = 1320.0, 0.9
		}
		start := int(t * metronomeSampleRate)
		for i := 0; i < metronomeSampleRate/25 && start+i < len(samples); i++ {
			x := float64(i) / metronomeSampleRate
			samples[start+i] += gain * math.Exp(-x*120) * math.Sin(2*math.Pi*freq*x)
		}
	}

	var pcm bytes.Buffer
	for _, s := range samples {
		binary.Write(&pcm, binary.LittleEndian, int16(max(-1, min(1, s))*32767))
	}
	var w bytes.Buffer
	w.WriteString("RIFF")
	binary.Write(&w, binary.LittleEndian, uint32(36+pcm.Len()))
	w.WriteString("WAVEfmt ")
	binary.Write(&w, binary.LittleEndian, uint32(16))
	binary.Write(&w, binary.LittleEndian, uint16(1)) // PCM
	binary.Write(&w, binary.LittleEndian, uint16(1)) // mono
	binary.Write(&w, binary.LittleEndian, uint32(metronomeSampleRate))
	binary.Write(&w, binary.LittleEndian, uint32(metronomeSampleRate*2))
	binary.Write(&w, binary.LittleEndian, uint16(2))
	binary.Write(&w, binary.LittleEndian, uint16(16))
	w.WriteString("data")
	binary.Write(&w, binary.LittleEndian, uint32(pcm.Len()))
	w.Write(pcm.Bytes())
	return w.Bytes(), locations
}

// addMetronomeTrack renders the click track next to the project file and
// links it as a muted-by-default instrumental track.
func addMetronomeTrack(project *SVPProject, score *Score, unrolled []playedMeasure, svpPath string) error {
	songEnd := int64(0)
	for _, p := range score.Parts {
		for _, n := range p.Notes {
			songEnd = max(songEnd, n.Onset+n.Duration)
		}
	}
	beats := metronomeBeats(unrolled, score.Meters, songEnd)
	wav, locations := renderClicks(beats, score.Tempos, songEnd)

	ext := filepath.Ext(svpPath)
	wavPath, err := filepath.Abs(svpPath[:len(svpPath)-len(ext)] + "-metronome.wav")
	if err != nil {
		return err
	}
	if err := os.WriteFile(wavPath, wav, 0644); err != nil {
		return err
	}

	bpm := 120.0
	if len(score.Tempos) > 0 {
		bpm = score.Tempos[0].BPM
	}
	// Sit just under the combined level of the vocal tracks.
	power := 0.0
	for _, t := range project.Tracks {
		power += math.Pow(10, t.Mixer.GainDecibel/10)
	}
	gain := 10*math.Log10(power) - 6

	mainGroup := &SVPGroup{
		Name:       "main",
		UUID:       newUUID(),
		Notes:      []*SVPNote{},
		Parameters: newEmptyParameters(),
	}
	project.Library = append(project.Library, mainGroup)
	project.Tracks = append(project.Tracks, &SVPTrack{
		Name:      "Metronome",
		DispColor: "ff888888",
		DispOrder: len(project.Tracks),
		UUID:      newUUID(),
		Mixer:     SVPMixer{GainDecibel: gain, Display: true},
		MainGroup: SVPGroupRef{GroupID: mainGroup.UUID, UUID: newUUID()},
		MainRef: SVPGroupRef{
			GroupID:          mainGroup.UUID,
			UUID:             newUUID(),
			IsInstrumental:   true,
			BlickAbsoluteEnd: int64(float64(len(wav)-44) / 2 / metronomeSampleRate * bpm / 60 * blicksPerQuarter),
			Audio: &SVPAudio{
				Filename:        wavPath,
				Duration:        float64(len(wav)-44) / 2 / metronomeSampleRate,
				BPM:             bpm,
				AlternativeBPMs: []float64{bpm},
				BeatLocations:   locations,
			},
		},
		Groups: []SVPGroupRef{},
	})
	return nil
}
