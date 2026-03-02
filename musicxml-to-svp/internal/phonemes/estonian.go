package phonemes

import "strings"

func newEstonian() *Converter {
	return &Converter{
		selectLang: selectEstonian,
		normalize:  func(s string) string { return strings.ReplaceAll(s, "ü", "y") },
		vowels:     "aeiouyõäöü",
		tables: map[string]*phoneTable{
			"mandarin":  estonianMandarin,
			"cantonese": estonianCantonese,
			"spanish":   estonianSpanish,
		},
		words: map[string]Result{
			"hmm": Result{
				Language: "cantonese",
				Phoneset: "xsampa",
				Phonemes: "h m=",
			},
		},
	}
}

func selectEstonian(word string) string {
	if strings.ContainsAny(word, "üy") {
		return "mandarin"
	}
	if strings.ContainsRune(word, 'ö') {
		return "cantonese"
	}
	if strings.ContainsRune(word, 'õ') {
		return "mandarin"
	}
	if strings.ContainsRune(word, 'ä') {
		return "cantonese"
	}
	return "spanish"
}

var estonianMandarin = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"s`"},
		"ng": {"N"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'h': "x", 'j': "j", 'l': "l", 'm': "m", 'n': "n",
		'r': "r\\`", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f", 'v': "w", 'z': "ts", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'õ': "7", 'ä': "A", 'ö': "@", 'y': "y",
	},
}

var estonianCantonese = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"s"},
		"ng": {"N"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'h': "h", 'j': "j", 'l': "l", 'm': "m", 'n': "n",
		'r': "l", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f", 'v': "w", 'z': "ts", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'õ': "8", 'ä': "E", 'ö': "9", 'y': "y",
	},
}

var estonianSpanish = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"sh"},
		"ng": {"N"},
		"ts": {"t", "s"},
		"rr": {"rr"}, // Estonian geminate r is a trill
	},
	singles: map[rune]string{
		'h': "x", 'j': "I", 'l': "l", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f", 'v': "B", 'z': "s", 'w': "U",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'õ': "o", 'ä': "a", 'ö': "e", 'y': "u",
	},
}
