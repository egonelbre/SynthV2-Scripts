/// <reference path="./docs/synthesizer-v-api.d.ts" />

/*

A side-panel utility with the usual quantize buttons for selected notes.

  - Quantize Start    : snap onsets to the grid; duration unchanged.
  - Quantize End      : snap ends to the grid; onset unchanged.
  - Quantize Both     : snap onsets and ends to the grid, preserving
                        continuity for notes that were touching.

*/

const SCRIPT_TITLE = "Quantize Notes";

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

// Grid resolutions in blicks. SV.QUARTER == one quarter note in blicks.
var GRID_CHOICES = [
	{ label: "1/1",   blicks: SV.QUARTER * 4 },
	{ label: "1/2",   blicks: SV.QUARTER * 2 },
	{ label: "1/4",   blicks: SV.QUARTER },
	{ label: "1/8",   blicks: SV.QUARTER / 2 },
	{ label: "1/16",  blicks: SV.QUARTER / 4 },
	{ label: "1/32",  blicks: SV.QUARTER / 8 },
	{ label: "1/4T",  blicks: SV.QUARTER * 2 / 3 },
	{ label: "1/8T",  blicks: SV.QUARTER / 3 },
	{ label: "1/16T", blicks: SV.QUARTER / 6 },
	{ label: "1/4.",  blicks: SV.QUARTER * 1.5 },
	{ label: "1/8.",  blicks: SV.QUARTER * 0.75 },
	{ label: "1/16.", blicks: SV.QUARTER * 0.375 }
];

var DEFAULT_GRID_INDEX = 4; // 1/16
var TOUCH_EPSILON = SV.QUARTER / 32;

// Quick-access grids exposed as their own rows of buttons.
var QUICK_GRIDS = [
	{ label: "1/8",  blicks: SV.QUARTER / 2 },
	{ label: "1/16", blicks: SV.QUARTER / 4 },
	{ label: "1/32", blicks: SV.QUARTER / 8 }
];

var gridValue = SV.create("WidgetValue");
gridValue.setValue(DEFAULT_GRID_INDEX);

var startButtonValue = SV.create("WidgetValue");
var endButtonValue = SV.create("WidgetValue");
var bothButtonValue = SV.create("WidgetValue");

for(var qi = 0; qi < QUICK_GRIDS.length; qi ++) {
	QUICK_GRIDS[qi].startValue = SV.create("WidgetValue");
	QUICK_GRIDS[qi].endValue = SV.create("WidgetValue");
	QUICK_GRIDS[qi].bothValue = SV.create("WidgetValue");
}

function snap(blicks, grid) {
	return Math.round(blicks / grid) * grid;
}

function getSelectedNotesSorted() {
	var notes = SV.getMainEditor().getSelection().getSelectedNotes().slice();
	notes.sort(function(a, b) {
		return a.getOnset() - b.getOnset();
	});
	return notes;
}

function getGridBlicks() {
	return GRID_CHOICES[gridValue.getValue()].blicks;
}

function quantizeStart(gridOverride) {
	var notes = getSelectedNotesSorted();
	if(notes.length === 0) return;
	var grid = gridOverride !== undefined ? gridOverride : getGridBlicks();
	SV.getProject().newUndoRecord();
	for(var i = 0; i < notes.length; i ++) {
		var note = notes[i];
		note.setTimeRange(snap(note.getOnset(), grid), note.getDuration());
	}
}

function quantizeEnd(gridOverride) {
	var notes = getSelectedNotesSorted();
	if(notes.length === 0) return;
	var grid = gridOverride !== undefined ? gridOverride : getGridBlicks();
	SV.getProject().newUndoRecord();
	for(var i = 0; i < notes.length; i ++) {
		var note = notes[i];
		var onset = note.getOnset();
		var newDuration = snap(note.getEnd(), grid) - onset;
		if(newDuration < grid / 8) newDuration = grid / 8;
		note.setTimeRange(onset, newDuration);
	}
}

