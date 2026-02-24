/// <reference path="./docs/synthesizer-v-api.d.ts" />

/*

This script converts Karelian lyrics to Synthesizer V phonemes using an adaptive
multi-language approach. It analyzes each word and selects the best phoneme set
(Mandarin, Cantonese, Japanese, or English) based on the word's characteristics.

Karelian is a Finnic language closely related to Finnish and Estonian.
Key phonological features:
- Vowels: a, e, i, o, u, y [y], ä [æ], ö [ø]
- Special consonants: č [tʃ], š [ʃ], ž [ʒ]
- Palatalized consonants marked with apostrophe: l', n', s', t', d', r'
- Geminate consonants (doubled letters)
- No õ (Estonian-specific vowel)

Strategy:
- Words with y (ü) → Mandarin/Cantonese (exact 'y' phoneme match)
- Words with ö → Cantonese (best '9' approximation)
- Words with ä → Cantonese (good 'E' match)
- Simple words → Japanese (cleanest basic vowels)
- Default → English (widely compatible)

*/

const SCRIPT_TITLE = "Karelian Phonemes (Adaptive)";

function getClientInfo() {
	return {
		name: SV.T(SCRIPT_TITLE),
		author: "Egon Elbre",
		category: "Language",
		versionNumber: 1,
		minEditorVersion: 65537,
	};
}

function getTranslations(langCode) {
	return [];
}

function main() {
	var form = {
		title: SV.T(SCRIPT_TITLE),
		buttons: "OkCancel",
		widgets: [
			{
				name: "scope",
				type: "ComboBox",
				label: SV.T("Scope"),
				choices: [
					SV.T("Selected Notes"),
					SV.T("Current Track"),
					SV.T("Entire Project"),
				],
				default: hasSelectedNotes() ? 0 : 2,
			},
			{
				name: "show_language",
				type: "CheckBox",
				label: SV.T("Show language selection info"),
				default: false,
			},
		],
	};

	var result = SV.showCustomDialog(form);
	if (result.status != 1) {
		SV.finish();
		return;
	}

	var scope = result.answers.scope;
	if (scope == 0) {
		processSelection(processNotes, result.answers);
	} else if (scope == 1) {
		processTrack(processNotes, result.answers);
	} else if (scope == 2) {
		processProjectWithRefs(processNotes, result.answers);
	} else {
		SV.showMessageBox(SV.T(SCRIPT_TITLE), SV.T("Invalid scope."));
	}
	SV.finish();
}

function processNotes(notes, group, options, groupRef) {
	if (notes.length == 0) {
		return;
	}

	var languageStats = {
		mandarin: 0,
		cantonese: 0,
		japanese: 0,
		english: 0,
		korean: 0,
	};

	for (var i = 0; i < notes.length; i++) {
		var note = notes[i];
		var lyrics = note.getLyrics();

		// Skip special markers and silence
		if (
			lyrics == "-" ||
			lyrics == "+" ||
			lyrics == "sil" ||
			lyrics == "br" ||
			lyrics == "SP" ||
			lyrics == "AP"
		) {
			continue;
		}

		var wordLower = lyrics.toLowerCase();

		// Analyze word and select best language
		var selectedLang = selectBestLanguage(wordLower);
		languageStats[selectedLang]++;

		// Set language and phonemes based on selection
		note.setLanguageOverride(selectedLang);

		var phonemes = "";
		switch (selectedLang) {
			case "mandarin":
				phonemes = karelianToMandarinPhonemes(wordLower);
				break;
			case "cantonese":
				phonemes = karelianToCantonesePhonemes(wordLower);
				break;
			case "japanese":
				phonemes = karelianToJapanesePhonemes(wordLower);
				break;
			case "english":
				phonemes = karelianToEnglishPhonemes(wordLower);
				break;
			case "korean":
				phonemes = karelianToKoreanPhonemes(wordLower);
				break;
		}

		if (phonemes) {
			note.setPhonemes(phonemes);
		}
	}

	// Show statistics if requested
	if (options.show_language) {
		var msg = "Language selection statistics:\n\n";
		msg += "Mandarin: " + languageStats.mandarin + " words\n";
		msg += "Cantonese: " + languageStats.cantonese + " words\n";
		msg += "Japanese: " + languageStats.japanese + " words\n";
		msg += "English: " + languageStats.english + " words\n";
		msg += "Korean: " + languageStats.korean + " words";
		SV.showMessageBox(SV.T(SCRIPT_TITLE), msg);
	}
}

