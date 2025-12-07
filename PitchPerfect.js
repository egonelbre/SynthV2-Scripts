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
  var selection = SV.getMainEditor().getSelection();
  var selectedNotes = selection.getSelectedNotes();

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
        "default" : selectedNotes.length > 0 ? 0 : 2
      }
    ],
  };

  var result = SV.showCustomDialog(form);
  if(result.status != 1) {
    SV.finish();
    return;
  }

  var options = result.answers;

  if(options.scope == 0) {
    // Selection
    if(selectedNotes.length == 0) {
      SV.finish();
      return;
    }

    var scope = SV.getMainEditor().getCurrentGroup();
    var group = scope.getTarget();

    selectedNotes = sortNotes(selectedNotes);
    processNotes(selectedNotes, group, options);
    SV.finish();
    return
  } else if (options.scope == 1) {
    // Track
    var track = SV.getMainEditor().getCurrentTrack();
    var trackGroupN = track.getNumGroups();
    var visited = [];
    for(var i = 0; i < trackGroupN; i ++) {
      var scope = track.getGroupReference(i);
      var group = scope.getTarget();
      if(visited.indexOf(group.getUUID()) < 0) {
        processNotes(groupAsNotesArray(group), group, options);
        visited.push(group.getUUID());
      }
    }
    SV.finish();
    return
  } else if (options.scope == 2) {
    // Project
    var project = SV.getProject();
    for(var i = 0; i < project.getNumNoteGroupsInLibrary(); i ++) {
      var group = project.getNoteGroup(i);
      processNotes(groupAsNotesArray(group), group, options);
    }
    for(var i = 0; i < project.getNumTracks(); i ++) {
      var track = project.getTrack(i);
      var mainGroup = track.getGroupReference(0).getTarget();
      processNotes(groupAsNotesArray(mainGroup), mainGroup, options);
    }
    SV.finish();
    return
  }
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
