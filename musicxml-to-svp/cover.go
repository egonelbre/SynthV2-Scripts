package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/voice"
)

// Cover modes, selected with -cover.
const (
	coverNone    = "none"
	coverFormant = "formant"
	coverPhoneme = "phoneme"
	coverBoth    = "both"
)

// Vowel modification ("covering"): above the passaggio singers darken and
// round the vowel instead of pushing the open sound higher. The closest
// equivalent in SynthV is the gender parameter, which shifts the formants;
// negative values darken.
const (
	coverLeadIn    = 5    // semitones below the passaggio where covering starts
	coverSemitones = 7    // semitones above the passaggio for a full cover
	coverMaxGender = -0.3 // gender offset at full cover
	coverRampGap   = int64(blicksPerQuarter / 8)

	// defaultGender brightens every voice a little, which sits better in a
	// choir texture.
	defaultGender = 0.15
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

// coverGender darkens the notes around and above the part's passaggio by
// layering onto the given gender curve. Unrecognized voice parts are returned
// unchanged.
func coverGender(part Part, gender []float64) []float64 {
	threshold, ok := coverThresholds[voice.ParseVoicePart(part.Name)]
	if !ok {
		fmt.Fprintf(os.Stderr, "  skipping cover for %q (unknown voice part)\n", part.Name)
		return gender
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
		return gender
	}
	return addStressShapes(gender, shapes)
}

// coverAmount ramps the gender offset from nothing coverLeadIn semitones below
// the passaggio -- so the passaggio itself is already partly covered -- to the
// full cover coverSemitones above it.
func coverAmount(pitch, threshold int) float64 {
	above := pitch - threshold + coverLeadIn
	if above <= 0 {
		return 0
	}
	span := coverLeadIn + coverSemitones
	if above > span {
		above = span
	}
	return coverMaxGender * float64(above) / float64(span)
}

// Notes this far above the passaggio get their vowel substituted; the swap is
// audible, so it starts well above where the formant shift comes in.
const coverPhonemeMargin = 6

// coverVowelTables maps a note's language to the vowel substitutions used for
// covering. Spanish (which Latin is sung with) only has five vowels, so [a]
// can only move to [o] -- a stronger modification than the [ɔ] a singer would
// use. Languages without a table are left alone.
var coverVowelTables = map[string]map[string]string{
	"spanish": {
		"a": "o",
	},
	"english": {
		"aa": "ao", // f_a_r  -> f_ough_t
		"ae": "eh", // b_a_t  -> b_e_t
		"eh": "ah", // b_e_t  -> d_u_ck
		"ey": "eh",
		"iy": "ih", // b_ea_t -> b_i_t
	},
}

// coverPhonemes substitutes darker vowels on notes well above the part's
// passaggio. It only touches notes that already carry phonemes and a language
// override, which is what -lang and -p produce.
func coverPhonemes(library []*SVPGroup) {
	for _, group := range library {
		threshold, ok := coverThresholds[voice.ParseVoicePart(group.Name)]
		if !ok {
			continue
		}

		covered, skipped := 0, 0
		for _, note := range group.Notes {
			if note.Phonemes == "" || note.Pitch < threshold+coverPhonemeMargin {
				continue
			}
			table, ok := coverVowelTables[note.Takes.LanguageOverride]
			if !ok {
				skipped++
				continue
			}
			phonemes := strings.Fields(note.Phonemes)
			for i, p := range phonemes {
				if darker, ok := table[p]; ok {
					phonemes[i] = darker
					covered++
				}
			}
			note.Phonemes = strings.Join(phonemes, " ")
		}
		if covered > 0 {
			fmt.Fprintf(os.Stderr, "  %s: covered %d vowels\n", group.Name, covered)
		}
		if skipped > 0 {
			fmt.Fprintf(os.Stderr, "  %s: %d high notes without a known language, left uncovered\n", group.Name, skipped)
		}
	}
}