// Analyze word and select best language for phoneme conversion
function selectBestLanguage(word) {
	// Priority 1: Words with y (front rounded vowel) → Mandarin (exact match)
	if (word.indexOf("y") >= 0) {
		return "mandarin";
	}

	// Priority 2: Words with ö → Cantonese (best approximation)
	if (word.indexOf("ö") >= 0) {
		return "cantonese";
	}

	// Priority 3: Words with ä → Cantonese (good E match)
	if (word.indexOf("ä") >= 0) {
		return "cantonese";
	}

	// Priority 4: Words with r → Korean (good R match)
	if (word.indexOf("r") >= 0) {
		return "korean";
	}

	// Priority 5: Simple basic vowels only → Japanese (cleanest)
	if (hasOnlyBasicVowels(word)) {
		return "japanese";
	}

	// Default: English (most compatible)
	return "english";
}

function hasOnlyBasicVowels(word) {
	// Check if word only contains a, e, i, o, u (no special Karelian vowels)
	for (var i = 0; i < word.length; i++) {
		var c = word[i];
		if (c == "ä" || c == "ö" || c == "y") {
			return false;
		}
	}
	return true;
}

// Strip palatalization marks (apostrophes) and return cleaned character
// Karelian uses ' after consonants to mark palatalization
function isPalatalizationMark(char) {
	return (
		char == "'" || char == "\u2018" || char == "\u2019" || char == "\u02BC"
	);
}

// Mandarin phoneme conversion
function karelianToMandarinPhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var afterNext = i + 2 < word.length ? word[i + 2] : "";

		// Skip palatalization marks
		if (isPalatalizationMark(char)) {
			i++;
			continue;
		}

		if (char == "č") {
			phonemes.push("ts`");
			i++;
		} else if (char == "š") {
			phonemes.push("s`");
			i++;
		} else if (char == "ž") {
			phonemes.push("s`");
			i++;
		} else if (char == "h") {
			phonemes.push("x");
			i++;
		} else if (char == "j") {
			phonemes.push("j");
			i++;
		} else if (char == "l") {
			if (nextChar == "l") {
				phonemes.push("l", "l");
				i += 2;
			} else {
				phonemes.push("l");
				i++;
			}
		} else if (char == "m") {
			if (nextChar == "m") {
				phonemes.push("m", "m");
				i += 2;
			} else {
				phonemes.push("m");
				i++;
			}
		} else if (char == "n") {
			if (nextChar == "g") {
				phonemes.push("N");
				i += 2;
			} else if (nextChar == "n") {
				phonemes.push("n", "n");
				i += 2;
			} else {
				phonemes.push("n");
				i++;
			}
		} else if (char == "r") {
			if (nextChar == "r") {
				phonemes.push("r\\`", "r\\`");
				i += 2;
			} else {
				phonemes.push("r\\`");
				i++;
			}
		} else if (char == "s") {
			if (nextChar == "s") {
				phonemes.push("s", "s");
				i += 2;
			} else {
				phonemes.push("s");
				i++;
			}
		} else if (char == "t") {
			if (nextChar == "s") {
				phonemes.push("ts");
				i += 2;
			} else if (nextChar == "t") {
				phonemes.push("t", "t");
				i += 2;
			} else {
				phonemes.push("t");
				i++;
			}
		} else if (char == "p") {
			if (nextChar == "p") {
				phonemes.push("p", "p");
				i += 2;
			} else {
				phonemes.push("p");
				i++;
			}
		} else if (char == "k") {
			if (nextChar == "k") {
				phonemes.push("k", "k");
				i += 2;
			} else {
				phonemes.push("k");
				i++;
			}
		} else if (char == "b") {
			phonemes.push("p");
			i++;
		} else if (char == "d") {
			phonemes.push("t");
			i++;
		} else if (char == "g") {
			phonemes.push("k");
			i++;
		} else if (char == "f") {
			phonemes.push("f");
			i++;
		} else if (char == "v") {
			phonemes.push("w");
			i++;
		} else if (char == "z") {
			phonemes.push("ts");
			i++;
		} else if (char == "w") {
			phonemes.push("w");
			i++;
		} else if (char == "a") {
			if (nextChar == "a") {
				phonemes.push("a", "a");
				i += 2;
			} else {
				phonemes.push("a");
				i++;
			}
		} else if (char == "e") {
			if (nextChar == "e") {
				phonemes.push("e", "e");
				i += 2;
			} else {
				phonemes.push("e");
				i++;
			}
		} else if (char == "i") {
			if (nextChar == "i") {
				phonemes.push("i", "i");
				i += 2;
			} else {
				phonemes.push("i");
				i++;
			}
		} else if (char == "o") {
			if (nextChar == "o") {
				phonemes.push("o", "o");
				i += 2;
			} else {
				phonemes.push("o");
				i++;
			}
		} else if (char == "u") {
			if (nextChar == "u") {
				phonemes.push("u", "u");
				i += 2;
			} else {
				phonemes.push("u");
				i++;
			}
		} else if (char == "ä") {
			if (nextChar == "ä") {
				phonemes.push("A", "A");
				i += 2;
			} else {
				phonemes.push("A");
				i++;
			}
		} else if (char == "ö") {
			if (nextChar == "ö") {
				phonemes.push("@", "@");
				i += 2;
			} else {
				phonemes.push("@");
				i++;
			}
		} else if (char == "y") {
			if (nextChar == "y") {
				phonemes.push("y", "y");
				i += 2;
			} else {
				phonemes.push("y");
				i++;
			}
		} else {
			i++;
		}
	}

	return phonemes.join(" ");
}

