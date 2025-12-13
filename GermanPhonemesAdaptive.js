/// <reference path="./docs/synthesizer-v-api.d.ts" />

/*

This script converts German lyrics to Synthesizer V phonemes using an adaptive
multi-language approach. It analyzes each word and selects the best phoneme set
(Mandarin, Cantonese, Japanese, Korean, Spanish, or English) based on the word's characteristics.

Strategy:
- Words with ü → Mandarin (exact 'y' phoneme for front rounded vowel)
- Words with ö → Cantonese (best '9' approximation for front rounded vowel)
- Words with ch after a/o/u → Spanish (good 'x' for [x] sound)
- Words with ch after i/e/ä/ö/ü or initial → Japanese (approximation with 'h')
- Words with r → Korean (good '4' flap approximation)
- Words with ß → English (good 's' match)
- Simple words (only a, e, i, o, u) → Japanese (cleanest basic vowels)
- Default → English (most compatible)

German phonetics covered:
- Umlauts: ä, ö, ü
- Special consonants: ch (ich-Laut vs ach-Laut), sch, ß, z, w, v
- Consonant clusters: pf, qu, ng, nk
- German diphthongs: ei, ie, eu, äu, au

*/

const SCRIPT_TITLE = "German Phonemes (Adaptive)";

function getClientInfo() {
	return {
		"name": SV.T(SCRIPT_TITLE),
		"author": "Egon Elbre",
		"category": "German",
		"versionNumber": 1,
		"minEditorVersion": 65537
	};
}

function getTranslations(langCode) {
	return [];
}

