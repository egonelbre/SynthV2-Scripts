/*

This script sets all selected notes to perfect pitch by adding
pitch control curves that lock the pitch to the note's base pitch.

*/

var SCRIPT_TITLE = "Pitch Perfect";

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
				name: "clear",
				type: "CheckBox",
				text: SV.T("Clear existing pitch controls"),
				default: false
			},
			{
				"name" : "scope", "type" : "ComboBox",
				"label" : SV.T("Scope"),
				"choices" : [
					SV.T("Selected Notes"),
					SV.T("Current Track"),
					SV.T("Entire Project")
				],
				"default" : hasSelectedNotes() ? 0 : 2
			}
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
	if(notes.length == 0) {
		return;
	}

	if(options.clear) {
		var first = notes[0];
		var last = notes[notes.length-1];

		var start = first.getOnset();
		var end = last.getEnd();

		// clear the overlapping pitch controls
		var N = group.getNumPitchControls();
		for(var i = N - 1; i >= 0; i--) {
			var control = group.getPitchControl(i);

			var at = control.getPosition();
			var min = at;
			var max = at;

			var points = control.getPoints();
			for(var k = 0; k < points.length; k ++) {
				min = Math.min(min, at + points[k][0]);
				max = Math.max(max, at + points[k][0]);
			}

			if(end < min) {
				continue;
			}
			if(max < start) {
				continue;
			}

			group.removePitchControl(i);
		}
	}

	// Add perfect pitch controls for each note
	for(var i = 0; i < notes.length; i++) {
		var note = notes[i];
		var onset = note.getOnset();
		var duration = note.getDuration();
		var pitch = note.getPitch();

		// allow for some transitions between notes
		if(i + 1 < notes.length) {
			var nextNote = notes[i+1];
			if(note.getEnd() == nextNote.getOnset()) {
				duration -= SV.QUARTER / 4;
				duration = Math.max(0, duration);
			}
		}

		// Create a pitch control curve that spans the entire note
		// with zero deviation (perfect pitch)
		var control = SV.create("PitchControlCurve");
		control.setPosition(onset);
		control.setPitch(pitch);
		control.setPoints([
			[0, 0],
			[duration, 0]
		]);

		group.addPitchControl(control);
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
	var track = SV.getMainEditor().getCurrentTrack();
	var groupCount = track.getNumGroups();

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