function quantizeBoth(gridOverride) {
	var notes = getSelectedNotesSorted();
	if(notes.length === 0) return;
	var grid = gridOverride !== undefined ? gridOverride : getGridBlicks();
	SV.getProject().newUndoRecord();

	// Walk runs of touching notes; share the snapped boundary so they stay continuous.
	var i = 0;
	while(i < notes.length) {
		var runStart = i;
		while(i + 1 < notes.length &&
				Math.abs(notes[i + 1].getOnset() - notes[i].getEnd()) < TOUCH_EPSILON) {
			i ++;
		}
		var run = notes.slice(runStart, i + 1);
		i ++;

		var edges = [run[0].getOnset()];
		for(var k = 0; k < run.length; k ++) {
			edges.push(run[k].getEnd());
		}

		var snapped = [];
		for(var k = 0; k < edges.length; k ++) {
			snapped.push(snap(edges[k], grid));
		}

		// Guarantee strictly increasing edges so no note collapses to zero.
		for(var k = 1; k < snapped.length; k ++) {
			if(snapped[k] <= snapped[k - 1]) {
				snapped[k] = snapped[k - 1] + grid;
			}
		}

		for(var k = 0; k < run.length; k ++) {
			run[k].setTimeRange(snapped[k], snapped[k + 1] - snapped[k]);
		}
	}
}

function updateButtonStates() {
	var hasSelection = SV.getMainEditor().getSelection().hasSelectedNotes();
	startButtonValue.setEnabled(hasSelection);
	endButtonValue.setEnabled(hasSelection);
	bothButtonValue.setEnabled(hasSelection);
	for(var qi = 0; qi < QUICK_GRIDS.length; qi ++) {
		QUICK_GRIDS[qi].startValue.setEnabled(hasSelection);
		QUICK_GRIDS[qi].endValue.setEnabled(hasSelection);
		QUICK_GRIDS[qi].bothValue.setEnabled(hasSelection);
	}
}

startButtonValue.setValueChangeCallback(function() { quantizeStart(); });
endButtonValue.setValueChangeCallback(function() { quantizeEnd(); });
bothButtonValue.setValueChangeCallback(function() { quantizeBoth(); });

// Bind quick buttons; capture grid value with an immediately invoked function.
for(var qi = 0; qi < QUICK_GRIDS.length; qi ++) {
	(function(grid) {
		QUICK_GRIDS[qi].startValue.setValueChangeCallback(function() { quantizeStart(grid); });
		QUICK_GRIDS[qi].endValue.setValueChangeCallback(function() { quantizeEnd(grid); });
		QUICK_GRIDS[qi].bothValue.setValueChangeCallback(function() { quantizeBoth(grid); });
	})(QUICK_GRIDS[qi].blicks);
}

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

function quickRow(field) {
	var width = 1.0 / QUICK_GRIDS.length;
	var columns = [];
	for(var qi = 0; qi < QUICK_GRIDS.length; qi ++) {
		columns.push({
			"type": "Button",
			"text": QUICK_GRIDS[qi].label,
			"value": QUICK_GRIDS[qi][field],
			"width": width
		});
	}
	return { "type": "Container", "columns": columns };
}

function getSidePanelSectionState() {
	var choices = [];
	for(var i = 0; i < GRID_CHOICES.length; i ++) {
		choices.push(GRID_CHOICES[i].label);
	}

	return {
		"title": SV.T(SCRIPT_TITLE),
		"rows": [
			{
				"type": "Label",
				"text": SV.T("Quick Start:")
			},
			quickRow("startValue"),
			{
				"type": "Label",
				"text": SV.T("Quick End:")
			},
			quickRow("endValue"),
			{
				"type": "Label",
				"text": SV.T("Quick Start & End:")
			},
			quickRow("bothValue"),
			{
				"type": "Label",
				"text": SV.T("Grid:")
			},
			{
				"type": "Container",
				"columns": [
					{
						"type": "ComboBox",
						"choices": choices,
						"value": gridValue,
						"width": 1.0
					}
				]
			},
			{
				"type": "Container",
				"columns": [
					{
						"type": "Button",
						"text": SV.T("Start"),
						"value": startButtonValue,
						"width": 1.0 / 3
					},
					{
						"type": "Button",
						"text": SV.T("End"),
						"value": endButtonValue,
						"width": 1.0 / 3
					},
					{
						"type": "Button",
						"text": SV.T("Both"),
						"value": bothButtonValue,
						"width": 1.0 / 3
					}
				]
			}
		]
	};
}