function main() {
	var form = {
		"title": SV.T(SCRIPT_TITLE),
		"buttons": "OkCancel",
		"widgets": [
			{
				"name": "scope",
				"type": "ComboBox",
				"label": SV.T("Scope"),
				"choices": [
					SV.T("Selected Notes"),
					SV.T("Current Track"),
					SV.T("Entire Project")
				],
				"default": hasSelectedNotes() ? 0 : 2
			},
			{
				"name": "show_language",
				"type": "CheckBox",
				"label": SV.T("Show language selection info"),
				"default": false
			}
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
		spanish: 0
	};

	for (var i = 0; i < notes.length; i++) {
		var note = notes[i];
		var lyrics = note.getLyrics();

		// Skip special markers and silence
		if (lyrics == "-" || lyrics == "+" || lyrics == "sil" || lyrics == "br" || lyrics == "SP" || lyrics == "AP") {
			continue;
		}

		var wordLower = lyrics.toLowerCase();

		// Analyze word and select best language
		var selectedLang = selectBestLanguage(wordLower);
		languageStats[selectedLang]++;

		// Set language and phonemes based on selection
		note.setLanguageOverride(selectedLang);

		var phonemes = "";
		switch(selectedLang) {
			case "mandarin":
				phonemes = germanToMandarinPhonemes(wordLower);
				break;
			case "cantonese":
				phonemes = germanToCantonesePhonemes(wordLower);
				break;
			case "japanese":
				phonemes = germanToJapanesePhonemes(wordLower);
				break;
			case "english":
				phonemes = germanToEnglishPhonemes(wordLower);
				break;
			case "korean":
				phonemes = germanToKoreanPhonemes(wordLower);
				break;
			case "spanish":
				phonemes = germanToSpanishPhonemes(wordLower);
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
		msg += "Korean: " + languageStats.korean + " words\n";
		msg += "Spanish: " + languageStats.spanish + " words";
		SV.showMessageBox(SV.T(SCRIPT_TITLE), msg);
	}
}

// Analyze word and select best language for phoneme conversion
function selectBestLanguage(word) {
	// Priority 1: Words with ü → Mandarin (exact 'y' phoneme match)
	if (word.indexOf('ü') >= 0) {
		return "mandarin";
	}

	// Priority 2: Words with ö → Cantonese (best '9' approximation)
	if (word.indexOf('ö') >= 0) {
		return "cantonese";
	}

	// Priority 3: Check for "ch" sound type
	// ach-Laut (after a, o, u, au) → Spanish x [x]
	// ich-Laut (after i, e, ä, ö, ü, consonants, or initial) → Japanese h
	if (word.indexOf('ch') >= 0) {
		var chIndex = word.indexOf('ch');
		if (chIndex > 0) {
			var prevChar = word[chIndex - 1];
			// Check for ach-Laut contexts
			if (prevChar == 'a' || prevChar == 'o' || prevChar == 'u') {
				return "spanish";
			}
			// Check for "au" before ch
			if (chIndex > 1 && word[chIndex - 2] == 'a' && prevChar == 'u') {
				return "spanish";
			}
		}
		// ich-Laut or initial ch → Japanese for softer sound
		return "japanese";
	}

	// Priority 4: Words with ä → Cantonese (good 'E' match)
	if (word.indexOf('ä') >= 0) {
		return "cantonese";
	}

	// Priority 5: Words with r → Korean (good flap/tap R)
	if (word.indexOf('r') >= 0) {
		return "korean";
	}

	// Priority 6: Simple basic vowels only → Japanese (cleanest)
	if (hasOnlyBasicVowels(word)) {
		return "japanese";
	}

	// Default: English (most compatible)
	return "english";
}

function hasOnlyBasicVowels(word) {
	// Check if word only contains basic vowels (no German umlauts)
	for (var i = 0; i < word.length; i++) {
		var c = word[i];
		if (c == 'ä' || c == 'ö' || c == 'ü') {
			return false;
		}
	}
	return true;
}

// ==================== MANDARIN PHONEME CONVERSION ====================
// Best for: ü (exact 'y' match)
function germanToMandarinPhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var nextNextChar = i + 2 < word.length ? word[i + 2] : "";

		// German consonant clusters and special consonants
		if (char == 's' && nextChar == 'c' && nextNextChar == 'h') {
			// sch → s` (sh sound)
			phonemes.push('s`');
			i += 3;
		} else if (char == 'c' && nextChar == 'h') {
			// ch → x (aspirated)
			phonemes.push('x');
			i += 2;
		} else if (char == 'n' && nextChar == 'g') {
			phonemes.push('N');
			i += 2;
		} else if (char == 'n' && nextChar == 'k') {
			phonemes.push('N', 'k');
			i += 2;
		} else if (char == 'p' && nextChar == 'f') {
			phonemes.push('p', 'f');
			i += 2;
		} else if (char == 'q' && nextChar == 'u') {
			phonemes.push('k', 'w');
			i += 2;
		} else if (char == 't' && nextChar == 's') {
			phonemes.push('ts');
			i += 2;
		} else if (char == 's' && nextChar == 't') {
			// st at beginning often pronounced sht
			phonemes.push('s`', 't');
			i += 2;
		} else if (char == 's' && nextChar == 'p') {
			// sp at beginning often pronounced shp
			phonemes.push('s`', 'p');
			i += 2;
		} else if (char == 's' && nextChar == 's') {
			phonemes.push('s', 's');
			i += 2;
		} else if (char == 'ß') {
			phonemes.push('s');
			i++;
		} else if (char == 'z') {
			// German z = ts
			phonemes.push('ts');
			i++;
		} else if (char == 'w') {
			// German w = v sound
			phonemes.push('w');
			i++;
		} else if (char == 'v') {
			// German v often = f, sometimes v
			phonemes.push('f');
			i++;
		} else if (char == 'j') {
			phonemes.push('j');
			i++;
		} else if (char == 'h') {
			phonemes.push('x');
			i++;
		} else if (char == 'l') {
			phonemes.push('l');
			i++;
		} else if (char == 'm') {
			phonemes.push('m');
			i++;
		} else if (char == 'n') {
			phonemes.push('n');
			i++;
		} else if (char == 'r') {
			phonemes.push('r\\`');
			i++;
		} else if (char == 's') {
			phonemes.push('s');
			i++;
		} else if (char == 't') {
			phonemes.push('t');
			i++;
		} else if (char == 'p') {
			phonemes.push('p');
			i++;
		} else if (char == 'k') {
			phonemes.push('k');
			i++;
		} else if (char == 'b') {
			phonemes.push('p');
			i++;
		} else if (char == 'd') {
			phonemes.push('t');
			i++;
		} else if (char == 'g') {
			phonemes.push('k');
			i++;
		} else if (char == 'f') {
			phonemes.push('f');
			i++;
		}
		// German diphthongs
		else if (char == 'e' && nextChar == 'i') {
			// ei → ai
			phonemes.push('a', 'i');
			i += 2;
		} else if (char == 'i' && nextChar == 'e') {
			// ie → long i
			phonemes.push('i', 'i');
			i += 2;
		} else if (char == 'e' && nextChar == 'u') {
			// eu → oi
			phonemes.push('o', 'i');
			i += 2;
		} else if (char == 'ä' && nextChar == 'u') {
			// äu → oi
			phonemes.push('o', 'i');
			i += 2;
		} else if (char == 'a' && nextChar == 'u') {
			// au → au
			phonemes.push('a', 'u');
			i += 2;
		}
		// German vowels
		else if (char == 'a') {
			phonemes.push('a');
			i++;
		} else if (char == 'e') {
			phonemes.push('e');
			i++;
		} else if (char == 'i') {
			phonemes.push('i');
			i++;
		} else if (char == 'o') {
			phonemes.push('o');
			i++;
		} else if (char == 'u') {
			phonemes.push('u');
			i++;
		} else if (char == 'ä') {
			phonemes.push('A');
			i++;
		} else if (char == 'ö') {
			phonemes.push('@');
			i++;
		} else if (char == 'ü') {
			// Exact match - front rounded vowel
			phonemes.push('y');
			i++;
		} else if (char == 'y') {
			phonemes.push('y');
			i++;
		} else {
			i++;
		}
	}

	return phonemes.join(' ');
}

