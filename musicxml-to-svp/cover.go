package main

import (
	"fmt"
	"os"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/voice"
)

// Vowel modification ("covering"): above the passaggio singers darken and
// round the vowel instead of pushing the open sound higher. The closest
// equivalent in SynthV is the gender parameter, which shifts the formants;
// negative values darken.
const (
	coverSemitones = 7    // semitones above the passaggio for a full cover
	coverMaxGender = -0.3 // gender offset at full cover
	coverRampGap   = int64(blicksPerQuarter / 8)
)

// coverThresholds maps a voice part to the MIDI pitch where covering starts,
// roughly the part's secondo passaggio.
var coverThresholds = map[voice.VoicePart]int{
	voice.Soprano:      77, // F5
	voice.MezzoSoprano: 76, // E5
	voice.Alto:         74, // D5
	voice.Tenor:        65, // F4
	voice.Baritone:     62, // D4
	voice.Bass:         60, // C4
}

// coverGender returns gender curve points that darken the notes above the
// part's passaggio, or nil when the part isn't a recognized voice.
func coverGender(part Part) []float64 {
	threshold, ok := coverThresholds[voice.ParseVoicePart(part.Name)]
	if !ok {
		fmt.Fprintf(os.Stderr, "  skipping cover for %q (unknown voice part)\n", part.Name)
		return nil
	}

	var shapes [][]float64
	for _, n := range part.Notes {
		amount := coverAmount(n.Pitch, threshold)
		if amount == 0 {
			continue
		}
		end := n.Onset + n.Duration
		shape := []float64{float64(n.Onset - coverRampGap), 0}
		if n.Duration > 3*coverRampGap {
			shape = append(shape,
				float64(n.Onset+coverRampGap), amount,
				float64(end-coverRampGap), amount)
		} else {
			shape = append(shape, float64(n.Onset+n.Duration/2), amount)
		}
		shape = append(shape, float64(end+coverRampGap), 0)
		shapes = append(shapes, shape)
	}
	if len(shapes) == 0 {
		return nil
	}
	return addStressShapes(nil, shapes)
}

// coverAmount ramps the gender offset from nothing at the passaggio to the
// full cover coverSemitones above it.
func coverAmount(pitch, threshold int) float64 {
	above := pitch - threshold + 1
	if above <= 0 {
		return 0
	}
	if above > coverSemitones {
		above = coverSemitones
	}
	return coverMaxGender * float64(above) / coverSemitones
}