// Cantonese phoneme conversion
function karelianToCantonesePhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";

		// Skip palatalization marks
		if (isPalatalizationMark(char)) {
			i++;
			continue;
		}

		if (char == "č") {
			phonemes.push("ts");
			i++;
		} else if (char == "š") {
			phonemes.push("s");
			i++;
		} else if (char == "ž") {
			phonemes.push("s");
			i++;
		} else if (char == "h") {
			phonemes.push("h");
			i++;
		} else if (char == "j") {
			phonemes.push("j");
			i++;
		} else if (char == "l") {
			if (nextChar == "l") {
				phonemes.push("l", "l");
				i += 2;
			} else {
				phonemes.push("l");
				i++;
			}
		} else if (char == "m") {
			if (nextChar == "m") {
				phonemes.push("m", "m");
				i += 2;
			} else {
				phonemes.push("m");
				i++;
			}
		} else if (char == "n") {
			if (nextChar == "g") {
				phonemes.push("N");
				i += 2;
			} else if (nextChar == "n") {
				phonemes.push("n", "n");
				i += 2;
			} else {
				phonemes.push("n");
				i++;
			}
		} else if (char == "r") {
			if (nextChar == "r") {
				phonemes.push("l", "l");
				i += 2;
			} else {
				phonemes.push("l");
				i++;
			}
		} else if (char == "s") {
			if (nextChar == "s") {
				phonemes.push("s", "s");
				i += 2;
			} else {
				phonemes.push("s");
				i++;
			}
		} else if (char == "t") {
			if (nextChar == "s") {
				phonemes.push("ts");
				i += 2;
			} else if (nextChar == "t") {
				phonemes.push("t", "t");
				i += 2;
			} else {
				phonemes.push("t");
				i++;
			}
		} else if (char == "p") {
			if (nextChar == "p") {
				phonemes.push("p", "p");
				i += 2;
			} else {
				phonemes.push("p");
				i++;
			}
		} else if (char == "k") {
			if (nextChar == "k") {
				phonemes.push("k", "k");
				i += 2;
			} else {
				phonemes.push("k");
				i++;
			}
		} else if (char == "b") {
			phonemes.push("p");
			i++;
		} else if (char == "d") {
			phonemes.push("t");
			i++;
		} else if (char == "g") {
			phonemes.push("k");
			i++;
		} else if (char == "f") {
			phonemes.push("f");
			i++;
		} else if (char == "v") {
			phonemes.push("w");
			i++;
		} else if (char == "z") {
			phonemes.push("ts");
			i++;
		} else if (char == "w") {
			phonemes.push("w");
			i++;
		} else if (char == "a") {
			if (nextChar == "a") {
				phonemes.push("a", "a");
				i += 2;
			} else {
				phonemes.push("a");
				i++;
			}
		} else if (char == "e") {
			if (nextChar == "e") {
				phonemes.push("e", "e");
				i += 2;
			} else {
				phonemes.push("e");
				i++;
			}
		} else if (char == "i") {
			if (nextChar == "i") {
				phonemes.push("i", "i");
				i += 2;
			} else {
				phonemes.push("i");
				i++;
			}
		} else if (char == "o") {
			if (nextChar == "o") {
				phonemes.push("o", "o");
				i += 2;
			} else {
				phonemes.push("o");
				i++;
			}
		} else if (char == "u") {
			if (nextChar == "u") {
				phonemes.push("u", "u");
				i += 2;
			} else {
				phonemes.push("u");
				i++;
			}
		} else if (char == "ä") {
			if (nextChar == "ä") {
				phonemes.push("E", "E");
				i += 2;
			} else {
				phonemes.push("E");
				i++;
			}
		} else if (char == "ö") {
			if (nextChar == "ö") {
				phonemes.push("9", "9");
				i += 2;
			} else {
				phonemes.push("9");
				i++;
			}
		} else if (char == "y") {
			if (nextChar == "y") {
				phonemes.push("y", "y");
				i += 2;
			} else {
				phonemes.push("y");
				i++;
			}
		} else {
			i++;
		}
	}

	return phonemes.join(" ");
}

