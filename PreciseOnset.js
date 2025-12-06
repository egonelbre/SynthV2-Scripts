/*

This script ensures note onsets and phrase offsets are precise.

*/


var SCRIPT_TITLE = "Precise Onset";

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
        name: "phrase",
        type: "CheckBox",
        text: SV.T("Only phrase starts."),
        default: false
      },
      {
        name: "clear",
        type: "CheckBox",
        text: SV.T("Clear all pitch controls."),
        default: true
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
      },
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
    if(selectedNotes.length <= 1) {
      SV.finish();
      return;
    }

    var scope = SV.getMainEditor().getCurrentGroup();
    var group = scope.getTarget();

    selectedNotes = sortNotes(selectedNotes);
    processNotes(selectedNotes, group, scope, options);
    SV.finish();
    return
  } else if (options.scope == 1) {
    // Track
    var track = SV.getMainEditor().getCurrentTrack();
    var trackGroupN = track.getNumGroups();
    for(var i = 0; i < trackGroupN; i ++) {
      var scope = track.getGroupReference(i);
      var group = scope.getTarget();
      processNotes(group, group, scope, options);
    }
    SV.finish();
    return
  } else if (options.scope == 2) {
    // Project

    // TODO: this may double visit groups
    var project = SV.getProject();
    for(var i = 0; i < project.getNumTracks(); i ++) {
      var track = project.getTrack(i);
      var trackGroupN = track.getNumGroups();
      for(var k = 0; k < trackGroupN; k ++) {
        var scope = track.getGroupReference(k);
        var group = scope.getTarget();
        processNotes(group, group, scope, options);
      }
    }
    SV.finish();
    return
  }
}

function processNotes(notes, group, scope, options) {
  var notesN = getLength(notes);
  if(notesN == 0) {
    return;
  }

  if(!options.clear) {
    var first = getNote(notes, 0);
    var last = getNote(notes, notesN-1);

    var start = first.getOnset();
    var end = last.getEnd();

    // clear the overlapping pitch controls.
    // TODO: this clears more than necessary, but good enough for me
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
  } else {
    var N = group.getNumPitchControls();
    for(var i = N - 1; i >= 0; i--) {
      group.removePitchControl(i);
    }
  }


  var eps = SV.QUARTER/16;

  if(options.phrase) {
    var lastEnd = 0;
    for(var i = 0; i < notesN; i++) {
      var note = getNote(notes, i);
      var onset = note.getOnset();

      // check whether note is part of a phrase
      if(lastEnd == onset) {
        lastEnd = note.getEnd();
        continue;
      }

      var prep = eps;
      if(onset - lastEnd < prep) {
        prep = onset - lastEnd;
      }
      var dur = eps;
      if(dur > note.getDuration()) {
        dur = note.getDuration();
      }

      // add onset control
      var control = SV.create("PitchControlCurve");
      control.setPosition(onset);
      control.setPitch(note.getPitch());
      control.setPoints([
        [-prep, 0],
        [dur, 0],
      ]);

      group.addPitchControl(control);

      // add offset control
      var endprep = eps;
      if(endprep > note.getDuration() - dur) {
        endprep = note.getDuration() - dur;
      }

      var enddur = eps;
      var skip = false;
      if(i + 1 < notesN) {
        var next = getNote(notes, i + 1);
        var nextOnset = next.getOnset();

        if (note.getEnd() + enddur > nextOnset) {
          skip = true;
        }
      }

      if(!skip && endprep + enddur > 0){
        var control = SV.create("PitchControlCurve");
        control.setPosition(note.getEnd());
        control.setPitch(note.getPitch());
        control.setPoints([
          [-endprep, 0],
          [enddur, 0],
        ]);

        group.addPitchControl(control);
      }

      lastEnd = note.getEnd();
    }
  } else {
    var lastEnd = 0;
    for(var i = 0; i < notesN; i++) {
      var note = getNote(notes, i);
      var onset = note.getOnset();

      var prep = eps;
      if(onset - lastEnd < prep) {
        prep = onset - lastEnd;
      }
      var dur = eps;
      if(dur > note.getDuration()) {
        dur = note.getDuration();
      }

      // otherwise add a control point
      var control = SV.create("PitchControlCurve");
      control.setPosition(onset);
      control.setPitch(note.getPitch());
      control.setPoints([
        [-prep, 0],
        [dur, 0],
      ]);

      group.addPitchControl(control);

      // add offset control
      var endprep = eps;
      if(endprep > note.getDuration() - dur) {
        endprep = note.getDuration() - dur;
      }

      var enddur = eps;
      var skip = false;
      if(i + 1 < notesN) {
        var next = getNote(notes, i + 1);
        var nextOnset = next.getOnset();

        if (note.getEnd() + enddur > nextOnset) {
          skip = true;
        }
      }

      if(!skip && endprep + enddur > 0){
        var control = SV.create("PitchControlCurve");
        control.setPosition(note.getEnd());
        control.setPitch(note.getPitch());
        control.setPoints([
          [-endprep, 0],
          [enddur, 0],
        ]);

        group.addPitchControl(control);
      }

      lastEnd = note.getEnd();
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
