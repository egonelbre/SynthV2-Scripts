/// <reference path="./docs/synthesizer-v-api.d.ts" />

/*

This script converts Estonian lyrics to Japanese ROMAJI phonemes for Synthesizer V.

Japanese ROMAJI phonemes available:
Vowels: a, i, u, e, o, N (ん)
Semivowels: w, v, y
Consonants: t, d, s, sh, j, z, ts, k, g, h, b, p, f, ch, n, m, r
Palatalized: ry, ky, py, dy, ty, ny, hy, my, gy, by

*/

const SCRIPT_TITLE = "Estonian to Japanese Phonemes";

function getClientInfo() {
	return {
		"name": SV.T(SCRIPT_TITLE),
		"author": "Egon Elbre",
		"category": "Estonian",
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

		// Set language override to Japanese for proper phoneme interpretation
		note.setLanguageOverride("japanese");

		var phonemes = estonianToJapanesePhonemes(lyrics.toLowerCase());
		if (phonemes) {
			note.setPhonemes(phonemes);
		}
	}
}

// Estonian to Japanese phoneme conversion
function estonianToJapanesePhonemes(word) {
	if (!word) return "";

	var phonemes = [];
	var i = 0;

	while (i < word.length) {
		var char = word[i];
		var nextChar = i + 1 < word.length ? word[i + 1] : "";
		var prevChar = i > 0 ? word[i - 1] : "";

		// Consonants - mapped to Japanese consonants
		if (char == 'h') {
			phonemes.push('h');
			i++;
		} else if (char == 'j') {
			// Estonian j is like y
			phonemes.push('y');
			i++;
		} else if (char == 'l') {
			// ll -> geminate (use r for Japanese approximation)
			if (nextChar == 'l') {
				phonemes.push('r', 'r');
				i += 2;
			} else {
				phonemes.push('r');
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
			// ng cluster -> N (moraic n)
			if (nextChar == 'g') {
				phonemes.push('N', 'g');
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
				phonemes.push('r', 'r');
				i += 2;
			} else {
				phonemes.push('r');
				i++;
			}
		} else if (char == 's') {
			// sh digraph
			if (nextChar == 'h') {
				phonemes.push('sh');
				i += 2;
			} else if (nextChar == 's') {
				phonemes.push('s', 's');
				i += 2;
			} else {
				phonemes.push('s');
				i++;
			}
		} else if (char == 't') {
			// ts cluster
			if (nextChar == 's') {
				phonemes.push('ts');
				i += 2;
			} else if (nextChar == 't') {
				phonemes.push('t', 't');
				i += 2;
			} else {
				phonemes.push('t');
				i++;
			}
		} else if (char == 'p') {
			if (nextChar == 'p') {
				phonemes.push('p', 'p');
				i += 2;
			} else {
				phonemes.push('p');
				i++;
			}
		} else if (char == 'k') {
			if (nextChar == 'k') {
				phonemes.push('k', 'k');
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
			phonemes.push('f');
			i++;
		} else if (char == 'v') {
			// Use Japanese v (ヴ)
			phonemes.push('v');
			i++;
		} else if (char == 'z') {
			phonemes.push('z');
			i++;
		} else if (char == 'w') {
			phonemes.push('w');
			i++;
		}

		// Vowels - mapped to Japanese vowels
		else if (char == 'a') {
			// Check for long vowel
			if (nextChar == 'a') {
				phonemes.push('a', 'a');
				i += 2;
			} else {
				phonemes.push('a');
				i++;
			}
		} else if (char == 'e') {
			if (nextChar == 'e') {
				phonemes.push('e', 'e');
				i += 2;
			} else {
				phonemes.push('e');
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
				phonemes.push('u', 'u');
				i += 2;
			} else {
				phonemes.push('u');
				i++;
			}
		} else if (char == 'õ') {
			// Estonian õ - approximate with o
			if (nextChar == 'õ') {
				phonemes.push('o', 'o');
				i += 2;
			} else {
				phonemes.push('o');
				i++;
			}
		} else if (char == 'ä') {
			// Estonian ä - approximate with a
			if (nextChar == 'ä') {
				phonemes.push('a', 'a');
				i += 2;
			} else {
				phonemes.push('a');
				i++;
			}
		} else if (char == 'ö') {
			// Estonian ö - approximate with o
			if (nextChar == 'ö') {
				phonemes.push('o', 'o');
				i += 2;
			} else {
				phonemes.push('o');
				i++;
			}
		} else if (char == 'ü') {
			// Estonian ü - approximate with u
			if (nextChar == 'ü') {
				phonemes.push('u', 'u');
				i += 2;
			} else {
				phonemes.push('u');
				i++;
			}
		} else if (char == 'y') {
			// Sometimes used instead of ü in older texts
			if (nextChar == 'y') {
				phonemes.push('u', 'u');
				i += 2;
			} else {
				phonemes.push('u');
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
