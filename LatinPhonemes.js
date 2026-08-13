/// <reference path="./docs/synthesizer-v-api.d.ts" />

/*

This script converts Latin lyrics to Synthesizer V phonemes using
Ecclesiastical (Church) Latin pronunciation — the standard for sung Latin.

The Spanish phoneme set covers Ecclesiastical Latin almost exactly:
pure vowels a e i o u, rolled r/rr, ch [tʃ], sh [ʃ], J [ɲ], so all
words use Spanish (no adaptive multi-language selection needed).

Rules covered:
- ae, oe → e
- c before e/i → [tʃ], sc before e/i → [ʃ], otherwise hard [k]
- g before e/i → [dʒ] approximation, gn → [ɲ]
- ti + vowel (not after s/t/x) → [tsi]
- qu → [kw], ngu + vowel → [ŋgw]
- x → [ks], xc before e/i → [kʃ]
- h silent (mihi/nihil → [k] as word overrides)
- ch → [k], ph → [f], th → [t], y → [i]
- initial i + vowel is consonantal [j]

*/

const SCRIPT_TITLE = "Latin Phonemes";

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

		var phonemes = latinToSpanishPhonemes(lyrics.toLowerCase());
		if (phonemes) {
			note.setLanguageOverride("spanish");
			note.setPhonemes(phonemes);
		}
	}
}

// Words with exceptional pronunciation ([k] for h).
var LATIN_WORDS = {
	mihi: "m i k i",
	nihil: "n i k i l",
};

// Rewrites Ecclesiastical Latin orthography into a phonetic intermediate
// form using placeholder chars č [tʃ], š [ʃ], ñ [ɲ], ǧ [dʒ] and w for
// the [w] glide, resolving context-dependent spelling before table lookup.
function normalizeLatin(s) {
	s = s
		.replace(/æ/g, "e")
		.replace(/œ/g, "e")
		.replace(/ae/g, "e")
		.replace(/oe/g, "e")
		.replace(/y/g, "i")
		.replace(/ch/g, "k")
		.replace(/ph/g, "f")
		.replace(/th/g, "t")
		.replace(/rh/g, "r")
		.replace(/h/g, "") // h is silent
		.replace(/^i([aeou])/, "j$1") // consonantal i: iesu → jesu
		.replace(/gn/g, "ñ")
		.replace(/qu/g, "kw")
		.replace(/ngu([aeiou])/g, "ngw$1") // sanguis → [saŋgwis]
		.replace(/([^stx])ti([aeou])/g, "$1tsi$2") // gratia → [gratsia]
		.replace(/xc([ei])/g, "kš$1") // excelsis → [ekʃelsis]
		.replace(/x/g, "ks")
		.replace(/sc([ei])/g, "š$1")
		.replace(/c([ei])/g, "č$1")
		.replace(/g([ei])/g, "ǧ$1");
	return s;
}

var LATIN_DIGRAPHS = {
	rr: ["rr"],
	ng: ["N", "g"],
};

var LATIN_SINGLES = {
	b: "b", c: "k", d: "d", f: "f", g: "g",
	j: "I", k: "k", l: "l", m: "m", n: "n",
	p: "p", r: "r", s: "s", t: "t", v: "B", w: "U",
	"č": "ch", "š": "sh", "ñ": "J", "ǧ": "y",
	a: "a", e: "e", i: "i", o: "o", u: "u",
};

var LATIN_VOWELS = "aeiou";

function latinToSpanishPhonemes(word) {
	if (LATIN_WORDS[word]) {
		return LATIN_WORDS[word];
	}

	var w = normalizeLatin(word);
	var phonemes = [];

	for (var i = 0; i < w.length; i++) {
		var ch = w[i];

		if (i + 1 < w.length) {
			var pair = w.substr(i, 2);
			if (LATIN_DIGRAPHS[pair]) {
				phonemes.push.apply(phonemes, LATIN_DIGRAPHS[pair]);
				i++;
				continue;
			}
		}

		if (ch == "z") {
			phonemes.push("d", "s"); // [dz]
			continue;
		}

		var ph = LATIN_SINGLES[ch];
		if (!ph) continue;

		// Geminate (doubled character).
		if (i + 1 < w.length && w[i + 1] == ch) {
			// Doubled vowels emit a single phoneme to avoid rearticulation.
			if (LATIN_VOWELS.indexOf(ch) >= 0) {
				phonemes.push(ph);
			} else {
				phonemes.push(ph, ph);
			}
			i++;
			continue;
		}

		phonemes.push(ph);
	}

	return phonemes.join(" ");
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
		if (visited.indexOf(group.getUUID()) >= 0) continue;
		visited.push(group.getUUID());

		process(groupAsNotesArray(group), group, options, groupRef);
	}
}

function processProjectWithRefs(process, options) {
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