// ==================== CANTONESE PHONEME CONVERSION ====================
// Best for: ö (best '9' approximation), ä (good 'E' match)
function germanToCantonesePhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var nextNextChar = i + 2 < word.length ? word[i + 2] : "";

		// German consonant clusters
		if (char == 's' && nextChar == 'c' && nextNextChar == 'h') {
			phonemes.push('s');
			i += 3;
		} else if (char == 'c' && nextChar == 'h') {
			phonemes.push('h');
			i += 2;
		} else if (char == 'n' && nextChar == 'g') {
			phonemes.push('N');
			i += 2;
		} else if (char == 'n' && nextChar == 'k') {
			phonemes.push('N', 'k');
			i += 2;
		} else if (char == 'p' && nextChar == 'f') {
			phonemes.push('p', 'f');
			i += 2;
		} else if (char == 'q' && nextChar == 'u') {
			phonemes.push('k', 'w');
			i += 2;
		} else if (char == 't' && nextChar == 's') {
			phonemes.push('ts');
			i += 2;
		} else if (char == 's' && nextChar == 's') {
			phonemes.push('s', 's');
			i += 2;
		} else if (char == 'ß') {
			phonemes.push('s');
			i++;
		} else if (char == 'z') {
			phonemes.push('ts');
			i++;
		} else if (char == 'w') {
			phonemes.push('w');
			i++;
		} else if (char == 'v') {
			phonemes.push('f');
			i++;
		} else if (char == 'j') {
			phonemes.push('j');
			i++;
		} else if (char == 'h') {
			phonemes.push('h');
			i++;
		} else if (char == 'l') {
			phonemes.push('l');
			i++;
		} else if (char == 'm') {
			phonemes.push('m');
			i++;
		} else if (char == 'n') {
			phonemes.push('n');
			i++;
		} else if (char == 'r') {
			phonemes.push('l');
			i++;
		} else if (char == 's') {
			phonemes.push('s');
			i++;
		} else if (char == 't') {
			phonemes.push('t');
			i++;
		} else if (char == 'p') {
			phonemes.push('p');
			i++;
		} else if (char == 'k') {
			phonemes.push('k');
			i++;
		} else if (char == 'b') {
			phonemes.push('p');
			i++;
		} else if (char == 'd') {
			phonemes.push('t');
			i++;
		} else if (char == 'g') {
			phonemes.push('k');
			i++;
		} else if (char == 'f') {
			phonemes.push('f');
			i++;
		}
		// German diphthongs
		else if (char == 'e' && nextChar == 'i') {
			phonemes.push('6', 'i');
			i += 2;
		} else if (char == 'i' && nextChar == 'e') {
			phonemes.push('i', 'i');
			i += 2;
		} else if (char == 'e' && nextChar == 'u') {
			phonemes.push('O', 'i');
			i += 2;
		} else if (char == 'ä' && nextChar == 'u') {
			phonemes.push('O', 'i');
			i += 2;
		} else if (char == 'a' && nextChar == 'u') {
			phonemes.push('a', 'u');
			i += 2;
		}
		// German vowels
		else if (char == 'a') {
			phonemes.push('a');
			i++;
		} else if (char == 'e') {
			phonemes.push('e');
			i++;
		} else if (char == 'i') {
			phonemes.push('i');
			i++;
		} else if (char == 'o') {
			phonemes.push('o');
			i++;
		} else if (char == 'u') {
			phonemes.push('u');
			i++;
		} else if (char == 'ä') {
			// Good match - open front vowel
			phonemes.push('E');
			i++;
		} else if (char == 'ö') {
			// Best match - front rounded vowel
			phonemes.push('9');
			i++;
		} else if (char == 'ü') {
			phonemes.push('y');
			i++;
		} else if (char == 'y') {
			phonemes.push('y');
			i++;
		} else {
			i++;
		}
	}

	return phonemes.join(' ');
}

