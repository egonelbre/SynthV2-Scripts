/// <reference path="./docs/synthesizer-v-api.d.ts" />

/*

A side-panel utility for moving or copying selected notes to the
note group of the track directly above or below the current one.

Useful when laying out harmony lines across separate voice tracks.

*/

const SCRIPT_TITLE = "Move/Copy to Voice";

function getClientInfo() {
	return {
		"name": SV.T(SCRIPT_TITLE),
		"category": "Utilities",
		"author": "Egon Elbre",
		"versionNumber": 1,
		"minEditorVersion": 131330,
		"type": "SidePanelSection"
	};
}

function getTranslations(langCode) {
	return [];
}

var copyUpValue = SV.create("WidgetValue");
var copyDownValue = SV.create("WidgetValue");
var moveUpValue = SV.create("WidgetValue");
var moveDownValue = SV.create("WidgetValue");

function findTrackByDisplayOrder(order) {
	var project = SV.getProject();
	for(var i = 0; i < project.getNumTracks(); i ++) {
		var track = project.getTrack(i);
		if(track.getDisplayOrder() === order) {
			return track;
		}
	}
	return null;
}

function getNeighborTrack(direction) {
	var current = SV.getMainEditor().getCurrentTrack();
	if(!current) return null;
	return findTrackByDisplayOrder(current.getDisplayOrder() + direction);
}

// Find a group reference in `targetTrack` that mirrors `sourceGroupRef`:
//   - If source is the track's main group, use the target's main group.
//   - Otherwise prefer a non-main library group that overlaps with the source's
//     absolute time range and does not reference the same NoteGroup as source.
function findTargetGroupRef(targetTrack, sourceGroupRef) {
	if(sourceGroupRef.isMain()) {
		for(var i = 0; i < targetTrack.getNumGroups(); i ++) {
			var ref = targetTrack.getGroupReference(i);
			if(ref.isMain()) return ref;
		}
		return null;
	}

	var sourceUUID = sourceGroupRef.getTarget().getUUID();
	var sourceStart = sourceGroupRef.getOnset();
	var sourceEnd = sourceGroupRef.getEnd();

	for(var i = 0; i < targetTrack.getNumGroups(); i ++) {
		var ref = targetTrack.getGroupReference(i);
		if(ref.isMain()) continue;
		if(ref.isInstrumental()) continue;
		if(ref.getTarget().getUUID() === sourceUUID) continue;
		if(ref.getOnset() < sourceEnd && ref.getEnd() > sourceStart) {
			return ref;
		}
	}
	return null;
}

function transferNotes(direction, isMove) {
	var selection = SV.getMainEditor().getSelection();
	var notes = selection.getSelectedNotes();
	if(notes.length === 0) return;

	var targetTrack = getNeighborTrack(direction);
	if(!targetTrack) return;

	var sourceGroupRef = SV.getMainEditor().getCurrentGroup();
	var sourceGroup = sourceGroupRef.getTarget();
	var sourceOffset = sourceGroupRef.getTimeOffset();

	var targetGroupRef = findTargetGroupRef(targetTrack, sourceGroupRef);
	if(!targetGroupRef) return;
	var targetGroup = targetGroupRef.getTarget();
	var targetOffset = targetGroupRef.getTimeOffset();

	// Capture source indices before any mutation so removal targets are stable.
	var sourceIndices = [];
	if(isMove) {
		for(var i = 0; i < notes.length; i ++) {
			sourceIndices.push(notes[i].getIndexInParent());
		}
		sourceIndices.sort(function(a, b) { return b - a; });
	}

	SV.getProject().newUndoRecord();

	for(var i = 0; i < notes.length; i ++) {
		var note = notes[i];
		var newNote = note.clone();
		var absoluteOnset = note.getOnset() + sourceOffset;
		newNote.setTimeRange(absoluteOnset - targetOffset, note.getDuration());
		targetGroup.addNote(newNote);
	}

	if(isMove) {
		for(var i = 0; i < sourceIndices.length; i ++) {
			sourceGroup.removeNote(sourceIndices[i]);
		}
	}
}

function canTransferTo(direction) {
	var track = getNeighborTrack(direction);
	if(!track) return false;
	var sourceGroupRef = SV.getMainEditor().getCurrentGroup();
	if(!sourceGroupRef) return false;
	return findTargetGroupRef(track, sourceGroupRef) !== null;
}

function updateButtonStates() {
	var hasSelection = SV.getMainEditor().getSelection().hasSelectedNotes();
	var canUp = hasSelection && canTransferTo(-1);
	var canDown = hasSelection && canTransferTo(1);
	copyUpValue.setEnabled(canUp);
	moveUpValue.setEnabled(canUp);
	copyDownValue.setEnabled(canDown);
	moveDownValue.setEnabled(canDown);
}

copyUpValue.setValueChangeCallback(function() { transferNotes(-1, false); });
copyDownValue.setValueChangeCallback(function() { transferNotes(1, false); });
moveUpValue.setValueChangeCallback(function() { transferNotes(-1, true); });
moveDownValue.setValueChangeCallback(function() { transferNotes(1, true); });

SV.getMainEditor().getSelection().registerSelectionCallback(function(selectionType) {
	if(selectionType == "note") {
		updateButtonStates();
	}
});

SV.getMainEditor().getSelection().registerClearCallback(function(selectionType) {
	if(selectionType == "notes") {
		updateButtonStates();
	}
});

updateButtonStates();

function getSidePanelSectionState() {
	return {
		"title": SV.T(SCRIPT_TITLE),
		"rows": [
			{
				"type": "Label",
				"text": SV.T("Copy:")
			},
			{
				"type": "Container",
				"columns": [
					{
						"type": "Button",
						"text": SV.T("Up"),
						"value": copyUpValue,
						"width": 0.5
					},
					{
						"type": "Button",
						"text": SV.T("Down"),
						"value": copyDownValue,
						"width": 0.5
					}
				]
			},
			{
				"type": "Label",
				"text": SV.T("Move:")
			},
			{
				"type": "Container",
				"columns": [
					{
						"type": "Button",
						"text": SV.T("Up"),
						"value": moveUpValue,
						"width": 0.5
					},
					{
						"type": "Button",
						"text": SV.T("Down"),
						"value": moveDownValue,
						"width": 0.5
					}
				]
			}
		]
	};
}