// Japanese phoneme conversion
function karelianToJapanesePhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";

		// Skip palatalization marks
		if (isPalatalizationMark(char)) {
			i++;
			continue;
		}

		if (char == "č") {
			phonemes.push("ch");
			i++;
		} else if (char == "š") {
			phonemes.push("sh");
			i++;
		} else if (char == "ž") {
			phonemes.push("j");
			i++;
		} else if (char == "h") {
			phonemes.push("h");
			i++;
		} else if (char == "j") {
			phonemes.push("y");
			i++;
		} else if (char == "l") {
			if (nextChar == "l") {
				phonemes.push("r", "r");
				i += 2;
			} else {
				phonemes.push("r");
				i++;
			}
		} else if (char == "m") {
			if (nextChar == "m") {
				phonemes.push("m", "m");
				i += 2;
			} else {
				phonemes.push("m");
				i++;
			}
		} else if (char == "n") {
			if (nextChar == "g") {
				phonemes.push("N", "g");
				i += 2;
			} else if (nextChar == "n") {
				phonemes.push("n", "n");
				i += 2;
			} else {
				phonemes.push("n");
				i++;
			}
		} else if (char == "r") {
			if (nextChar == "r") {
				phonemes.push("r", "r");
				i += 2;
			} else {
				phonemes.push("r");
				i++;
			}
		} else if (char == "s") {
			if (nextChar == "s") {
				phonemes.push("s", "s");
				i += 2;
			} else {
				phonemes.push("s");
				i++;
			}
		} else if (char == "t") {
			if (nextChar == "s") {
				phonemes.push("ts");
				i += 2;
			} else if (nextChar == "t") {
				phonemes.push("t", "t");
				i += 2;
			} else {
				phonemes.push("t");
				i++;
			}
		} else if (char == "p") {
			if (nextChar == "p") {
				phonemes.push("p", "p");
				i += 2;
			} else {
				phonemes.push("p");
				i++;
			}
		} else if (char == "k") {
			if (nextChar == "k") {
				phonemes.push("k", "k");
				i += 2;
			} else {
				phonemes.push("k");
				i++;
			}
		} else if (char == "b") {
			phonemes.push("b");
			i++;
		} else if (char == "d") {
			phonemes.push("d");
			i++;
		} else if (char == "g") {
			phonemes.push("g");
			i++;
		} else if (char == "f") {
			phonemes.push("f");
			i++;
		} else if (char == "v") {
			phonemes.push("v");
			i++;
		} else if (char == "z") {
			phonemes.push("z");
			i++;
		} else if (char == "w") {
			phonemes.push("w");
			i++;
		} else if (char == "a") {
			if (nextChar == "a") {
				phonemes.push("a", "a");
				i += 2;
			} else {
				phonemes.push("a");
				i++;
			}
		} else if (char == "e") {
			if (nextChar == "e") {
				phonemes.push("e", "e");
				i += 2;
			} else {
				phonemes.push("e");
				i++;
			}
		} else if (char == "i") {
			if (nextChar == "i") {
				phonemes.push("i", "i");
				i += 2;
			} else {
				phonemes.push("i");
				i++;
			}
		} else if (char == "o") {
			if (nextChar == "o") {
				phonemes.push("o", "o");
				i += 2;
			} else {
				phonemes.push("o");
				i++;
			}
		} else if (char == "u") {
			if (nextChar == "u") {
				phonemes.push("u", "u");
				i += 2;
			} else {
				phonemes.push("u");
				i++;
			}
		} else if (char == "ö") {
			if (nextChar == "ö") {
				phonemes.push("o", "o");
				i += 2;
			} else {
				phonemes.push("o");
				i++;
			}
		} else if (char == "ä") {
			if (nextChar == "ä") {
				phonemes.push("a", "a");
				i += 2;
			} else {
				phonemes.push("a");
				i++;
			}
		} else if (char == "y") {
			if (nextChar == "y") {
				phonemes.push("u", "u");
				i += 2;
			} else {
				phonemes.push("u");
				i++;
			}
		} else {
			i++;
		}
	}

	return phonemes.join(" ");
}

