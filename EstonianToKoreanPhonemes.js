/// <reference path="./docs/synthesizer-v-api.d.ts" />

/*

This script converts Estonian lyrics to Korean X-SAMPA phonemes for Synthesizer V.

Korean X-SAMPA phonemes available:
Vowels: 6 (a), e_o (e), i, M (eu), o, V (eo)
Semivowels: w, M_ (ui), j (y)
Consonants: 4 (r/l flap), b, d, g, h, ts\h (jj), k, k_t (kk), l, m, n, p, pp, s, s_t (ss), t, tt, ts\_h (ch)
Codas: N (ng)

*/

const SCRIPT_TITLE = "Estonian to Korean Phonemes";

function getClientInfo() {
	return {
		"name": SV.T(SCRIPT_TITLE),
		"author": "Egon Elbre",
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

	for (var i = 0; i < notes.length; i++) {
		var note = notes[i];
		var lyrics = note.getLyrics();
		
		// Skip special markers and silence
		if (lyrics == "-" || lyrics == "+" || lyrics == "sil" || lyrics == "br" || lyrics == "SP" || lyrics == "AP") {
			continue;
		}

		// Set language override to Korean for proper phoneme interpretation
		note.setLanguageOverride("korean");

		var phonemes = estonianToKoreanPhonemes(lyrics.toLowerCase());
		if (phonemes) {
			note.setPhonemes(phonemes);
		}
	}
}

// Estonian to Korean phoneme conversion
function estonianToKoreanPhonemes(word) {
	if (!word) return "";
	
	var phonemes = [];
	var i = 0;
	
	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var prevChar = i > 0 ? word[i - 1] : "";
		
		// Consonants - mapped to Korean consonants
		if (char == 'h') {
			phonemes.push('h');
			i++;
		} else if (char == 'j') {
			// Estonian j is like English y, Korean has j (semivowel)
			phonemes.push('j');
			i++;
		} else if (char == 'l') {
			// ll -> geminate l
			if (nextChar == 'l') {
				phonemes.push('l', 'l');
				i += 2;
			} else {
				phonemes.push('4');  // Korean flap r/l
				i++;
			}
		} else if (char == 'm') {
			if (nextChar == 'm') {
				phonemes.push('m', 'm');
				i += 2;
			} else {
				phonemes.push('m');
				i++;
			}
		} else if (char == 'n') {
			// ng cluster -> N (coda)
			if (nextChar == 'g') {
				phonemes.push('N');
				i += 2;
			} else if (nextChar == 'n') {
				phonemes.push('n', 'n');
				i += 2;
			} else {
				phonemes.push('n');
				i++;
			}
		} else if (char == 'r') {
			if (nextChar == 'r') {
				phonemes.push('4', '4');  // Korean r/l flap
				i += 2;
			} else {
				phonemes.push('4');
				i++;
			}
		} else if (char == 's') {
			// sh digraph -> just use s (Korean doesn't have sh)
			if (nextChar == 'h') {
				phonemes.push('s');
				i += 2;
			} else if (nextChar == 's') {
				phonemes.push('s_t');  // Korean tense s
				i += 2;
			} else {
				phonemes.push('s');
				i++;
			}
		} else if (char == 't') {
			// ts cluster -> Korean ts\_h (ch)
			if (nextChar == 's') {
				phonemes.push('ts\\_h');
				i += 2;
			} else if (nextChar == 't') {
				phonemes.push('tt');  // Korean tense t
				i += 2;
			} else {
				phonemes.push('t');
				i++;
			}
		} else if (char == 'p') {
			if (nextChar == 'p') {
				phonemes.push('pp');  // Korean tense p
				i += 2;
			} else {
				phonemes.push('p');
				i++;
			}
		} else if (char == 'k') {
			if (nextChar == 'k') {
				phonemes.push('k_t');  // Korean tense k
				i += 2;
			} else {
				phonemes.push('k');
				i++;
			}
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
			// Korean doesn't have f, use p
			phonemes.push('p');
			i++;
		} else if (char == 'v') {
			// Korean doesn't have v, use b
			phonemes.push('b');
			i++;
		} else if (char == 'z') {
			// Korean doesn't have z, use s
			phonemes.push('s');
			i++;
		} else if (char == 'w') {
			phonemes.push('w');
			i++;
		}
		
		// Vowels - mapped to Korean vowels
		else if (char == 'a') {
			// Check for long vowel
			if (nextChar == 'a') {
				phonemes.push('6', '6');  // Korean 'a'
				i += 2;
			} else {
				phonemes.push('6');
				i++;
			}
		} else if (char == 'e') {
			if (nextChar == 'e') {
				phonemes.push('e_o', 'e_o');  // Korean 'e'
				i += 2;
			} else {
				phonemes.push('e_o');
				i++;
			}
		} else if (char == 'i') {
			if (nextChar == 'i') {
				phonemes.push('i', 'i');
				i += 2;
			} else {
				phonemes.push('i');
				i++;
			}
		} else if (char == 'o') {
			if (nextChar == 'o') {
				phonemes.push('o', 'o');
				i += 2;
			} else {
				phonemes.push('o');
				i++;
			}
		} else if (char == 'u') {
			if (nextChar == 'u') {
				phonemes.push('M', 'M');  // Korean 'eu'
				i += 2;
			} else {
				phonemes.push('M');
				i++;
			}
		} else if (char == 'õ') {
			// Estonian õ is a close-mid back unrounded vowel
			// Korean 'V' (eo) is somewhat similar
			if (nextChar == 'õ') {
				phonemes.push('V', 'V');
				i += 2;
			} else {
				phonemes.push('V');
				i++;
			}
		} else if (char == 'ä') {
			// Estonian ä is like [æ]
			// Use Korean '6' (a)
			if (nextChar == 'ä') {
				phonemes.push('6', '6');
				i += 2;
			} else {
				phonemes.push('6');
				i++;
			}
		} else if (char == 'ö') {
			// Estonian ö [ø] - front rounded vowel
			// Approximate with Korean 'V' (eo)
			if (nextChar == 'ö') {
				phonemes.push('V', 'V');
				i += 2;
			} else {
				phonemes.push('V');
				i++;
			}
		} else if (char == 'ü') {
			// Estonian ü [y] is close front rounded vowel
			// Approximate with Korean 'M' (eu)
			if (nextChar == 'ü') {
				phonemes.push('M', 'M');
				i += 2;
			} else {
				phonemes.push('M');
				i++;
			}
		} else if (char == 'y') {
			// Sometimes used instead of ü in older texts
			if (nextChar == 'y') {
				phonemes.push('M', 'M');
				i += 2;
			} else {
				phonemes.push('M');
				i++;
			}
		} 
		
		// Skip unknown characters
		else {
			i++;
		}
	}
	
	return phonemes.join(' ');
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
