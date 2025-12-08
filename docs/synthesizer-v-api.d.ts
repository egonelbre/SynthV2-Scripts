/**
 * Type definitions for Dreamtonics Synthesizer V Studio Scripting API
 * Generated from official documentation
 * https://resource.dreamtonics.com/scripting/index.html
 */

declare class ArrangementSelectionState {
  /**
   * Unselects all object types supported by this selection state. Return true if the selection has changed.
   * @returns boolean
   */
  clearAll(): boolean;

  /**
   * Unselect all `NoteGroupReference`. Return true if the selection has changed.
   * @returns boolean
   */
  clearGroups(): boolean;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get an array of selected `NoteGroupReference` following the order of selection.
   */
  getSelectedGroups(): void;

  /**
   * Check if there's anything selected.
   */
  hasSelectedContent(): void;

  /**
   * Check if there is at least one `NoteGroupReference` selected.
   */
  hasSelectedGroups(): void;

  /**
   * Check if there's any unfinished edit on the selected objects.
   */
  hasUnfinishedEdits(): void;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Attach a script function to be called when the selection is cleared. The callback function will receive one argument: the type of the cleared objects.
   * @param callback
   */
  registerClearCallback(callback: Function): void;

  /**
   * Attach a script function to be called when objects are selected or deselected by the user. The callback function will receive two arguments: the type of the selected objects and whether this is a selection or deselection operation.
   * @param callback
   */
  registerSelectionCallback(callback: Function): void;

  /**
   * Add a `NoteGroupReference` to the selection.
   */
  selectGroup(): void;

  /**
   * Unselect a `NoteGroupReference`. Return true if the selection has changed.
   * @returns boolean
   */
  unselectGroup(): boolean;

}

declare class ArrangementView {
  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the coordinate system for the track arrangement area.
   * @returns CoordinateSystem
   */
  getNavigation(): CoordinateSystem;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the selection state object for arrangement view.
   * @returns ArrangementSelectionState
   */
  getSelection(): ArrangementSelectionState;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

}

declare class Automation {
  /**
   * Add a control point with position `b` (blicks) and parameter value `v`. If there is already a point on `b`, the parameter value will get updated to `v`.
   * @param b
   * @param v
   * @returns boolean
   */
  add(b: number, v: number): boolean;

  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns Automation
   */
  clone(): Automation;

  /**
   * Get the interpolated parameter value at position `b` (blicks). If a point exists at `b`, the interpolation is guaranteed to return the value for the point, regardless of the interpolation method.
   * @param b
   * @returns number
   */
  get(b: number): number;

  /**
   * A version of `Automation#getPoints` with unlimited range.
   * @returns any[]
   */
  getAllPoints(): any[];

  /**
   * Get a definition object with the following properties,
   * @returns any
   */
  getDefinition(): any;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Returns how values between control points are interpolated:
   * @returns string
   */
  getInterpolationMethod(): string;

  /**
   * A version of `Automation#get` that uses linear interpolation (even if `Automation#getInterpolationMethod` is not "Linear").
   * @param b
   * @returns number
   */
  getLinear(b: number): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get an array of control points whose positions are between `begin` and `end` (blicks). Each element in the array is an array of two elements: a `number` for the position (blicks) and a `number` for the parameter value. For example, `[[0, 0.1], [5000, 0], [10000, -0.1]]`.
   * @param begin
   * @param end
   * @returns any[]
   */
  getPoints(begin: number, end: number): any[];

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Get the parameter type for this `Automation`. See the `typeName` column of the table in `Automation#getDefinition`.
   * @returns string
   */
  getType(): string;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove the control point at position `b` (blicks) if there is one.
   * @param b
   * @returns boolean
   */
  remove(b: number): boolean;

  /**
   * Remove all control points in the `Automation`.
   */
  removeAll(): void;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

  /**
   * Simplify the parameter curve from position `begin` (blicks) to position `end` (blicks) by removing control points that do not significantly contribute to the curve's shape. If `threshold` is not provided, it will be set to 0.002. Higher values of `threshold` will result in more simplification.
   * @param begin
   * @param end
   */
  simplify(begin: number, end: number): void;

}

declare class CoordinateSystem {
  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the scaling factor in the horizontal direction.
   * @returns number
   */
  getTimePxPerUnit(): number;

  /**
   * Get the current visible time range. It returns an array with two `number` elements corresponding to the starting time and ending time. The time unit is blicks.
   * @returns any[]
   */
  getTimeViewRange(): any[];

  /**
   * Get the scaling factor in the vertical direction.
   */
  getValuePxPerUnit(): void;