// ==================== JAPANESE PHONEME CONVERSION ====================
// Best for: simple words, ich-Laut (ch after front vowels)
function germanToJapanesePhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var nextNextChar = i + 2 < word.length ? word[i + 2] : "";

		// German consonant clusters
		if (char == 's' && nextChar == 'c' && nextNextChar == 'h') {
			phonemes.push('sh');
			i += 3;
		} else if (char == 'c' && nextChar == 'h') {
			// ich-Laut approximation
			phonemes.push('h');
			i += 2;
		} else if (char == 'n' && nextChar == 'g') {
			phonemes.push('N', 'g');
			i += 2;
		} else if (char == 'n' && nextChar == 'k') {
			phonemes.push('N', 'k');
			i += 2;
		} else if (char == 'p' && nextChar == 'f') {
			phonemes.push('p', 'f');
			i += 2;
		} else if (char == 'q' && nextChar == 'u') {
			phonemes.push('k', 'w');
			i += 2;
		} else if (char == 't' && nextChar == 's') {
			phonemes.push('ts');
			i += 2;
		} else if (char == 's' && nextChar == 's') {
			phonemes.push('s', 's');
			i += 2;
		} else if (char == 'ß') {
			phonemes.push('s');
			i++;
		} else if (char == 'z') {
			phonemes.push('ts');
			i++;
		} else if (char == 'w') {
			phonemes.push('v');
			i++;
		} else if (char == 'v') {
			phonemes.push('f');
			i++;
		} else if (char == 'j') {
			phonemes.push('y');
			i++;
		} else if (char == 'h') {
			phonemes.push('h');
			i++;
		} else if (char == 'l') {
			phonemes.push('r');
			i++;
		} else if (char == 'm') {
			phonemes.push('m');
			i++;
		} else if (char == 'n') {
			phonemes.push('n');
			i++;
		} else if (char == 'r') {
			phonemes.push('r');
			i++;
		} else if (char == 's') {
			phonemes.push('s');
			i++;
		} else if (char == 't') {
			phonemes.push('t');
			i++;
		} else if (char == 'p') {
			phonemes.push('p');
			i++;
		} else if (char == 'k') {
			phonemes.push('k');
			i++;
		} else if (char == 'b') {
			phonemes.push('b');
			i++;
		} else if (char == 'd') {
			phonemes.push('d');
			i++;
		} else if (char == 'g') {
			phonemes.push('g');
			i++;
		} else if (char == 'f') {
			phonemes.push('f');
			i++;
		}
		// German diphthongs
		else if (char == 'e' && nextChar == 'i') {
			phonemes.push('a', 'i');
			i += 2;
		} else if (char == 'i' && nextChar == 'e') {
			phonemes.push('i', 'i');
			i += 2;
		} else if (char == 'e' && nextChar == 'u') {
			phonemes.push('o', 'i');
			i += 2;
		} else if (char == 'ä' && nextChar == 'u') {
			phonemes.push('o', 'i');
			i += 2;
		} else if (char == 'a' && nextChar == 'u') {
			phonemes.push('a', 'u');
			i += 2;
		}
		// German vowels
		else if (char == 'a') {
			phonemes.push('a');
			i++;
		} else if (char == 'e') {
			phonemes.push('e');
			i++;
		} else if (char == 'i') {
			phonemes.push('i');
			i++;
		} else if (char == 'o') {
			phonemes.push('o');
			i++;
		} else if (char == 'u') {
			phonemes.push('u');
			i++;
		} else if (char == 'ä') {
			phonemes.push('a');
			i++;
		} else if (char == 'ö') {
			phonemes.push('o');
			i++;
		} else if (char == 'ü') {
			phonemes.push('u');
			i++;
		} else if (char == 'y') {
			phonemes.push('u');
			i++;
		} else {
			i++;
		}
	}

	return phonemes.join(' ');
}