// English phoneme conversion
function karelianToEnglishPhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";

		// Skip palatalization marks
		if (isPalatalizationMark(char)) {
			i++;
			continue;
		}

		if (char == "č") {
			phonemes.push("ch");
			i++;
		} else if (char == "š") {
			phonemes.push("sh");
			i++;
		} else if (char == "ž") {
			phonemes.push("zh");
			i++;
		} else if (char == "h") {
			phonemes.push("hh");
			i++;
		} else if (char == "j") {
			phonemes.push("y");
			i++;
		} else if (char == "l") {
			if (nextChar == "l") {
				phonemes.push("l", "l");
				i += 2;
			} else {
				phonemes.push("l");
				i++;
			}
		} else if (char == "m") {
			if (nextChar == "m") {
				phonemes.push("m", "m");
				i += 2;
			} else {
				phonemes.push("m");
				i++;
			}
		} else if (char == "n") {
			if (nextChar == "g") {
				phonemes.push("ng", "g");
				i += 2;
			} else if (nextChar == "n") {
				phonemes.push("n", "n");
				i += 2;
			} else {
				phonemes.push("n");
				i++;
			}
		} else if (char == "r") {
			if (nextChar == "r") {
				phonemes.push("r", "r");
				i += 2;
			} else {
				phonemes.push("r");
				i++;
			}
		} else if (char == "s") {
			if (nextChar == "s") {
				phonemes.push("s", "s");
				i += 2;
			} else {
				phonemes.push("s");
				i++;
			}
		} else if (char == "t") {
			if (nextChar == "s") {
				phonemes.push("t", "s");
				i += 2;
			} else if (nextChar == "t") {
				phonemes.push("t", "t");
				i += 2;
			} else {
				phonemes.push("t");
				i++;
			}
		} else if (char == "p") {
			if (nextChar == "p") {
				phonemes.push("p", "p");
				i += 2;
			} else {
				phonemes.push("p");
				i++;
			}
		} else if (char == "k") {
			if (nextChar == "k") {
				phonemes.push("k", "k");
				i += 2;
			} else {
				phonemes.push("k");
				i++;
			}
		} else if (char == "b") {
			phonemes.push("b");
			i++;
		} else if (char == "d") {
			phonemes.push("d");
			i++;
		} else if (char == "g") {
			phonemes.push("g");
			i++;
		} else if (char == "f") {
			phonemes.push("f");
			i++;
		} else if (char == "v") {
			phonemes.push("v");
			i++;
		} else if (char == "z") {
			phonemes.push("z");
			i++;
		} else if (char == "w") {
			phonemes.push("w");
			i++;
		} else if (char == "a") {
			if (nextChar == "a") {
				phonemes.push("aa", "aa");
				i += 2;
			} else {
				phonemes.push("aa");
				i++;
			}
		} else if (char == "e") {
			if (nextChar == "e") {
				phonemes.push("eh", "eh");
				i += 2;
			} else {
				phonemes.push("eh");
				i++;
			}
		} else if (char == "i") {
			if (nextChar == "i") {
				phonemes.push("iy", "iy");
				i += 2;
			} else {
				phonemes.push("iy");
				i++;
			}
		} else if (char == "o") {
			if (nextChar == "o") {
				phonemes.push("ow", "ow");
				i += 2;
			} else {
				phonemes.push("ow");
				i++;
			}
		} else if (char == "u") {
			if (nextChar == "u") {
				phonemes.push("uw", "uw");
				i += 2;
			} else {
				phonemes.push("uw");
				i++;
			}
		} else if (char == "ä") {
			if (nextChar == "ä") {
				phonemes.push("ae", "ae");
				i += 2;
			} else {
				phonemes.push("ae");
				i++;
			}
		} else if (char == "ö") {
			if (nextChar == "ö") {
				phonemes.push("er", "er");
				i += 2;
			} else {
				phonemes.push("er");
				i++;
			}
		} else if (char == "y") {
			if (nextChar == "y") {
				phonemes.push("iy", "iy");
				i += 2;
			} else {
				phonemes.push("iy");
				i++;
			}
		} else {
			i++;
		}
	}

	return phonemes.join(" ");
}