  /**
   * Get the current visible value range. It returns an array with two `number` elements corresponding to the lower value and upper value. For the piano roll, the unit is MIDI number (semitones); for arrangement view, its value bears no meaning.
   */
  getValueViewRange(): void;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Move the visible area so the left end is at `time`.
   * @param time
   */
  setTimeLeft(time: number): void;

  /**
   * Move the visible area so the right end is at `time`.
   * @param time
   */
  setTimeRight(time: number): void;

  /**
   * Set the horizontal scaling factor to `scale`.
   * @param scale
   */
  setTimeScale(scale: number): void;

  /**
   * Move the visible area so the vertical center is at `v`.
   * @param v
   */
  setValueCenter(v: number): void;

  /**
   * Round a time position `b` based on snapping settings.
   * @param b
   * @returns number
   */
  snap(b: number): number;

  /**
   * Convert a time position to an x-position (pixels).
   * @param t
   * @returns number
   */
  t2x(t: number): number;

  /**
   * Convert a value to a y-position (pixels).
   * @param v
   * @returns number
   */
  v2y(v: number): number;

  /**
   * Convert an x-position (pixels) to a time position.
   * @param x
   * @returns number
   */
  x2t(x: number): number;

  /**
   * Convert a y-position (pixels) to a value.
   * @param y
   * @returns number
   */
  y2v(y: number): number;

}

declare class GroupSelection {
  /**
   * Unselect all `NoteGroupReference`. Return true if the selection has changed.
   * @returns boolean
   */
  clearGroups(): boolean;

  /**
   * Get an array of selected `NoteGroupReference` following the order of selection.
   */
  getSelectedGroups(): void;

  /**
   * Check if there is at least one `NoteGroupReference` selected.
   */
  hasSelectedGroups(): void;

  /**
   * Add a `NoteGroupReference` to the selection.
   */
  selectGroup(): void;

  /**
   * Unselect a `NoteGroupReference`. Return true if the selection has changed.
   * @returns boolean
   */
  unselectGroup(): boolean;

}

declare class MainEditorView {
  /**
   * Get the current `NoteGroupReference` that the user is working inside. If the user has not entered a `NoteGroupReference`, return the main group of the current track.
   * @returns NoteGroupReference
   */
  getCurrentGroup(): NoteGroupReference;

  /**
   * Get the current `Track` opened in the piano roll.
   * @returns Track
   */
  getCurrentTrack(): Track;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the `CoordinateSystem` of the piano roll.
   * @returns CoordinateSystem
   */
  getNavigation(): CoordinateSystem;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the selection state object for the piano roll.
   * @returns TrackInnerSelectionState
   */
  getSelection(): TrackInnerSelectionState;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

}

declare class NestedObject {
  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

}

declare class Note {
  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns Note
   */
  clone(): Note;

  /**
   * Get an object holding note properties. The object has the following properties.
   * @returns any
   */
  getAttributes(): any;

  /**
   * Get the pitch adjustment in cents. 100 cents equals one semitone. This adjustment is applied on top of the base pitch of the note. (supported since 2.1.1)
   * @returns number
   */
  getDetune(): number;

  /**
   * Get the duration of the note. The unit is blicks.
   * @returns number
   */
  getDuration(): number;

  /**
   * Get the end position (start + duration) of the note. The unit is blicks.
   * @returns number
   */
  getEnd(): number;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the note-specific language. This returns empty when the note is using the group/track's default language. (supported since 1.9.0b2)
   */
  getLanguageOverride(): void;

  /**
   * Get the lyrics for this note.
   * @returns string
   */
  getLyrics(): string;

  /**
   * Get the type of the note ("sing" or "rap"). (supported since 1.9.0b2)
   * @returns string
   */
  getMusicalType(): string;

