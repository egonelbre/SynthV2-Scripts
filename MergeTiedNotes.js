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
  var selection = SV.getMainEditor().getSelection();
  var selectedNotes = selection.getSelectedNotes();

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
        "default" : selectedNotes.length > 0 ? 0 : 2
      },
    ],
  };

  var result = SV.showCustomDialog(form);
  if(result.status != 1) {
    SV.finish();
    return;
  }

  if(result.answers.scope == 0) {
    // Selection
    if(selectedNotes.length <= 1) {
      SV.finish();
      return;
    }

    var scope = SV.getMainEditor().getCurrentGroup();
    var group = scope.getTarget();

    selectedNotes = sortNotes(selectedNotes);
    mergeTiedNotes(selectedNotes, group);
    SV.finish();
    return
  } else if (result.answers.scope == 1) {
    // Track
    var track = SV.getMainEditor().getCurrentTrack();
    var trackGroupN = track.getNumGroups();
    var visited = [];
    for(var i = 0; i < trackGroupN; i ++) {
      var group = track.getGroupReference(i).getTarget();
      if(visited.indexOf(group.getUUID()) < 0) {
        mergeTiedNotes(group, group);
        visited.push(group.getUUID());
      }
    }
    SV.finish();
    return
  } else if (result.answers.scope == 2) {
    // Project
    var project = SV.getProject();
    for(var i = 0; i < project.getNumNoteGroupsInLibrary(); i ++) {
      var group = project.getNoteGroup(i);
      mergeTiedNotes(group, group);
    }
    for(var i = 0; i < project.getNumTracks(); i ++) {
      var track = project.getTrack(i);
      var mainGroup = track.getGroupReference(0).getTarget();
      mergeTiedNotes(mainGroup, mainGroup);
    }
    SV.finish();
    return
  }
}

function mergeTiedNotes(notes, group) {
  if(getLength(notes) <= 1) {
    return;
  }

  var toRemove = [];

  var fragIdx = -1;
  var fragPitch = 0;
  var fragEnd = 0;

  for(var i = 0; i < getLength(notes); i++) {
    var note = getNote(notes, i);
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
      var first = getNote(notes, fragIdx);
      first.setDuration(fragEnd - first.getOnset());
      // remove the other notes
      for(var k = fragIdx+1; k < i; k++) {
        toRemove.push(getNote(notes, k));
      }
    }

    // reset the fragment
    fragIdx = i;
    fragPitch = pitch;
    fragEnd = note.getEnd();
  }

  if(fragIdx + 1 != getLength(notes)) {
      // extend the first note
      var first = getNote(notes, fragIdx);
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

function sortNotes(notes) {
  return notes.sort(function(a,b) {
    if(a.getOnset() < b.getOnset()) return -1;
    if(a.getOnset() > b.getOnset()) return 1;
    return 0;
  });
}

function getNote(arr, index) {
  if(Array.isArray(arr)) {
    return arr[index];
  } else {
    // the input is a NoteGroup
    return arr.getNote(index);
  }
}

function getLength(arr) {
  if(Array.isArray(arr)) {
    return arr.length;
  } else {
    // the input is a NoteGroup
    return arr.getNumNotes();
  }
}