// ==================== ENGLISH PHONEME CONVERSION ====================
// Default fallback - most compatible
function germanToEnglishPhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var nextNextChar = i + 2 < word.length ? word[i + 2] : "";

		// German consonant clusters
		if (char == 's' && nextChar == 'c' && nextNextChar == 'h') {
			phonemes.push('sh');
			i += 3;
		} else if (char == 'c' && nextChar == 'h') {
			phonemes.push('hh');
			i += 2;
		} else if (char == 'n' && nextChar == 'g') {
			phonemes.push('ng');
			i += 2;
		} else if (char == 'n' && nextChar == 'k') {
			phonemes.push('ng', 'k');
			i += 2;
		} else if (char == 'p' && nextChar == 'f') {
			phonemes.push('p', 'f');
			i += 2;
		} else if (char == 'q' && nextChar == 'u') {
			phonemes.push('k', 'w');
			i += 2;
		} else if (char == 't' && nextChar == 's') {
			phonemes.push('t', 's');
			i += 2;
		} else if (char == 's' && nextChar == 's') {
			phonemes.push('s', 's');
			i += 2;
		} else if (char == 'ß') {
			phonemes.push('s');
			i++;
		} else if (char == 'z') {
			phonemes.push('t', 's');
			i++;
		} else if (char == 'w') {
			phonemes.push('v');
			i++;
		} else if (char == 'v') {
			phonemes.push('f');
			i++;
		} else if (char == 'j') {
			phonemes.push('y');
			i++;
		} else if (char == 'h') {
			phonemes.push('hh');
			i++;
		} else if (char == 'l') {
			phonemes.push('l');
			i++;
		} else if (char == 'm') {
			phonemes.push('m');
			i++;
		} else if (char == 'n') {
			phonemes.push('n');
			i++;
		} else if (char == 'r') {
			phonemes.push('r');
			i++;
		} else if (char == 's') {
			phonemes.push('s');
			i++;
		} else if (char == 't') {
			phonemes.push('t');
			i++;
		} else if (char == 'p') {
			phonemes.push('p');
			i++;
		} else if (char == 'k') {
			phonemes.push('k');
			i++;
		} else if (char == 'b') {
			phonemes.push('b');
			i++;
		} else if (char == 'd') {
			phonemes.push('d');
			i++;
		} else if (char == 'g') {
			phonemes.push('g');
			i++;
		} else if (char == 'f') {
			phonemes.push('f');
			i++;
		}
		// German diphthongs
		else if (char == 'e' && nextChar == 'i') {
			phonemes.push('ay');
			i += 2;
		} else if (char == 'i' && nextChar == 'e') {
			phonemes.push('iy', 'iy');
			i += 2;
		} else if (char == 'e' && nextChar == 'u') {
			phonemes.push('oy');
			i += 2;
		} else if (char == 'ä' && nextChar == 'u') {
			phonemes.push('oy');
			i += 2;
		} else if (char == 'a' && nextChar == 'u') {
			phonemes.push('aw');
			i += 2;
		}
		// German vowels
		else if (char == 'a') {
			phonemes.push('aa');
			i++;
		} else if (char == 'e') {
			phonemes.push('eh');
			i++;
		} else if (char == 'i') {
			phonemes.push('iy');
			i++;
		} else if (char == 'o') {
			phonemes.push('ow');
			i++;
		} else if (char == 'u') {
			phonemes.push('uw');
			i++;
		} else if (char == 'ä') {
			phonemes.push('ae');
			i++;
		} else if (char == 'ö') {
			phonemes.push('er');
			i++;
		} else if (char == 'ü') {
			phonemes.push('iy');
			i++;
		} else if (char == 'y') {
			phonemes.push('iy');
			i++;
		} else {
			i++;
		}
	}

	return phonemes.join(' ');
}

