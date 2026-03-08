package phonemes

import (
	"strings"
	"unicode/utf8"
)

func newGerman() *Converter {
	return &Converter{
		selectLang: selectGerman,
		normalize:  normalizeGerman,
		vowels:     "aeiouyäöü",
		tables: map[string]*phoneTable{
			"mandarin":  germanMandarin,
			"cantonese": germanCantonese,
			"japanese":  germanJapanese,
			"english":   germanEnglish,
			"korean":    germanKorean,
			"spanish":   germanSpanish,
		},
	}
}

// selectGerman picks the best target language based on word characteristics.
//
// Priority:
//  1. Words with ü → Mandarin (exact 'y' phoneme for front rounded vowel)
//  2. Words with ö → Cantonese (best '9' approximation)
//  3. Words with standalone ch → Spanish (ach-Laut after a/o/u) or Japanese (ich-Laut)
//  4. Words with ä → Cantonese (good 'E' match)
//  5. Words with r → Korean (good '4' flap approximation)
//  6. Words with ß → English (good 's' match)
//  7. Simple words → Japanese (cleanest basic vowels)
//  8. Default → English
func selectGerman(word string) string {
	if strings.ContainsRune(word, 'ü') {
		return "mandarin"
	}
	if strings.ContainsRune(word, 'ö') {
		return "cantonese"
	}
	if idx := findStandaloneCh(word); idx >= 0 {
		if idx > 0 {
			prev, _ := utf8.DecodeLastRuneInString(word[:idx])
			if prev == 'a' || prev == 'o' || prev == 'u' {
				return "spanish"
			}
		}
		return "japanese"
	}
	if strings.ContainsRune(word, 'ä') {
		return "cantonese"
	}
	if strings.ContainsRune(word, 'r') {
		return "korean"
	}
	if strings.ContainsRune(word, 'ß') {
		return "english"
	}
	return "japanese"
}

// findStandaloneCh finds "ch" that is not part of "sch".
// Returns byte index of 'c', or -1 if not found.
func findStandaloneCh(word string) int {
	for i := 0; i < len(word); {
		idx := strings.Index(word[i:], "ch")
		if idx < 0 {
			return -1
		}
		idx += i
		if idx == 0 || word[idx-1] != 's' {
			return idx
		}
		i = idx + 2
	}
	return -1
}

// normalizeGerman replaces trigraph "sch" with placeholder ʃ (U+0283),
// ß with plain s, and x with ks before table lookup. Word-initial "sp"
// and "st" are replaced with ʃp/ʃt since German palatalizes these only
// at syllable onset.
func normalizeGerman(s string) string {
	s = strings.ReplaceAll(s, "sch", "\u0283")
	if strings.HasPrefix(s, "sp") {
		s = "\u0283p" + s[2:]
	} else if strings.HasPrefix(s, "st") {
		s = "\u0283t" + s[2:]
	}
	s = strings.ReplaceAll(s, "ß", "s")
	s = strings.ReplaceAll(s, "x", "ks")
	return s
}

// germanMandarin is best for words with ü (exact 'y' phoneme match).
var germanMandarin = &phoneTable{
	digraphs: map[string][]string{
		"ch":      {"x"},
		"ck":      {"k"},
		"ng":      {"N"},
		"nk":      {"N", "k"},
		"pf":      {"p", "f"},
		"qu":      {"k", "w"},
		"ts":      {"ts"},
		"\u0283t": {"s`", "t"},
		"\u0283p": {"s`", "p"},
		"ss":      {"s", "s"},
		"ei":      {"a", "i"},
		"ie":      {"i", "i"},
		"eu":      {"o", "i"},
		"äu":      {"o", "i"},
		"au":      {"a", "u"},
	},
	singles: map[rune]string{
		'z': "ts", 'w': "w", 'v': "f", 'j': "j",
		'h': "x", 'l': "l", 'm': "m", 'n': "n",
		'r': "r\\`", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ä': "A", 'ö': "@", 'ü': "y", 'y': "y",
		'\u0283': "s`", // sch
	},
}

