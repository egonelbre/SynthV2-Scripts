/// <reference path="./docs/synthesizer-v-api.d.ts"

/*

This script splits overlapping notes from a single track into
separate tracks, one per voice. It detects voices by finding
notes that overlap in time -- overlapping notes must belong to
different voices.

For example, a track with two voices interleaved:

Track 1:
  Voice 1: [C4----] [D4----]
  Voice 2:    [E3----]  [F3----]

Becomes:

Track 1: [C4----] [D4----]
Track 2:    [E3----]  [F3----]

*/

var SCRIPT_TITLE = "Split Voices";

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
	var track = SV.getMainEditor().getCurrentTrack();

	// Collect notes from all groups on the track.
	// Each entry tracks which group it belongs to for later removal.
	var entries = []; // {note, group, groupRef}
	for(var g = 0; g < track.getNumGroups(); g++) {
		var groupRef = track.getGroupReference(g);
		var group = groupRef.getTarget();
		var timeOffset = groupRef.getTimeOffset();
		for(var i = 0; i < group.getNumNotes(); i++) {
			entries.push({
				note: group.getNote(i),
				group: group,
				groupRef: groupRef,
				// Effective onset accounts for group time offset.
				effectiveOnset: group.getNote(i).getOnset() + timeOffset,
				effectiveEnd: group.getNote(i).getEnd() + timeOffset
			});
		}
	}

	if(entries.length == 0) {
		SV.showMessageBox(SV.T(SCRIPT_TITLE), SV.T("No notes found in the current track."));
		SV.finish();
		return;
	}

	// Sort by effective onset, then by pitch (descending) for stable ordering.
	entries.sort(function(a, b) {
		if(a.effectiveOnset < b.effectiveOnset) return -1;
		if(a.effectiveOnset > b.effectiveOnset) return 1;
		// Same onset: higher pitch first (voice 1 is typically the top voice).
		var pa = a.note.getPitch(), pb = b.note.getPitch();
		if(pa > pb) return -1;
		if(pa < pb) return 1;
		return 0;
	});

	// Assign each note to a voice using a greedy algorithm.
	// Each voice tracks its latest end time.
	var voiceEnds = []; // voiceEnds[v] = effective end blick of last note assigned to voice v
	var voiceAssignment = []; // voiceAssignment[i] = voice index for entries[i]

	for(var i = 0; i < entries.length; i++) {
		var onset = entries[i].effectiveOnset;
		var assigned = -1;

		// Find the first voice where this note doesn't overlap.
		for(var v = 0; v < voiceEnds.length; v++) {
			if(onset >= voiceEnds[v]) {
				assigned = v;
				break;
			}
		}

		if(assigned < 0) {
			// No existing voice fits, create a new one.
			assigned = voiceEnds.length;
			voiceEnds.push(0);
		}

		voiceEnds[assigned] = entries[i].effectiveEnd;
		voiceAssignment.push(assigned);
	}

	var numVoices = voiceEnds.length;
	if(numVoices <= 1) {
		SV.showMessageBox(SV.T(SCRIPT_TITLE), SV.T("No overlapping notes found. Nothing to split."));
		SV.finish();
		return;
	}

	var result = SV.showCustomDialog({
		"title": SV.T(SCRIPT_TITLE),
		"buttons": "OkCancel",
		"widgets": [
			{
				"name": "info", "type": "TextArea",
				"label": SV.T("Summary"),
				"height": 60,
				"default": "Found " + numVoices + " voices in " + entries.length + " notes.\n" +
					"Voice 1 stays on the current track.\n" +
					(numVoices - 1) + " new track(s) will be created."
			}
		]
	});
	if(result.status != 1) {
		SV.finish();
		return;
	}

	var project = SV.getProject();
	var trackName = track.getName() || "Track";

	// Create new tracks for voices 2..N.
	var newTracks = [];
	for(var v = 1; v < numVoices; v++) {
		var newTrack = SV.create("Track");
		newTrack.setName(trackName + " (Voice " + (v + 1) + ")");
		project.addTrack(newTrack);
		newTracks.push(newTrack);
	}

	// Rename original track.
	track.setName(trackName + " (Voice 1)");

	// Clone notes into new tracks and collect notes to remove from originals.
	// Group removals by source group so we can remove in reverse index order.
	var removals = {}; // groupUUID -> [note, ...]
	for(var i = 0; i < entries.length; i++) {
		var v = voiceAssignment[i];
		if(v == 0) continue; // stays on original track

		var entry = entries[i];
		var cloned = entry.note.clone();
		var newMainGroup = newTracks[v - 1].getGroupReference(0).getTarget();
		newMainGroup.addNote(cloned);

		var uuid = entry.group.getUUID();
		if(!removals[uuid]) {
			removals[uuid] = { group: entry.group, notes: [] };
		}
		removals[uuid].notes.push(entry.note);
	}

	// Remove moved notes from their original groups (reverse index order).
	for(var uuid in removals) {
		var group = removals[uuid].group;
		var notes = removals[uuid].notes;
		notes.sort(function(a, b) {
			return b.getIndexInParent() - a.getIndexInParent();
		});
		for(var i = 0; i < notes.length; i++) {
			group.removeNote(notes[i].getIndexInParent());
		}
	}

	SV.finish();
}
