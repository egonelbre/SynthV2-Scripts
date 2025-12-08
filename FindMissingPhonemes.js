/// <reference path="./docs/synthesizer-v-api.d.ts"

/*

This fixes issues with words spanning notes with rests.

*/

var SCRIPT_TITLE = "Find Missing Phonemes";

function getClientInfo() {
	return {
		"name" : SV.T(SCRIPT_TITLE),
		"author" : "Egon Elbre",
		"versionNumber" : 1,
		"minEditorVersion" : 65537
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
				"name" : "scope", "type" : "ComboBox",
				"label" : SV.T("Scope"),
				"choices" : [
					SV.T("Selected Notes"),
					SV.T("Current Track"),
					SV.T("Entire Project")
				],
				"default" : hasSelectedNotes() > 0 ? 0 : 2
			},
		],
	};

	var result = SV.showCustomDialog(form);
	if(result.status != 1) {
		SV.finish();
		return;
	}

	var scope = result.answers.scope;
	if(scope == 0) {
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

// @param notes: Note[]
// @param group: NoteGroup
// @param options: Object
// @param groupRef: NoteGroupReference
function processNotes(notes, group, options, groupRef) {
	if(notes.length <= 1) {
		return;
	}

	var measures = [];
	var timeAxis = SV.getProject().getTimeAxis();

	var allPhonemes = SV.getPhonemesForGroup(groupRef);
	for (var i = 0; i < notes.length; i++) {
		var phoneme = allPhonemes[i]
		if (phoneme) continue;

		var lyrics = notes[i].getLyrics();
		if(lyrics == "-") {
			continue;
		}

		var measure = timeAxis.getMeasureAt(notes[i].getOnset()+1);
		if(measures.indexOf(measure) <= 0) {
			measures.push(measure);
		}
	}

	if(measures.length > 0) {
		SV.showMessageBox(SV.T(SCRIPT_TITLE), SV.T("Missing phonemes in measure(s):") + measures.join(", "));
	}
}

// @param notes: Note[]
// @param group: Group
// @param options: Object
function processNotesX(notes, group, options, groupRef) {
	if(notes.length <= 1) {
		return;
	}

	var allPhonemes = SV.getPhonemesForGroup(groupRef);

	var first = 0;
	while (first < notes.length) {
		var note = notes[first];
		var next = first + 1;
		var timeEnd = note.getEnd();
		var disjoint = false;

		while(next < notes.length) {
			var syl = notes[next];
			var lyrics = syl.getLyrics();
			if(lyrics != "-" && lyrics != "+") {
				break;
			}
			disjoint = disjoint || (syl.getOnset() != timeEnd);
			timeEnd = syl.getEnd();
			next++;
		}

		if(!disjoint) {
			first = next;
			continue;
		}

		if(!distributePhonemes(notes, allPhonemes, first, next)) {
			return;
		}

		first = next;
	}
}

function distributePhonemes(notes, allPhonemes, from, to) {
	var first = notes[from];
	var phonemes = allPhonemes[from].split(" ");
	var groups = [];

	var i = 0;
	while (i < phonemes.length) {
		var group = [];
		groups.push(group);

		// consume clusives
		while (i < phonemes.length && isConsonant(phonemes[i])) {
			group.push(phonemes[i]);
			i++;
		}

		if (i >= phonemes.length) break;

		// add the vowel
		group.push(phonemes[i]);
		i++

		// if it's followed by two consonants then one of them belongs to the
		// current group.
		if(i+1 < phonemes.length &&
			isConsonant(phonemes[i]) &&
			isConsonant(phonemes[i+1])){
			group.push(phonemes[i]);
			i++;
		}
	}

	var result = "";
	for(var i = 0; i < groups.length; i++) {
		result += groups[i].join(" ") + "; ";
	}

	var timeAxis = SV.getProject().getTimeAxis();
	var measure = timeAxis.getMeasureAt(first.getOnset());

	return result == "OK";
}

function isConsonant(p) {
	return p == "p" || p == "t" || p == "k" ||
		p == "b" || p == "d" || p == "g" ||
		p == "m" || p == "n" ||
		p == "w" || p == "h" ||
		p == "r" || p == "l" || p == "s" || p=="z" || p=="sh" || p=="ch" ||
		p == "dh" || p == "th" || p == "f" || p == "v" || p == "j" || p == "y" ||
		p == "x" || p == "q" || p == "ng";
}

// * Common * //

function hasSelectedNotes() {
	return SV.getMainEditor().getSelection().hasSelectedNotes();
}

function processSelection(process, options) {
	var selection = SV.getMainEditor().getSelection();
	var selectedNotes = selection.getSelectedNotes();
	selectedNotes = sortNotes(selectedNotes);

	var groupRef = SV.getMainEditor().getCurrentGroup()
	var group = groupRef.getTarget();
	process(selectedNotes, group, options, groupRef);
}

function processTrack(process, options) {
	var track = SV.getMainEditor().getCurrentTrack();
	var groupCount = track.getNumGroups();
	var visited = [];
	for(var i = 0; i < groupCount; i ++) {
		var groupRef = track.getGroupReference(i);
		var group = groupRef.getTarget();

		// some note groups may be shared between or within tracks
		if(visited.indexOf(group.getUUID()) >= 0)
			continue;
		visited.push(group.getUUID());

		process(groupAsNotesArray(group), group, options, groupRef);
	}
}

function processProject(process, options) {
	// process all groups that may be shared between tracks
	var project = SV.getProject();
	for(var i = 0; i < project.getNumNoteGroupsInLibrary(); i ++) {
		var group = project.getNoteGroup(i);
		process(groupAsNotesArray(group), group, options);
	}

	// process unique groups for each track
	for(var i = 0; i < project.getNumTracks(); i ++) {
		var track = project.getTrack(i);
		var mainGroup = track.getGroupReference(0).getTarget();
		process(groupAsNotesArray(mainGroup), mainGroup, options);
	}
}

function processProjectWithRefs(process, options) {
	var visited = [];
	var project = SV.getProject();

	// process unique groups for each track
	for(var i = 0; i < project.getNumTracks(); i ++) {
		var track = project.getTrack(i);
		var groupCount = track.getNumGroups();
		var visited = [];
		for(var k = 0; k < groupCount; k ++) {
			var groupRef = track.getGroupReference(k);
			var group = groupRef.getTarget();

			// some note groups may be shared between or within tracks
			if(visited.indexOf(group.getUUID()) >= 0)
				continue;
			visited.push(group.getUUID());

			process(groupAsNotesArray(group), group, options, groupRef);
		}
	}
}

function sortNotes(notes) {
	return notes.sort(function(a,b) {
		if(a.getOnset() < b.getOnset()) return -1;
		if(a.getOnset() > b.getOnset()) return 1;
		return 0;
	});
}

function groupAsNotesArray(noteGroup) {
	return new Proxy(noteGroup, {
		get: function(target, prop) {
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