// germanCantonese is best for words with ö (best '9' approximation) and ä (good 'E' match).
var germanCantonese = &phoneTable{
	digraphs: map[string][]string{
		"ch": {"h"},
		"ck": {"k"},
		"ng": {"N"},
		"nk": {"N", "k"},
		"pf": {"p", "f"},
		"qu": {"k", "w"},
		"ts": {"ts"},
		"ss": {"s", "s"},
		"ei": {"6", "i"},
		"ie": {"i", "i"},
		"eu": {"O", "i"},
		"äu": {"O", "i"},
		"au": {"a", "u"},
	},
	singles: map[rune]string{
		'z': "ts", 'w': "w", 'v': "f", 'j': "j",
		'h': "h", 'l': "l", 'm': "m", 'n': "n",
		'r': "l", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ä': "E", 'ö': "9", 'ü': "y", 'y': "y",
		'\u0283': "s", // sch
	},
}

// germanJapanese is best for simple words and ich-Laut (ch after front vowels).
var germanJapanese = &phoneTable{
	digraphs: map[string][]string{
		"ch": {"h"},
		"ck": {"k"},
		"ng": {"N", "g"},
		"nk": {"N", "k"},
		"pf": {"p", "f"},
		"qu": {"k", "w"},
		"ts": {"ts"},
		"ss": {"s", "s"},
		"ei": {"a", "i"},
		"ie": {"i", "i"},
		"eu": {"o", "i"},
		"äu": {"o", "i"},
		"au": {"a", "u"},
	},
	singles: map[rune]string{
		'z': "ts", 'w': "v", 'v': "f", 'j': "y",
		'h': "h", 'l': "r", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ä': "a", 'ö': "o", 'ü': "u", 'y': "u",
		'\u0283': "sh", // sch
	},
}

// germanEnglish is the default fallback (most compatible).
var germanEnglish = &phoneTable{
	digraphs: map[string][]string{
		"ch": {"hh"},
		"ck": {"k"},
		"ng": {"ng"},
		"nk": {"ng", "k"},
		"pf": {"p", "f"},
		"qu": {"k", "w"},
		"ts": {"t", "s"},
		"ss": {"s", "s"},
		"ei": {"ay"},
		"ie": {"iy", "iy"},
		"eu": {"oy"},
		"äu": {"oy"},
		"au": {"aw"},
	},
	expanded: map[rune][]string{
		'z': {"t", "s"},
	},
	singles: map[rune]string{
		'w': "v", 'v': "f", 'j': "y",
		'h': "hh", 'l': "l", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f",
		'a': "aa", 'e': "eh", 'i': "iy", 'o': "ow", 'u': "uw",
		'ä': "ae", 'ö': "er", 'ü': "iy", 'y': "iy",
		'\u0283': "sh", // sch
	},
}

// germanKorean is best for words with r (good '4' flap approximation).
var germanKorean = &phoneTable{
	digraphs: map[string][]string{
		"ch": {"h"},
		"ck": {"k"},
		"ng": {"N"},
		"nk": {"N", "k"},
		"pf": {"p", "p"},
		"qu": {"k", "w"},
		"ts": {"ts_h"},
		"ss": {"s_t"},
		"ei": {"6", "i"},
		"ie": {"i", "i"},
		"eu": {"o", "i"},
		"äu": {"o", "i"},
		"au": {"6", "M"},
	},
	singles: map[rune]string{
		'z': "ts_h", 'w': "b", 'v': "p", 'j': "j",
		'h': "h", 'l': "l", 'm': "m", 'n': "n",
		'r': "4", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "p",
		'a': "6", 'e': "e_o", 'i': "i", 'o': "o", 'u': "M",
		'ä': "6", 'ö': "V", 'ü': "M", 'y': "M",
		'\u0283': "s", // sch
	},
}

// germanSpanish is best for ach-Laut (ch after back vowels) using Spanish 'x' [x].
var germanSpanish = &phoneTable{
	digraphs: map[string][]string{
		"ch": {"x"},
		"ck": {"k"},
		"ng": {"N"},
		"nk": {"N", "k"},
		"pf": {"p", "f"},
		"qu": {"k", "U"},
		"ts": {"t", "s"},
		"ss": {"s", "s"},
		"rr": {"rr"},
		"ei": {"a", "I"},
		"ie": {"i", "i"},
		"eu": {"o", "I"},
		"äu": {"o", "I"},
		"au": {"a", "U"},
	},
	expanded: map[rune][]string{
		'z': {"t", "s"},
	},
	singles: map[rune]string{
		'w': "B", 'v': "f", 'j': "I",
		'h': "x", 'l': "l", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ä': "e", 'ö': "e", 'ü': "i", 'y': "i",
		'\u0283': "sh", // sch
	},
}