// Korean phoneme conversion
function karelianToKoreanPhonemes(word) {
	if (!word) return "";

	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";

		// Skip palatalization marks
		if (isPalatalizationMark(char)) {
			i++;
			continue;
		}

		// Consonants
		if (char == "č") {
			phonemes.push("ts\\_h");
			i++;
		} else if (char == "š") {
			phonemes.push("s");
			i++;
		} else if (char == "ž") {
			phonemes.push("s");
			i++;
		} else if (char == "h") {
			phonemes.push("h");
			i++;
		} else if (char == "j") {
			phonemes.push("j");
			i++;
		} else if (char == "l") {
			if (nextChar == "l") {
				phonemes.push("l", "l");
				i += 2;
			} else {
				phonemes.push("4");
				i++;
			}
		} else if (char == "m") {
			if (nextChar == "m") {
				phonemes.push("m", "m");
				i += 2;
			} else {
				phonemes.push("m");
				i++;
			}
		} else if (char == "n") {
			if (nextChar == "g") {
				phonemes.push("N");
				i += 2;
			} else if (nextChar == "n") {
				phonemes.push("n", "n");
				i += 2;
			} else {
				phonemes.push("n");
				i++;
			}
		} else if (char == "r") {
			if (nextChar == "r") {
				phonemes.push("4", "4");
				i += 2;
			} else {
				phonemes.push("4");
				i++;
			}
		} else if (char == "s") {
			if (nextChar == "s") {
				phonemes.push("s_t");
				i += 2;
			} else {
				phonemes.push("s");
				i++;
			}
		} else if (char == "t") {
			if (nextChar == "s") {
				phonemes.push("ts\\_h");
				i += 2;
			} else if (nextChar == "t") {
				phonemes.push("tt");
				i += 2;
			} else {
				phonemes.push("t");
				i++;
			}
		} else if (char == "p") {
			if (nextChar == "p") {
				phonemes.push("pp");
				i += 2;
			} else {
				phonemes.push("p");
				i++;
			}
		} else if (char == "k") {
			if (nextChar == "k") {
				phonemes.push("k_t");
				i += 2;
			} else {
				phonemes.push("k");
				i++;
			}
		} else if (char == "b") {
			phonemes.push("b");
			i++;
		} else if (char == "d") {
			phonemes.push("d");
			i++;
		} else if (char == "g") {
			phonemes.push("g");
			i++;
		} else if (char == "f") {
			phonemes.push("p");
			i++;
		} else if (char == "v") {
			phonemes.push("b");
			i++;
		} else if (char == "z") {
			phonemes.push("s");
			i++;
		} else if (char == "w") {
			phonemes.push("w");
			i++;
		}

		// Vowels
		else if (char == "a") {
			if (nextChar == "a") {
				phonemes.push("6", "6");
				i += 2;
			} else {
				phonemes.push("6");
				i++;
			}
		} else if (char == "e") {
			if (nextChar == "e") {
				phonemes.push("e_o", "e_o");
				i += 2;
			} else {
				phonemes.push("e_o");
				i++;
			}
		} else if (char == "i") {
			if (nextChar == "i") {
				phonemes.push("i", "i");
				i += 2;
			} else {
				phonemes.push("i");
				i++;
			}
		} else if (char == "o") {
			if (nextChar == "o") {
				phonemes.push("o", "o");
				i += 2;
			} else {
				phonemes.push("o");
				i++;
			}
		} else if (char == "u") {
			if (nextChar == "u") {
				phonemes.push("M", "M");
				i += 2;
			} else {
				phonemes.push("M");
				i++;
			}
		} else if (char == "ä") {
			if (nextChar == "ä") {
				phonemes.push("6", "6");
				i += 2;
			} else {
				phonemes.push("6");
				i++;
			}
		} else if (char == "ö") {
			if (nextChar == "ö") {
				phonemes.push("V", "V");
				i += 2;
			} else {
				phonemes.push("V");
				i++;
			}
		} else if (char == "y") {
			if (nextChar == "y") {
				phonemes.push("M", "M");
				i += 2;
			} else {
				phonemes.push("M");
				i++;
			}
		}

		// Skip unknown characters
		else {
			i++;
		}
	}

	return phonemes.join(" ");
}

