/*

This script ensures note onsets and phrase offsets are precise.

*/

const SCRIPT_TITLE = "Precise Onset and Offset";
const FLAG = "precise-pitch-control";

const ONSET_EPSILON = SV.QUARTER / 16;

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
				name: "phrase",
				type: "CheckBox",
				text: SV.T("Only first and last note of the phrase."),
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
	if(notes.length == 0) {
		return;
	}

	clearPitchControlsInRange(group, notes[0].getOnset(), notes[notes.length-1].getEnd(), FLAG);

	function touching (a, b) {
		return a.getEnd() + ONSET_EPSILON > b.getOnset();
	}

	for (var i = 0; i < notes.length; i++) {
		var note = notes[i];
		var prevNote = i - 1 >= 0 ? notes[i - 1] : null;
		var nextNote = i + 1 < notes.length ? notes[i + 1] : null;

		// create onset control
		var prepare = ONSET_EPSILON;
		var onset = ONSET_EPSILON;

		var phraseContinuation = prevNote && touching(prevNote, note);
		if(prevNote && phraseContinuation) {
			prepare = note.getOnset() - prevNote.getEnd();
		}
		if(onset > note.getDuration()) {
			onset = note.getDuration();
		}

		// only control onset if not part of a phrase (with option)
		var controlOnset = true;
		if (options.phrase) controlOnset = !phraseContinuation;

		if(controlOnset) {
			var control = SV.create("PitchControlCurve");
			control.setScriptData(FLAG, true);
			control.setPosition(note.getOnset());
			control.setPitch(note.getPitch());
			control.setPoints([
				[-prepare, 0],
				[onset, 0],
			]);

			group.addPitchControl(control);
		}

		// When it's the last note in a phrase add control to avoid glissando down.
		if(nextNote && !touching(note, nextNote)) {
			var end = ONSET_EPSILON;
			if(end > note.getDuration() - onset) {
				end = note.getDuration() - onset;
			}

			var control = SV.create("PitchControlCurve");
			control.setScriptData(FLAG, true);
			control.setPosition(note.getEnd());
			control.setPitch(note.getPitch());
			control.setPoints([
				[-end, 0],
				[ONSET_EPSILON, 0],
			]);

			group.addPitchControl(control);
		}
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

// * Pitch Control Helpers * //

function clearPitchControlsInRange(group, startPos, endPos, flag) {
	var pitchControlsToRemove = [];

	for(var i = 0; i < group.getNumPitchControls(); i ++) {
		var pitchControl = group.getPitchControl(i);
		var controlStart, controlEnd;

		if(pitchControl.getScriptData(flag) != true) {
			continue;
		}

		if(pitchControl.type === "PitchControlCurve") {
			// For curves, check the range of points
			var points = pitchControl.getPoints();
			if(points && points.length > 0) {
				var curvePos = pitchControl.getPosition();
				controlStart = curvePos + points[0][0];
				controlEnd = curvePos + points[points.length - 1][0];
			} else {
				// No points, treat as single position
				controlStart = controlEnd = pitchControl.getPosition();
			}
		} else {
			// For points, use single position
			controlStart = controlEnd = pitchControl.getPosition();
		}

		if(!(controlEnd < startPos || controlStart > endPos)) {
			pitchControlsToRemove.push(i);
		}
	}

	for(var i = pitchControlsToRemove.length - 1; i >= 0; i --) {
		group.removePitchControl(pitchControlsToRemove[i]);
	}
}
