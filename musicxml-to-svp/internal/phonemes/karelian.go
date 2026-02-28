package phonemes

import "strings"

func newKarelian() *Converter {
	return &Converter{
		selectLang: selectKarelian,
		skip:       isKarelianPalatalization,
		tables: map[string]*phoneTable{
			"mandarin":  karelianMandarin,
			"cantonese": karelianCantonese,
			"japanese":  karelianJapanese,
			"english":   karelianEnglish,
			"korean":    karelianKorean,
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
	if strings.ContainsRune(word, 'r') {
		return "korean"
	}
	if hasOnlyBasicVowels(word, "äöy") {
		return "japanese"
	}
	return "english"
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

var karelianJapanese = &phoneTable{
	digraphs: map[string][]string{
		"ng": {"N", "g"},
		"ts": {"ts"},
	},
	singles: map[rune]string{
		'č': "ch", 'š': "sh", 'ž': "j",
		'h': "h", 'j': "y", 'l': "r", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f", 'v': "v", 'z': "z", 'w': "w",
		'a': "a", 'e': "e", 'i': "i", 'o': "o", 'u': "u",
		'ö': "o", 'ä': "a", 'y': "u",
	},
}

var karelianEnglish = &phoneTable{
	digraphs: map[string][]string{
		"ng": {"ng", "g"},
		"ts": {"t", "s"},
	},
	singles: map[rune]string{
		'č': "ch", 'š': "sh", 'ž': "zh",
		'h': "hh", 'j': "y", 'l': "l", 'm': "m", 'n': "n",
		'r': "r", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "f", 'v': "v", 'z': "z", 'w': "w",
		'a': "aa", 'e': "eh", 'i': "iy", 'o': "ow", 'u': "uw",
		'ä': "ae", 'ö': "er", 'y': "iy",
	},
}

var karelianKorean = &phoneTable{
	digraphs: map[string][]string{
		"ng": {"N"},
		"ts": {"ts\\_h"},
		"ss": {"s_t"},
		"tt": {"tt"},
		"pp": {"pp"},
		"kk": {"k_t"},
		"ll": {"l", "l"},
	},
	singles: map[rune]string{
		'č': "ts\\_h", 'š': "s", 'ž': "s",
		'h': "h", 'j': "j", 'l': "4", 'm': "m", 'n': "n",
		'r': "4", 's': "s", 't': "t", 'p': "p", 'k': "k",
		'b': "b", 'd': "d", 'g': "g", 'f': "p", 'v': "b", 'z': "s", 'w': "w",
		'a': "6", 'e': "e_o", 'i': "i", 'o': "o", 'u': "M",
		'ä': "6", 'ö': "V", 'y': "M",
	},
}