// * Common Helper Functions * //

function hasSelectedNotes() {
	return SV.getMainEditor().getSelection().hasSelectedNotes();
}

function processSelection(process, options) {
	var selection = SV.getMainEditor().getSelection();
	var selectedNotes = selection.getSelectedNotes();
	selectedNotes = sortNotes(selectedNotes);

	var groupRef = SV.getMainEditor().getCurrentGroup();
	var group = groupRef.getTarget();
	process(selectedNotes, group, options, groupRef);
}

function processTrack(process, options) {
	var track = SV.getMainEditor().getCurrentTrack();
	var groupCount = track.getNumGroups();
	var visited = [];
	for (var i = 0; i < groupCount; i++) {
		var groupRef = track.getGroupReference(i);
		var group = groupRef.getTarget();

		// some note groups may be shared between or within tracks
		if (visited.indexOf(group.getUUID()) >= 0) continue;
		visited.push(group.getUUID());

		process(groupAsNotesArray(group), group, options, groupRef);
	}
}

function processProjectWithRefs(process, options) {
	var visited = [];
	var project = SV.getProject();

	// process unique groups for each track
	for (var i = 0; i < project.getNumTracks(); i++) {
		var track = project.getTrack(i);
		var groupCount = track.getNumGroups();
		var visited = [];
		for (var k = 0; k < groupCount; k++) {
			var groupRef = track.getGroupReference(k);
			var group = groupRef.getTarget();

			// some note groups may be shared between or within tracks
			if (visited.indexOf(group.getUUID()) >= 0) continue;
			visited.push(group.getUUID());

			process(groupAsNotesArray(group), group, options, groupRef);
		}
	}
}

function sortNotes(notes) {
	return notes.sort(function (a, b) {
		if (a.getOnset() < b.getOnset()) return -1;
		if (a.getOnset() > b.getOnset()) return 1;
		return 0;
	});
}

function groupAsNotesArray(noteGroup) {
	return new Proxy(noteGroup, {
		get: function (target, prop) {
			if (prop === "length") {
				return target.getNumNotes();
			}
			if (typeof prop == "number") {
				return target.getNote(prop);
			}
			return target[prop];
		},
	});
}
