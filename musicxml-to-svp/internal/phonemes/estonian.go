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
			"japanese":  estonianJapanese,
			"english":   estonianEnglish,
			"korean":    estonianKorean,
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
	if strings.ContainsRune(word, 'r') {
		return "korean"
	}
	if hasOnlyBasicVowels(word, "äöüõy") {
		return "japanese"
	}
	return "english"
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

var estonianJapanese = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"sh"},
		"ng": {"N", "g"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'h': "h", 'j': "y", 'l': "r", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f", 'v': "v", 'z': "z", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'õ': "o", 'ä': "a", 'ö': "o", 'y': "u",
	},
}

var estonianEnglish = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"sh"},
		"ng": {"ng", "g"},
		"ts": {"t", "s"},
	},
	singles: map[rune]string{
		'h': "hh", 'j': "y", 'l': "l", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f", 'v': "v", 'z': "z", 'w': "w",
		'a': "aa", 'e': "eh", 'i': "iy", 'o': "ow", 'u': "uw",
		'õ': "uh", 'ä': "ae", 'ö': "er", 'y': "iy",
	},
}

var estonianKorean = &phoneTable{
	digraphs: map[string][]string{
		"sh": {"s"},
		"ng": {"N"},
		"ts": {"ts\\_h"},
		"ss": {"s_t"},
		"tt": {"tt"},
		"pp": {"pp"},
		"kk": {"k_t"},
		"ll": {"l", "l"},
	},
	singles: map[rune]string{
		'h': "h", 'j': "j", 'l': "4", 'm': "m", 'n': "n",
		'r': "4", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "p", 'v': "b", 'z': "s", 'w': "w",
		'a': "6", 'e': "e_o", 'i': "i", 'o': "o", 'u': "M",
		'õ': "V", 'ä': "6", 'ö': "V", 'y': "M",
	},
}
