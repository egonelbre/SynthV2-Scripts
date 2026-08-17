package phonemes

import "strings"

func newEstonian() *Converter {
	return &Converter{
		selectLangs: selectEstonian,
		normalize:   func(s string) string { return strings.ReplaceAll(s, "ü", "y") },
		vowels:      "aeiouyõäöü",
		tables: map[string]*phoneTable{
			"mandarin":  estonianMandarin,
			"cantonese": estonianCantonese,
			"spanish":   estonianSpanish,
		},
		words: map[string]Result{
			"h(m)": Result{
				Language: "cantonese",
				Phoneset: "xsampa",
				Phonemes: "h m=",
			},
			"hm": Result{
				Language: "cantonese",
				Phoneset: "xsampa",
				Phonemes: "h m=",
			},
			"hmm": Result{
				Language: "cantonese",
				Phoneset: "xsampa",
				Phonemes: "h m=",
			},
			"ŋ": Result{
				Language: "english",
				Phoneset: "arpabet",
				Phonemes: "ng",
			},
		},
	}
}

// selectEstonian returns the languages that can sing a word, best first.
// Front rounded vowels pin the choice; plain words can go to any of the three,
// except that Mandarin and Cantonese mangle r and v.
func selectEstonian(word string) []string {
	switch {
	case strings.ContainsAny(word, "üy"):
		return []string{"mandarin"} // exact 'y'
	case strings.ContainsRune(word, 'ö'):
		return []string{"cantonese"} // '9'
	case strings.ContainsRune(word, 'õ'):
		return []string{"mandarin"} // '7'
	case strings.ContainsRune(word, 'ä'):
		return []string{"cantonese"} // 'E'
	case strings.ContainsAny(word, "rvz"):
		return []string{"spanish"}
	}
	return []string{"spanish", "mandarin", "cantonese"}
}

var estonianMandarin = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"s`"},
		"ng": {"N", "k"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'h': "x", 'j': "j", 'l': "l", 'm': "m", 'n': "n",
		'r': "r\\`", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f", 'v': "w", 'z': "s", 'w': "w",
		'š': "s`", 'ž': "s`",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'õ': "7", 'ä': "A", 'ö': "@", 'y': "y",
	},
}

var estonianCantonese = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"s"},
		"ng": {"N", "k"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'h': "h", 'j': "j", 'l': "l", 'm': "m", 'n': "n",
		'r': "l", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f", 'v': "w", 'z': "s", 'w': "w",
		'š': "s", 'ž': "s",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'õ': "8", 'ä': "E", 'ö': "9", 'y': "y",
	},
}

var estonianSpanish = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"sh"},
		"ng": {"N", "k"},
		"ts": {"t", "s"},
		"rr": {"rr"}, // Estonian geminate r is a trill
	},
	singles: map[rune]string{
		'h': "x", 'j': "I", 'l': "l", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f", 'v': "B", 'z': "s", 'w': "U",
		'š': "sh", 'ž': "sh",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'õ': "o", 'ä': "a", 'ö': "e", 'y': "u",
	},
}