  /**
   * Get the start position of the note. The unit is blicks.
   * @returns number
   */
  getOnset(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Returns the user-specified phonemes, delimited by spaces. For example, "hh ah ll ow".
   * @returns string
   */
  getPhonemes(): string;

  /**
   * Get the pitch as a MIDI number. C4 maps to 60.
   * @returns number
   */
  getPitch(): number;

  /**
   * Get the pitch mode of the note: true for auto, false for manual. (supported since 1.9.0b2)
   * @returns boolean
   */
  getPitchAutoMode(): boolean;

  /**
   * Get the accent (if available) for a rap note. (supported since 1.9.0b2)
   * @returns string
   */
  getRapAccent(): string;

  /**
   * Get the retake list for this note.
   * @returns RetakeList
   */
  getRetakes(): RetakeList;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set note properties based on an attribute object. The attribute object does not have to be complete; only the given properties will be updated. For example,
   */
  setAttributes(): void;

  /**
   * Set the pitch adjustment in cents. 100 cents equals one semitone. This adjustment is applied on top of the base pitch of the note. (supported since 2.1.1)
   */
  setDetune(): void;

  /**
   * Resize the note to duration `t`. The unit is blicks. This changes the end as well, but not the onset.
   * @param t
   */
  setDuration(t: number): void;

  /**
   * Set the language for the note to override the track/group level language settings. Available options : "mandarin", "japanese", "english", "cantonese" (supported since 1.9.0b2)
   * @param language
   */
  setLanguageOverride(language: string): void;

  /**
   * Change the lyrics.
   * @param lyrics
   */
  setLyrics(lyrics: string): void;

  /**
   * Set the note type ("sing" or "rap"). (supported since 1.9.0b2)
   * @param type
   */
  setMusicalType(type: string): void;

  /**
   * Move the note to start at `t`. The unit is blicks. This does not change the duration.
   * @param t
   */
  setOnset(t: number): void;

  /**
   * Change the phonemes to `phoneme_str`. For example, "hh ah ll ow".
   */
  setPhonemes(): void;

  /**
   * Set the note pitch to `pitchNumber`, a MIDI number.
   * @param pitchNumber
   */
  setPitch(pitchNumber: number): void;

  /**
   * Set whether the note has auto pitch mode (true) or manual pitch mode (false). (supported since 1.9.0b2)
   * @param isAuto
   */
  setPitchAutoMode(isAuto: boolean): void;

  /**
   * Set the accent for rap notes. Note that rap accent is only used in Mandarin Chinese where there are five accent types (1, 2, 3, 4, 5). (supported since 1.9.0b2)
   * @param accent
   */
  setRapAccent(accent: string): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

  /**
   * Set both onset and duration. This is a shorthand for calling `setOnset(onset)` and `setDuration(duration)`.
   * @param onset
   * @param duration
   */
  setTimeRange(onset: number, duration: number): void;

}

declare class NoteGroup {
  /**
   * Add a note to this `NoteGroup` and return the index of the added note. The notes are kept sorted by ascending onset positions.
   * @returns number
   */
  addNote(): number;

  /**
   * Add a pitch control object to this `NoteGroup` and return the index of the added object. The pitch control objects are kept sorted by ascending anchor positions.
   * @returns number
   */
  addPitchControl(): number;

  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns NoteGroup
   */
  clone(): NoteGroup;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the user-specified name of this `NoteGroup`.
   * @returns string
   */
  getName(): string;

  /**
   * Get the note at `index`. The notes inside a `NoteGroup` are always sorted by onset positions.
   * @param index
   * @returns Note
   */
  getNote(index: number): Note;

  /**
   * Get the number of notes in the `NoteGroup`.
   * @returns number
   */
  getNumNotes(): number;

  /**
   * Get the number of pitch control objects in the `NoteGroup`.
   */
  getNumPitchControls(): void;

  /**
   * Get the `Automation` object for parameter `type`. It is case-insensitive.
   * @param type
   * @returns Automation
   */
  getParameter(type: string): Automation;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the pitch control object at `index`. The pitch control objects inside a `NoteGroup` are kept sorted by anchor positions.
   * @param index
   * @returns PitchControlPoint|PitchControlCurve
   */
  getPitchControl(index: number): PitchControlPoint|PitchControlCurve;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Get the Universally Unique Identifier. Unlike the name, a UUID is unique across the project and can be used to associate a `NoteGroupReference` with a `NoteGroup`.
   * @returns string
   */
  getUUID(): string;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove the note at `index`.
   * @param index
   */
  removeNote(index: number): void;

