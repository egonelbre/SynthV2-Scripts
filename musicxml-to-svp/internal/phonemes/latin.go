package phonemes

import (
	"regexp"
	"strings"
)

// newLatin converts Ecclesiastical (Church) Latin, the pronunciation used in
// sung Latin. The Spanish phoneme set covers it almost exactly (pure vowels,
// rolled r, ch [tʃ], sh [ʃ], J [ɲ]), so no adaptive selection is needed.
func newLatin() *Converter {
	return &Converter{
		selectLangs: func(string) []string { return []string{"spanish"} },
		normalize:   normalizeLatin,
		vowels:      "aeiou",
		tables: map[string]*phoneTable{
			"spanish": latinSpanish,
		},
		words: map[string]Result{
			// mihi/nihil have an exceptional [k] for h.
			"mihi":  {Language: "spanish", Phoneset: "xsampa", Phonemes: "m i k i"},
			"nihil": {Language: "spanish", Phoneset: "xsampa", Phonemes: "n i k i l"},
		},
	}
}

var (
	latinInitialI = regexp.MustCompile(`^i([aeou])`)         // consonantal i: iesu → jesu
	latinNGW      = regexp.MustCompile(`ngu([aeiou])`)       // sanguis → [saŋgwis]
	latinTI       = regexp.MustCompile(`([^stx])ti([aeou])`) // gratia → [gratsia]
	latinXCSoft   = regexp.MustCompile(`xc([ei])`)           // excelsis → [ekʃelsis]
	latinSCSoft   = regexp.MustCompile(`sc([ei])`)           // → [ʃ]
	latinCSoft    = regexp.MustCompile(`c([ei])`)            // → [tʃ]
	latinGSoft    = regexp.MustCompile(`g([ei])`)            // → [dʒ]
)

// normalizeLatin rewrites Ecclesiastical Latin orthography into a phonetic
// intermediate form using placeholder runes č [tʃ], š [ʃ], ñ [ɲ], ǧ [dʒ]
// and w for the [w] glide, resolving all context-dependent spelling rules
// before table lookup.
func normalizeLatin(s string) string {
	// Ligatures and Greek-style spellings.
	s = strings.NewReplacer(
		"æ", "e", "œ", "e", "ae", "e", "oe", "e", "y", "i",
		"ch", "k", "ph", "f", "th", "t", "rh", "r",
	).Replace(s)
	s = strings.ReplaceAll(s, "h", "") // h is silent
	s = latinInitialI.ReplaceAllString(s, "j$1")
	s = strings.ReplaceAll(s, "gn", "ñ")
	s = strings.ReplaceAll(s, "qu", "kw")
	s = latinNGW.ReplaceAllString(s, "ngw$1")
	s = latinTI.ReplaceAllString(s, "${1}tsi$2")
	s = latinXCSoft.ReplaceAllString(s, "kš$1")
	s = strings.ReplaceAll(s, "x", "ks")
	s = latinSCSoft.ReplaceAllString(s, "š$1")
	s = latinCSoft.ReplaceAllString(s, "č$1")
	s = latinGSoft.ReplaceAllString(s, "ǧ$1")
	return s
}

var latinSpanish = &phoneTable{
	digraphs: map[string][]string{
		"rr": {"rr"},
		"ng": {"N", "g"},
	},
	expanded: map[rune][]string{
		'z': {"d", "s"}, // [dz]
	},
	singles: map[rune]string{
		'b': "b", 'c': "k", 'd': "d", 'f': "f", 'g': "g",
		'j': "I", 'k': "k", 'l': "l", 'm': "m", 'n': "n",
		'p': "p", 'r': "r", 's': "s", 't': "t", 'v': "B", 'w': "U",
		'č': "ch", 'š': "sh", 'ñ': "J", 'ǧ': "y",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
	},
}