// ==================== KOREAN PHONEME CONVERSION ====================
// Best for: words with 'r' (good flap/tap sound)
function germanToKoreanPhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var nextNextChar = i + 2 < word.length ? word[i + 2] : "";

		// German consonant clusters
		if (char == 's' && nextChar == 'c' && nextNextChar == 'h') {
			phonemes.push('s');
			i += 3;
		} else if (char == 'c' && nextChar == 'h') {
			phonemes.push('h');
			i += 2;
		} else if (char == 'n' && nextChar == 'g') {
			phonemes.push('N');
			i += 2;
		} else if (char == 'n' && nextChar == 'k') {
			phonemes.push('N', 'k');
			i += 2;
		} else if (char == 'p' && nextChar == 'f') {
			phonemes.push('p', 'p');
			i += 2;
		} else if (char == 'q' && nextChar == 'u') {
			phonemes.push('k', 'w');
			i += 2;
		} else if (char == 't' && nextChar == 's') {
			phonemes.push('ts\\_h');
			i += 2;
		} else if (char == 's' && nextChar == 's') {
			phonemes.push('s_t');
			i += 2;
		} else if (char == 'ß') {
			phonemes.push('s');
			i++;
		} else if (char == 'z') {
			phonemes.push('ts\\_h');
			i++;
		} else if (char == 'w') {
			phonemes.push('b');
			i++;
		} else if (char == 'v') {
			phonemes.push('p');
			i++;
		} else if (char == 'j') {
			phonemes.push('j');
			i++;
		} else if (char == 'h') {
			phonemes.push('h');
			i++;
		} else if (char == 'l') {
			phonemes.push('l');
			i++;
		} else if (char == 'm') {
			phonemes.push('m');
			i++;
		} else if (char == 'n') {
			phonemes.push('n');
			i++;
		} else if (char == 'r') {
			// Korean flap - best match for German r
			phonemes.push('4');
			i++;
		} else if (char == 's') {
			phonemes.push('s');
			i++;
		} else if (char == 't') {
			phonemes.push('t');
			i++;
		} else if (char == 'p') {
			phonemes.push('p');
			i++;
		} else if (char == 'k') {
			phonemes.push('k');
			i++;
		} else if (char == 'b') {
			phonemes.push('b');
			i++;
		} else if (char == 'd') {
			phonemes.push('d');
			i++;
		} else if (char == 'g') {
			phonemes.push('g');
			i++;
		} else if (char == 'f') {
			phonemes.push('p');
			i++;
		}
		// German diphthongs
		else if (char == 'e' && nextChar == 'i') {
			phonemes.push('6', 'i');
			i += 2;
		} else if (char == 'i' && nextChar == 'e') {
			phonemes.push('i', 'i');
			i += 2;
		} else if (char == 'e' && nextChar == 'u') {
			phonemes.push('o', 'i');
			i += 2;
		} else if (char == 'ä' && nextChar == 'u') {
			phonemes.push('o', 'i');
			i += 2;
		} else if (char == 'a' && nextChar == 'u') {
			phonemes.push('6', 'M');
			i += 2;
		}
		// German vowels
		else if (char == 'a') {
			phonemes.push('6');
			i++;
		} else if (char == 'e') {
			phonemes.push('e_o');
			i++;
		} else if (char == 'i') {
			phonemes.push('i');
			i++;
		} else if (char == 'o') {
			phonemes.push('o');
			i++;
		} else if (char == 'u') {
			phonemes.push('M');
			i++;
		} else if (char == 'ä') {
			phonemes.push('6');
			i++;
		} else if (char == 'ö') {
			phonemes.push('V');
			i++;
		} else if (char == 'ü') {
			phonemes.push('M');
			i++;
		} else if (char == 'y') {
			phonemes.push('M');
			i++;
		} else {
			i++;
		}
	}

	return phonemes.join(' ');
}

