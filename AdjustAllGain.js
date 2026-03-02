/// <reference path="./docs/synthesizer-v-api.d.ts"

/*

This script adjusts the gain of every track in the project by a specified amount.

*/

var SCRIPT_TITLE = "Adjust All Gain";

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
	var project = SV.getProject();
	var numTracks = project.getNumTracks();

	// Build a summary of current track gains.
	var trackInfo = "";
	for (var i = 0; i < numTracks; i++) {
		var track = project.getTrack(i);
		var mixer = track.getMixer();
		var name = track.getName() || ("Track " + (i + 1));
		trackInfo += name + ": " + mixer.getGainDecibel().toFixed(1) + " dB\n";
	}

	var form = {
		"title": SV.T(SCRIPT_TITLE),
		"buttons": "OkCancel",
		"widgets": [
			{
				"name": "info",
				"type": "TextArea",
				"label": SV.T("Current gains"),
				"height": 80,
				"default": trackInfo
			},
			{
				"name": "mode",
				"type": "ComboBox",
				"label": SV.T("Mode"),
				"choices": [
					SV.T("Adjust by (relative)"),
					SV.T("Set to (absolute)")
				],
				"default": 0
			},
			{
				"name": "gain",
				"type": "Slider",
				"label": SV.T("Gain (dB)"),
				"format": "%1.1f",
				"minValue": -24.0,
				"maxValue": 24.0,
				"interval": 0.1,
				"default": 0
			}
		]
	};

	var result = SV.showCustomDialog(form);
	if (result.status != 1) {
		SV.finish();
		return;
	}

	var mode = result.answers.mode;
	var gain = result.answers.gain;

	for (var i = 0; i < numTracks; i++) {
		var track = project.getTrack(i);
		var mixer = track.getMixer();

		if (mode == 0) {
			// Relative: add to current gain
			var newGain = mixer.getGainDecibel() + gain;
			newGain = Math.max(-24.0, Math.min(24.0, newGain));
			mixer.setGainDecibel(newGain);
		} else {
			// Absolute: set to exact value
			mixer.setGainDecibel(gain);
		}
	}

	SV.finish();
}
