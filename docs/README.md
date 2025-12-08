# Synthesizer V Studio Scripting Tools

Tools for working with the Dreamtonics Synthesizer V Studio Scripting API.

## Scripts

### download_docs.py

Downloads all documentation from the official Dreamtonics Scripting API website.

**Usage:**
```bash
python3 download_docs.py
```

**Output:**
- Downloads 24 HTML files to `docs/dreamtonics-api/`
- Includes the main index page and all 23 class documentation pages

**Features:**
- Respectful rate limiting (0.5s delay between requests)
- Progress feedback with success/failure indicators
- Error handling for network issues

### generate_types.py

Parses the downloaded HTML documentation and generates TypeScript type definitions.

**Usage:**
```bash
python3 generate_types.py
```

**Prerequisites:**
- Run `download_docs.py` first to download the documentation

**Output:**
- Generates `synthesizer-v-api.d.ts` with TypeScript type definitions
- 23 classes with 371+ methods
- ~2,200 lines of TypeScript definitions

**Features:**
- Extracts method signatures with parameter types
- Preserves documentation comments
- Handles return types and inheritance
- Converts documentation types to TypeScript types

## Generated Files

### synthesizer-v-api.d.ts

TypeScript type definitions for the Synthesizer V Studio Scripting API. This file provides:

- Type checking for all API classes
- IntelliSense/autocomplete support in VS Code and other editors
- Method signatures with parameter and return types
- JSDoc comments with descriptions

**Example usage in a TypeScript project:**
```typescript
/// <reference path="./synthesizer-v-api.d.ts" />

function main() {
  const project: Project = SV.getProject();
  const track: Track = project.getTrack(0);
  const group: NoteGroup = track.getGroupReference(0).getTarget();
  
  for (let i = 0; i < group.getNumNotes(); i++) {
    const note: Note = group.getNote(i);
    console.log(`Note: ${note.getLyrics()} at ${note.getOnset()}`);
  }
}
```

## API Documentation

The scripts work with the official Dreamtonics Scripting API documentation:
https://resource.dreamtonics.com/scripting/index.html

### Documented Classes

- ArrangementSelectionState
- ArrangementView
- Automation
- CoordinateSystem
- GroupSelection
- MainEditorView
- NestedObject
- Note
- NoteGroup
- NoteGroupReference
- PitchControlCurve
- PitchControlPoint
- PlaybackControl
- Project
- RetakeList
- SV (main API entry point)
- ScriptableNestedObject
- SelectionStateBase
- TimeAxis
- Track
- TrackInnerSelectionState
- TrackMixer
- WidgetValue

## Requirements

- Python 3.6+
- No external dependencies (uses only standard library)

## License

These tools are for working with the Dreamtonics Synthesizer V Studio Scripting API.
Refer to the official Synthesizer V Studio documentation for API license terms.