  /**
   * Remove the pitch control object at `index`.
   * @param index
   */
  removePitchControl(index: number): void;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set the name of this `NoteGroup`.
   * @param name
   */
  setName(name: string): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class NoteGroupReference {
  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns NoteGroupReference
   */
  clone(): NoteGroupReference;

  /**
   * The duration of this `NoteGroupReference` (blicks).
   * @returns number
   */
  getDuration(): number;

  /**
   * Get the ending position (blicks), that is, the end of the last note in the target `NoteGroup` plus the time offset.
   * @returns number
   */
  getEnd(): number;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the beginning position (blicks), that is, the onset of the first `Note` in the target `NoteGroup` plus the time offset.
   * @returns number
   */
  getOnset(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the pitch shift (semitones) applied to all notes in the target `NoteGroup`}.
   * @returns number
   */
  getPitchOffset(): number;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Get the target `NoteGroup`.
   * @returns NoteGroup
   */
  getTarget(): NoteGroup;

  /**
   * Get the time offset (blicks) applied to all notes in the target `NoteGroup`.
   * @returns number
   */
  getTimeOffset(): number;

  /**
   * Get an object holding the default voice properties for this group, similar to `Note#getAttributes`.
   * @returns any
   */
  getVoice(): any;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Whether this `NoteGroupReference` refers to an external audio file. If so, it must not refer to a `NoteGroup`.
   * @returns boolean
   */
  isInstrumental(): boolean;

  /**
   * Whether this `NoteGroupReference` refers to the parent `Track`'s main group.
   * @returns boolean
   */
  isMain(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Check if this group is muted.
   * @returns boolean
   */
  isMuted(): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set the mute status of this group.
   * @param muted
   */
  setMuted(muted: boolean): void;

  /**
   * Set the pitch offset to `pitchOffset` (semitones).
   * @param pitchOffset
   */
  setPitchOffset(pitchOffset: number): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

  /**
   * Set the target `NoteGroup`.
   */
  setTarget(): void;

  /**
   * Set the time offset to `blickOffset` (blicks).
   * @param blickOffset
   */
  setTimeOffset(blickOffset: number): void;

  /**
   * Set the absolute onset and duration of a group. This does not affect the time offset and can be used to shorten or extend a group from either left or right.
   * @param onset
   * @param duration
   */
  setTimeRange(onset: number, duration: number): void;

  /**
   * Set voice properties based on an attribute object (for the definition, see `NoteGroupReference#getVoice`). The attribute object does not have to be complete; only the given properties will be updated (see `Note#setAttributes`).
   * @param attributes
   */
  setVoice(attributes: any): void;

}

declare class PitchControlCurve {
  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns PitchControlCurve
   */
  clone(): PitchControlCurve;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the anchor pitch value of this pitch control curve in semitones relative to the pitch offset of the note group.
   * @returns number
   */
  getPitch(): number;

  /**
   * Get all control points of this pitch control curve.
   * @returns Array.&lt;Array.&lt;number&gt;&gt;
   */
  getPoints(): Array.&lt;Array.&lt;number&gt;&gt;;

  /**
   * Get the anchor position of this pitch control curve relative to the time offset of the note group (in blicks).
   * @returns number
   */
  getPosition(): number;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Get the interpolated pitch value at a specific time position.
   * @returns number
   */
  getValueAt(): number;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set the anchor pitch value of this pitch control curve.
   */
  setPitch(): void;

  /**
   * Set all control points of this pitch control curve.
   */
  setPoints(): void;

  /**
   * Set the anchor position of this pitch control curve.
   */
  setPosition(): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class PitchControlPoint {
  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns PitchControlPoint
   */
  clone(): PitchControlPoint;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the pitch value of this pitch control point in semitones relative to the pitch offset of the note group.
   * @returns number
   */
  getPitch(): number;

  /**
   * Get the position of this pitch control point relative to the time offset of the note group (in blicks).
   * @returns number
   */
  getPosition(): number;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set the pitch value of this pitch control point.
   */
  setPitch(): void;

  /**
   * Set the time position of this pitch control point.
   */
  setPosition(): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class PlaybackControl {
  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get the current playhead position in seconds.
   * @returns number
   */
  getPlayhead(): number;

  /**
   * Get the current playback status. It can be one of the following.
   * @returns string
   */
  getStatus(): string;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Start looping between `tBegin` and `tEnd` in seconds.
   * @param tBegin
   * @param tEnd
   */
  loop(tBegin: number, tEnd: number): void;

  /**
   * Stop playing but without resetting the playhead.
   */
  pause(): void;

  /**
   * Start playing audio.
   */
  play(): void;

  /**
   * Set the playhead position to `t` in seconds.
   * @param t
   */
  seek(t: number): void;

  /**
   * Stop playing and reset the playhead to where the playback started.
   */
  stop(): void;

}

declare class Project {
  /**
   * Insert a `NoteGroup` to the project library at `suggestedIndex`. If `suggestedIndex` is not given, the `NoteGroup` is added at the end. Return the index of the added `NoteGroup`.
   */
  addNoteGroup(): void;

  /**
   * Add a `Track` to the `Project`. Return the index of the added `Track`.
   * @returns number
   */
  addTrack(): number;

  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * Get the duration of the `Project` (blicks), defined as the duration of the longest `Track`.
   * @returns number
   */
  getDuration(): number;

  /**
   * Get the absolute path of the project on the file system.
   * @returns string
   */
  getFileName(): string;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * If `id` is a number, get the `id`-th `NoteGroup` in the project library.
   * @returns NoteGroup|undefined
   */
  getNoteGroup(): NoteGroup|undefined;

  /**
   * Get the number of `NoteGroup` in the project library.
   * @returns number
   */
  getNumNoteGroupsInLibrary(): number;

  /**
   * Get the number of tracks.
   * @returns number
   */
  getNumTracks(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Get the `TimeAxis` object of this `Project`.
   * @returns TimeAxis
   */
  getTimeAxis(): TimeAxis;

  /**
   * Get the `index`-th `Track`. The indexing is based on the storage order rather than display order.
   * @param index
   * @returns Track
   */
  getTrack(index: number): Track;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Add a new undo record for this `Project`. This means that all edits following the last undo record will be undone/redone together when the users press `Ctrl + Z` or `Ctrl + Y`.
   */
  newUndoRecord(): void;

  /**
   * Remove `index`-th `NoteGroup` from the project library. This also removes all `NoteGroupReference` that refer to the `NoteGroup`.
   * @param index
   */
  removeNoteGroup(index: number): void;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Remove the `index`-th `Track` from the `Project`.
   * @param index
   */
  removeTrack(index: number): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class RetakeList {
  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * Delete a retake by its ID.
   */
  deleteTake(): void;

  /**
   * Generate a new retake with the specified variation parameters.
   */
  generateTake(): void;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the number of retakes in this list.
   * @returns number
   */
  getNumTakes(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set the active retake by its ID.
   */
  setActiveTake(): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class SV {
  /**
   * Number of blicks in a quarter. The value is 705600000.
   */
  QUARTER(): void;

  /**
   * Get a localized version of `text` based on the current UI language settings.
   * @param text
   * @returns string
   */
  T(text: string): string;

  /**
   * Check whether the key (passed in as a MIDI number) is a black key on a piano.
   * @param k
   * @returns boolean
   */
  blackKey(k: number): boolean;

  /**
   * Convert `b` from number of blicks into number of quarters.
   * @param b
   * @returns number
   */
  blick2Quarter(b: number): number;

  /**
   * Convert `b` from blicks into seconds with the specified beats per minute `bpm`.
   * @param b
   * @param bpm
   * @returns number
   */
  blick2Seconds(b: number, bpm: number): number;

  /**
   * Rounded division of `dividend` (blicks) over `divisor` (blicks).
   * @param dividend
   * @param divisor
   */
  blickRoundDiv(dividend: number, divisor: number): void;

  /**
   * Returns the closest multiple of `interval` (blicks) from `b` (blick).
   * @param b
   * @param interval
   * @returns number
   */
  blickRoundTo(b: number, interval: number): number;

  /**
   * Create a new object. `type` can be one of the following type-specifying strings.
   * @returns any
   */
  create(): any;

  /**
   * Mark the finish of a script. All subsequent async callbacks will not be executed. Note that this does not cause the current script to exit immediately.
   */
  finish(): void;

  /**
   * Convert a frequency in Hz to a MIDI number (semitones, where C4 is 60).
   * @param f
   * @returns number
   */
  freq2Pitch(f: number): number;

  /**
   * Get the UI state object for arrangement view.
   * @returns ArrangementView
   */
  getArrangement(): ArrangementView;

  /**
   * Get computed attributes for all notes in a group (passed in as a group reference). (supported since 2.1.1)
   * @returns any[]
   */
  getComputedAttributesForGroup(): any[];

  /**
   * Get computed pitch values for a group (passed in as a group reference) over a specified time range. (supported since 2.1.1)
   */
  getComputedPitchForGroup(): void;

  /**
   * Get the text on the system clipboard.
   * @returns string
   */
  getHostClipboard(): string;

  /**
   * Get an object with the following properties.
   * @returns any
   */
  getHostInfo(): any;

  /**
   * Get the UI state object for the piano roll.
   * @returns MainEditorView
   */
  getMainEditor(): MainEditorView;

  /**
   * Get the phonemes for all notes in a group (passed in as a group reference). The group must be part of the currently open project.
   * @returns any[]
   */
  getPhonemesForGroup(): any[];

  /**
   * Get the UI state object for controlling the playback.
   * @returns PlayBackControl
   */
  getPlayback(): PlayBackControl;

  /**
   * Get the currently open project.
   * @returns Project
   */
  getProject(): Project;

  /**
   * Convert a MIDI number (semitones, where C4 is 60) to a frequency in Hz.
   * @param p
   * @returns number
   */
  pitch2freq(p: number): number;

  /**
   * Print any number of arguments to the standard output stream.
   */
  print(): void;

  /**
   * Convert `q` from number of quarters into number of blick.
   * @param q
   * @returns number
   */
  quarter2Blick(q: number): number;

  /**
   * Force Synthesizer V Studio to reload the side panel section for the current script.
   */
  refreshSidePanel(): void;

  /**
   * Convert `s` from seconds into blicks with the specified beats per minute `bpm`.
   * @param s
   * @param bpm
   * @returns number
   */
  seconds2Blick(s: number, bpm: number): number;

  /**
   * Set the system clipboard.
   * @param text
   */
  setHostClipboard(text: string): void;

  /**
   * Schedule a delayed call to `callback` after `timeOut` milliseconds.
   * @param timeOut
   * @param callback
   */
  setTimeout(timeOut: number, callback: Function): void;

  /**
   * The synchronous version of `SV#showCustomDialogAsync` that blocks the script execution until the user closes the dialog. It returns the inputs (the completed form) from the user.
   * @param form
   * @returns any
   */
  showCustomDialog(form: any): any;

  /**
   * Display a custom dialog defined in `form`, without blocking the script execution.
   * @param form
   * @param callback
   */
  showCustomDialogAsync(form: any, callback: Function): void;

  /**
   * The synchronous version of `SV#showInputBoxAsync` that blocks the script execution until the user closes the dialog. It returns the text input from the user.
   * @param title
   * @param message
   * @param defaultText
   */
  showInputBox(title: string, message: string, defaultText: string): void;

  /**
   * Display a dialog with a text box and an "OK" button, without blocking the script execution.
   * @param title
   * @param message
   * @param defaultText
   * @param callback
   */
  showInputBoxAsync(title: string, message: string, defaultText: string, callback: Function): void;

  /**
   * The synchronous version of `SV#showMessageBoxAsync` that blocks the script execution until the user closes the message box.
   * @param title
   * @param message
   */
  showMessageBox(title: string, message: string): void;

  /**
   * Cause a message box to pop up without blocking the script execution.
   * @param title
   * @param message
   */
  showMessageBoxAsync(title: string, message: string): void;

  /**
   * The synchronous version of `SV#showOkCancelBoxAsync` that blocks the script execution until the user closes the message box. It returns true if "OK" button is pressed.
   * @param title
   * @param message
   * @returns boolean
   */
  showOkCancelBox(title: string, message: string): boolean;

  /**
   * Display a message box with an "OK" button and a "Cancel" button, without blocking the script execution.
   * @param title
   * @param message
   * @param callback
   */
  showOkCancelBoxAsync(title: string, message: string, callback: Function): void;

  /**
   * The synchronous version of `SV#showYesNoCancelBoxAsync` that blocks the script execution until the user closes the message box. It returns "yes", "no" or "cancel".
   * @param title
   * @param message
   * @returns string
   */
  showYesNoCancelBox(title: string, message: string): string;

  /**
   * Display a message box with a "Yes" button, an "No" button and a "Cancel" button, without blocking the script execution.
   * @param title
   * @param message
   * @param callback
   */
  showYesNoCancelBoxAsync(title: string, message: string, callback: Function): void;

}

declare class ScriptableNestedObject {
  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class SelectionStateBase {
  /**
   * Unselects all object types supported by this selection state. Return true if the selection has changed.
   * @returns boolean
   */
  clearAll(): boolean;

  /**
   * Check if there's anything selected.
   */
  hasSelectedContent(): void;

  /**
   * Check if there's any unfinished edit on the selected objects.
   */
  hasUnfinishedEdits(): void;

  /**
   * Attach a script function to be called when the selection is cleared. The callback function will receive one argument: the type of the cleared objects.
   * @param callback
   */
  registerClearCallback(callback: Function): void;

  /**
   * Attach a script function to be called when objects are selected or deselected by the user. The callback function will receive two arguments: the type of the selected objects and whether this is a selection or deselection operation.
   * @param callback
   */
  registerSelectionCallback(callback: Function): void;

}

declare class TimeAxis {
  /**
   * Insert a `nomin`/`denom` measure mark at position `measure` (a measure number). If a measure mark exists at `measure`, update the information.
   * @param measure
   * @param nomin
   * @param denom
   */
  addMeasureMark(measure: number, nomin: number, denom: number): void;

  /**
   * Insert a tempo mark with beats per minute of `bpm` at position `b` (blicks). If a tempo mark exists at position `b`, update the BPM.
   * @param b
   * @param bpm
   */
  addTempoMark(b: number, bpm: number): void;

  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns TimeAxis
   */
  clone(): TimeAxis;

  /**
   * Get all measure marks in this `TimeAxis`. See `TimeAxis#getMeasureMarkAt`.
   */
  getAllMeasureMarks(): void;

  /**
   * Get all tempo marks in this `TimeAxis`. See `TimeAxis#getTempoMarkAt`.
   * @returns any[]
   */
  getAllTempoMarks(): any[];

  /**
   * Convert physical time `t` (second) to musical time (blicks).
   * @param t
   */
  getBlickFromSeconds(t: number): void;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the measure number at position `b` (blicks).
   * @param b
   * @returns number
   */
  getMeasureAt(b: number): number;

  /**
   * Get the measure mark at measure `measureNumber`.
   * @param measureNumber
   * @returns any
   */
  getMeasureMarkAt(measureNumber: number): any;

  /**
   * Get the measure mark that is effective at position `b` (blicks). For the returned object, see `TimeAxis#getMeasureMarkAt`.
   * @param b
   * @returns any
   */
  getMeasureMarkAtBlick(b: number): any;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Convert musical time `b` (blicks) to physical time (seconds).
   * @param b
   */
  getSecondsFromBlick(b: number): void;

  /**
   * Get the tempo mark that is effective at position `b` (blicks).
   * @param b
   * @returns TempoMark
   */
  getTempoMarkAt(b: number): TempoMark;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove the measure mark at measure number `measure`. If a measure mark exists at `measure`, return true.
   * @param measure
   * @returns boolean
   */
  removeMeasureMark(measure: number): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Remove the tempo mark at position `b` (blicks). If a tempo mark exists at position `b`, return true.
   * @param b
   * @returns boolean
   */
  removeTempoMark(b: number): boolean;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class Track {
  /**
   * Add a `NoteGroupReference` to this `Track` and return the index of the added group. It keeps all groups sorted by onset position.
   * @returns number
   */
  addGroupReference(): number;

  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * A deep copy of the current object.
   * @returns Track
   */
  clone(): Track;

  /**
   * Get the track's color as a hex string.
   * @returns string
   */
  getDisplayColor(): string;

  /**
   * Get the display order of the track inside the parent `Project`. A track's display order can be different from its storage index. The order of tracks as displayed in arrangement view is always based on the display order.
   * @returns number
   */
  getDisplayOrder(): number;

  /**
   * Get the duration of the `Track` in blicks, defined as the ending position of the last `NoteGroupReference`.
   * @returns number
   */
  getDuration(): number;

  /**
   * Get the `index`-th `NoteGroupReference`. The first is always the main group, followed by groups that refer to `NoteGroup` in the project library. The groups are sorted in ascending onset positions.
   * @param index
   * @returns NoteGroupReference
   */
  getGroupReference(index: number): NoteGroupReference;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the track's mixer.
   * @returns TrackMixer
   */
  getMixer(): TrackMixer;

  /**
   * Get the track name.
   * @returns string
   */
  getName(): string;

  /**
   * Get the number of `NoteGroupReference` in this `Track`, including the main group.
   * @returns number
   */
  getNumGroups(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * An option for whether or not to be exported to files, shown in Render Panel.
   * @returns boolean
   */
  isBounced(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Remove the `index`-th `NoteGroupReference` from this `Track`.
   * @param index
   */
  removeGroupReference(index: number): void;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set whether or not to have the `Track` exported to files. See `Track#isBounced`.
   * @param enabled
   */
  setBounced(enabled: boolean): void;

  /**
   * Set the display color of the `Track` to a hex string.
   * @param colorStr
   */
  setDisplayColor(colorStr: string): void;

  /**
   * Set the name of the `Track`.
   * @param name
   */
  setName(name: string): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

}

declare class TrackInnerSelectionState {
  /**
   * Unselects all object types supported by this selection state. Return true if the selection has changed.
   * @returns boolean
   */
  clearAll(): boolean;

  /**
   * Unselect all `NoteGroupReference`. Return true if the selection has changed.
   * @returns boolean
   */
  clearGroups(): boolean;

  /**
   * Unselect all notes. Return `true` if the selection has changed.
   * @returns boolean
   */
  clearNotes(): boolean;

  /**
   * Clear all pitch control selections.
   */
  clearPitchControls(): void;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Get an array of selected `NoteGroupReference` following the order of selection.
   */
  getSelectedGroups(): void;

  /**
   * Get an array of selected `Note` following the order of selection.
   * @returns any[]
   */
  getSelectedNotes(): any[];

  /**
   * Get an array of selected pitch control objects.
   * @returns any[]
   */
  getSelectedPitchControls(): any[];

  /**
   * Get an array of selected automation points for the specified parameter type.
   * @returns any[]
   */
  getSelectedPoints(): any[];

  /**
   * Check if there's anything selected.
   */
  hasSelectedContent(): void;

  /**
   * Check if there is at least one `NoteGroupReference` selected.
   */
  hasSelectedGroups(): void;

  /**
   * Check if there is at least one `Note` selected.
   * @returns boolean
   */
  hasSelectedNotes(): boolean;

  /**
   * Check if there are any selected pitch control objects.
   * @returns boolean
   */
  hasSelectedPitchControls(): boolean;

  /**
   * Check if there's any unfinished edit on the selected objects.
   */
  hasUnfinishedEdits(): void;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Attach a script function to be called when the selection is cleared. The callback function will receive one argument: the type of the cleared objects.
   * @param callback
   */
  registerClearCallback(callback: Function): void;

  /**
   * Attach a script function to be called when objects are selected or deselected by the user. The callback function will receive two arguments: the type of the selected objects and whether this is a selection or deselection operation.
   * @param callback
   */
  registerSelectionCallback(callback: Function): void;

  /**
   * Add a `NoteGroupReference` to the selection.
   */
  selectGroup(): void;

  /**
   * Select a `Note`. The note must be inside the current `NoteGroupReference` opened in the piano roll (see `MainEditorView#getCurrentGroup`).
   */
  selectNote(): void;

  /**
   * Select pitch control objects.
   */
  selectPitchControls(): void;

  /**
   * Select automation points for the specified parameter type.
   */
  selectPoints(): void;

  /**
   * Unselect a `NoteGroupReference`. Return true if the selection has changed.
   * @returns boolean
   */
  unselectGroup(): boolean;

  /**
   * Unselect a `Note`. Return `true` if the selection has changed.
   * @returns boolean
   */
  unselectNote(): boolean;

  /**
   * Unselect pitch control objects.
   */
  unselectPitchControls(): void;

  /**
   * Unselect automation points for the specified parameter type.
   */
  unselectPoints(): void;

}

declare class TrackMixer {
  /**
   * Remove all script data from the object's storage. Note: use with caution as this could also remove data created by other scripts.
   */
  clearScriptData(): void;

  /**
   * Get the gain in decibels.
   * @returns number
   */
  getGainDecibel(): number;

  /**
   * Get index of the current object in its parent. In Lua, this index starts from 1. In JavaScript, this index starts from 0.
   * @returns number
   */
  getIndexInParent(): number;

  /**
   * Get the pan position.
   * @returns number
   */
  getPan(): number;

  /**
   * Get the parent `NestedObject`. Return `undefined` if the current object is not attached to a parent.
   * @returns NestedObject|undefined
   */
  getParent(): NestedObject|undefined;

  /**
   * Retrieve a value from the object's script data storage by key. Returns `undefined` if the key does not exist.
   * @returns any
   */
  getScriptData(): any;

  /**
   * Get all keys currently stored in the object's script data storage.
   */
  getScriptDataKeys(): void;

  /**
   * Check whether a key exists in the object's script data storage.
   * @returns boolean
   */
  hasScriptData(): boolean;

  /**
   * Check whether or not the current object is memory managed (i.e. garbage collected by the script environment).
   * @returns boolean
   */
  isMemoryManaged(): boolean;

  /**
   * Check if the track is muted.
   * @returns boolean
   */
  isMuted(): boolean;

  /**
   * Check if the track is soloed.
   * @returns boolean
   */
  isSolo(): boolean;

  /**
   * Remove a key-value pair from the object's script data storage.
   */
  removeScriptData(): void;

  /**
   * Set the gain in decibels.
   */
  setGainDecibel(): void;

  /**
   * Set the mute state of the track.
   */
  setMuted(): void;

  /**
   * Set the pan position.
   */
  setPan(): void;

  /**
   * Store a value with the specified key in the object's script data storage. The value must be JSON-serializable.
   */
  setScriptData(): void;

  /**
   * Set the solo state of the track.
   */
  setSolo(): void;

}

declare class WidgetValue {
  /**
   * Get the enable/disable status currently set for the UI widget.
   * @returns boolean
   */
  getEnabled(): boolean;

  /**
   * Get the value currently set for the UI widget.
   * @returns any
   */
  getValue(): any;

  /**
   * Enable or disable the UI widget. In the disabled state, the user will not be able to change the widget's value.
   * @param enabled
   */
  setEnabled(enabled: boolean): void;

  /**
   * Update the UI widget with a new value.
   */
  setValue(): void;

  /**
   * Set a script function to be called when the UI widget's value is changed by the user. The callback function will receive the new value as its sole argument.
   * @param callback
   */
  setValueChangeCallback(callback: Function): void;

}
