package phonemes

import "strings"

// newLatvian converts Latvian. Vowel length is carried by the note, so
// macrons are dropped; a plain o is the native diphthong [uo], while ō is
// always the loanword [o:] (write ō, or use -p, for loanwords like "opera").
func newLatvian() *Converter {
	return &Converter{
		selectLangs: selectLatvian,
		normalize: strings.NewReplacer(
			"ā", "a", "ē", "e", "ī", "i", "ū", "u", "ŗ", "r",
		).Replace,
		vowels: "aeiou",
		tables: map[string]*phoneTable{
			"spanish":   latvianSpanish,
			"mandarin":  latvianMandarin,
			"cantonese": latvianCantonese,
		},
	}
}

// selectLatvian returns the languages that can sing a word, best first.
// Spanish has the trill, [v]-ish B, and palatal J; Mandarin and Cantonese
// mangle those, so only plain words may go to them.
func selectLatvian(word string) []string {
	if strings.ContainsAny(word, "rŗvzņļ") {
		return []string{"spanish"}
	}
	return []string{"spanish", "mandarin", "cantonese"}
}

// Palatal ķ/ģ/ļ have no Spanish counterpart and fall back to the plain
// stop/lateral.
var latvianSpanish = &phoneTable{
	digraphs: map[string][]string{
		"dz": {"d", "s"},
		"dž": {"d", "y"},
		"rr": {"rr"},
	},
	expanded: map[rune][]string{
		'o': {"u", "o"},
	},
	singles: map[rune]string{
		'b': "b", 'c': "ts", 'č': "ch", 'd': "d", 'f': "f", 'g': "g", 'ģ': "g",
		'h': "x", 'j': "I", 'k': "k", 'ķ': "k", 'l': "l", 'ļ': "l", 'm': "m",
		'n': "n", 'ņ': "J", 'p': "p", 'r': "r", 's': "s", 'š': "sh", 't': "t",
		'v': "B", 'z': "s", 'ž': "sh", 'w': "U",
		'a': "a", 'e': "e", 'i': "i", 'ō': "o", 'u': "u",
	},
}

var latvianMandarin = &phoneTable{
	digraphs: map[string][]string{
		"dz": {"ts"},
		"dž": {"ts`"},
	},
	expanded: map[rune][]string{
		'o': {"u", "o"},
	},
	singles: map[rune]string{
		'b': "p", 'c': "ts", 'č': "ts`", 'd': "t", 'f': "f", 'g': "k", 'ģ': "k",
		'h': "x", 'j': "j", 'k': "k", 'ķ': "k", 'l': "l", 'ļ': "l", 'm': "m",
		'n': "n", 'ņ': "n", 'p': "p", 'r': "r\\`", 's': "s", 'š': "s`", 't': "t",
		'v': "w", 'z': "s", 'ž': "s`", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'ō': "o", 'u': "u",
	},
}

var latvianCantonese = &phoneTable{
	digraphs: map[string][]string{
		"dz": {"ts"},
		"dž": {"ts"},
	},
	expanded: map[rune][]string{
		'o': {"u", "o"},
	},
	singles: map[rune]string{
		'b': "p", 'c': "ts", 'č': "ts", 'd': "t", 'f': "f", 'g': "k", 'ģ': "k",
		'h': "h", 'j': "j", 'k': "k", 'ķ': "k", 'l': "l", 'ļ': "l", 'm': "m",
		'n': "n", 'ņ': "n", 'p': "p", 'r': "l", 's': "s", 'š': "s", 't': "t",
		'v': "w", 'z': "s", 'ž': "s", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'ō': "o", 'u': "u",
	},
}
