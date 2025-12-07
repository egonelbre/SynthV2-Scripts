/*

This script merges any notes that are on the same pitch
and tied togehter.

For example:

[- ]
[Lorem] [+ ]     [- ] [- ] [+ ]

Will be merged into:

[- ]
[Lorem] [+ ]     [-      ] [+ ]

*/


var SCRIPT_TITLE = "Merge Tied Notes";

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
		processProject(processNotes, result.answers);
	} else {
		SV.showMessageBox(SV.T(SCRIPT_TITLE), SV.T("Invalid scope."));
	}
	SV.finish();
}

function processNotes(notes, group, options) {
	if(notes.length <= 1) {
		return;
	}

	var toRemove = [];

	var fragIdx = -1;
	var fragPitch = 0;
	var fragEnd = 0;

	for(var i = 0; i < notes.length; i++) {
		var note = notes[i];
		var lyrics = note.getLyrics();
		var pitch = note.getPitch()
		var currOnset = note.getOnset();

		var melisma = lyrics == "-";
		var samePitch = fragPitch == pitch;
		var continues = fragEnd == note.getOnset();

		// don't reset the fragment start when all the conditions are held.
		if(melisma && samePitch && continues) {
			fragEnd = note.getEnd();
			continue;
		}

		if(fragIdx + 1 != i) {
			// extend the first note
			var first = notes[fragIdx];
			first.setDuration(fragEnd - first.getOnset());
			// remove the other notes
			for(var k = fragIdx+1; k < i; k++) {
				toRemove.push(notes[k]);
			}
		}

		// reset the fragment
		fragIdx = i;
		fragPitch = pitch;
		fragEnd = note.getEnd();
	}

	if(fragIdx + 1 != notes.length) {
		// extend the first note
		var first = notes[fragIdx];
		first.setDuration(fragEnd - first.getOnset());
		// remove the other notes
		for(var k = fragIdx+1; k < i; k++) {
			toRemove.push(notes[k]);
		}
	}

	for(var k = 0; k < toRemove.length; k++) {
		group.removeNote(toRemove[k].getIndexInParent());
	}
}

// * Common * //

function hasSelectedNotes() {
	return SV.getMainEditor().getSelection().getSelectedNotes().length > 0;
}

function processSelection(process, options) {
	var selection = SV.getMainEditor().getSelection();
	var selectedNotes = selection.getSelectedNotes();
	selectedNotes = sortNotes(selectedNotes);

	var group = SV.getMainEditor().getCurrentGroup().getTarget();
	process(selectedNotes, group, options);
}

function processTrack(process, options) {
	var track = SV.getMainEditor().getCurrentTrack();
	var groupCount = track.getNumGroups();
	var visited = [];
	for(var i = 0; i < groupCount; i ++) {
		var group = track.getGroupReference(i).getTarget();

		// some note groups may be shared between or within tracks
		if(visited.indexOf(group.getUUID()) >= 0)
			continue;
		visited.push(group.getUUID());

		process(groupAsNotesArray(group), group, options);
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
