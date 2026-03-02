package phonemes

import "strings"

func newKarelian() *Converter {
	return &Converter{
		selectLang: selectKarelian,
		skip:       isKarelianPalatalization,
		vowels:     "aeiouyäöü",
		tables: map[string]*phoneTable{
			"mandarin":  karelianMandarin,
			"cantonese": karelianCantonese,
			"spanish":   karelianSpanish,
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

func selectKarelian(word string) string {
	if strings.ContainsRune(word, 'y') {
		return "mandarin"
	}
	if strings.ContainsRune(word, 'ö') {
		return "cantonese"
	}
	if strings.ContainsRune(word, 'ä') {
		return "cantonese"
	}
	return "spanish"
}

func isKarelianPalatalization(r rune) bool {
	return r == '\'' || r == '\u2018' || r == '\u2019' || r == '\u02BC'
}

var karelianMandarin = &phoneTable{
	digraphs: map[string][]string{
		"ng": {"N"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'č': "ts`", 'š': "s`", 'ž': "s`",
		'h': "x", 'j': "j", 'l': "l", 'm': "m", 'n': "n",
		'r': "r\\`", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f", 'v': "w", 'z': "ts", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ä': "A", 'ö': "@", 'y': "y",
	},
}

var karelianCantonese = &phoneTable{
	digraphs: map[string][]string{
		"ng": {"N"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'č': "ts", 'š': "s", 'ž': "s",
		'h': "h", 'j': "j", 'l': "l", 'm': "m", 'n': "n",
		'r': "l", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "p", 'd': "t", 'g': "k", 'f': "f", 'v': "w", 'z': "ts", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ä': "E", 'ö': "9", 'y': "y",
	},
}

var karelianSpanish = &phoneTable{
	digraphs: map[string][]string{
		"ng": {"N"},
		"ts": {"t", "s"},
		"rr": {"rr"}, // Karelian geminate r is a trill
	},
	singles: map[rune]string{
		'č': "ch", 'š': "sh", 'ž': "sh",
		'h': "x", 'j': "I", 'l': "l", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f", 'v': "B", 'z': "s", 'w': "U",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ä': "a", 'ö': "e", 'y': "u",
	},
}