// ==================== SPANISH PHONEME CONVERSION ====================
// Best for: ach-Laut (ch after back vowels) using Spanish 'x' [x]
function germanToSpanishPhonemes(word) {
	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var nextNextChar = i + 2 < word.length ? word[i + 2] : "";

		// German consonant clusters
		if (char == 's' && nextChar == 'c' && nextNextChar == 'h') {
			phonemes.push('sh');
			i += 3;
		} else if (char == 'c' && nextChar == 'h') {
			// ach-Laut - Spanish x [x] is a good match
			phonemes.push('x');
			i += 2;
		} else if (char == 'n' && nextChar == 'g') {
			phonemes.push('N');
			i += 2;
		} else if (char == 'n' && nextChar == 'k') {
			phonemes.push('N', 'k');
			i += 2;
		} else if (char == 'p' && nextChar == 'f') {
			phonemes.push('p', 'f');
			i += 2;
		} else if (char == 'q' && nextChar == 'u') {
			phonemes.push('k', 'U');
			i += 2;
		} else if (char == 't' && nextChar == 's') {
			phonemes.push('t', 's');
			i += 2;
		} else if (char == 's' && nextChar == 's') {
			phonemes.push('s', 's');
			i += 2;
		} else if (char == 'r' && nextChar == 'r') {
			phonemes.push('rr');
			i += 2;
		} else if (char == 'ß') {
			phonemes.push('s');
			i++;
		} else if (char == 'z') {
			phonemes.push('t', 's');
			i++;
		} else if (char == 'w') {
			phonemes.push('b');
			i++;
		} else if (char == 'v') {
			phonemes.push('f');
			i++;
		} else if (char == 'j') {
			phonemes.push('y');
			i++;
		} else if (char == 'h') {
			phonemes.push('x');
			i++;
		} else if (char == 'l') {
			phonemes.push('l');
			i++;
		} else if (char == 'm') {
			phonemes.push('m');
			i++;
		} else if (char == 'n') {
			phonemes.push('n');
			i++;
		} else if (char == 'r') {
			phonemes.push('r');
			i++;
		} else if (char == 's') {
			phonemes.push('s');
			i++;
		} else if (char == 't') {
			phonemes.push('t');
			i++;
		} else if (char == 'p') {
			phonemes.push('p');
			i++;
		} else if (char == 'k') {
			phonemes.push('k');
			i++;
		} else if (char == 'b') {
			phonemes.push('b');
			i++;
		} else if (char == 'd') {
			phonemes.push('d');
			i++;
		} else if (char == 'g') {
			phonemes.push('g');
			i++;
		} else if (char == 'f') {
			phonemes.push('f');
			i++;
		}
		// German diphthongs
		else if (char == 'e' && nextChar == 'i') {
			phonemes.push('a', 'I');
			i += 2;
		} else if (char == 'i' && nextChar == 'e') {
			phonemes.push('i', 'i');
			i += 2;
		} else if (char == 'e' && nextChar == 'u') {
			phonemes.push('o', 'I');
			i += 2;
		} else if (char == 'ä' && nextChar == 'u') {
			phonemes.push('o', 'I');
			i += 2;
		} else if (char == 'a' && nextChar == 'u') {
			phonemes.push('a', 'U');
			i += 2;
		}
		// German vowels
		else if (char == 'a') {
			phonemes.push('a');
			i++;
		} else if (char == 'e') {
			phonemes.push('e');
			i++;
		} else if (char == 'i') {
			phonemes.push('i');
			i++;
		} else if (char == 'o') {
			phonemes.push('o');
			i++;
		} else if (char == 'u') {
			phonemes.push('u');
			i++;
		} else if (char == 'ä') {
			phonemes.push('e');
			i++;
		} else if (char == 'ö') {
			phonemes.push('e');
			i++;
		} else if (char == 'ü') {
			phonemes.push('i');
			i++;
		} else if (char == 'y') {
			phonemes.push('i');
			i++;
		} else {
			i++;
		}
	}

	return phonemes.join(' ');
}

// ==================== COMMON HELPER FUNCTIONS ====================

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
		if (visited.indexOf(group.getUUID()) >= 0)
			continue;
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
			if (visited.indexOf(group.getUUID()) >= 0)
				continue;
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
			if (prop === 'length') {
				return target.getNumNotes();
			}
			if (typeof prop == "number") {
				return target.getNote(prop);
			}
			return target[prop];
		}
	});
}
