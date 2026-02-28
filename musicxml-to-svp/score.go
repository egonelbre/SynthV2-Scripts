package main

// Score represents played music with all repeats unrolled and timing resolved.
type Score struct {
	Meters []MeterChange
	Tempos []TempoChange
	Parts  []Part
}

type MeterChange struct {
	MeasureIndex int
	Numerator    int
	Denominator  int
}

type TempoChange struct {
	Position int64   // blicks
	BPM      float64
}

type Part struct {
	Name     string
	Notes    []Note
	Dynamics []dynEvent
}

type Articulation int

const (
	ArticulationStaccato      Articulation = 1 << iota
	ArticulationStaccatissimo
	ArticulationTenuto
	ArticulationAccent
	ArticulationStrongAccent // marcato
)

type Note struct {
	Onset         int64        // blicks, absolute
	Duration      int64        // blicks (original notated duration)
	Pitch         int          // MIDI number
	Detune        int          // cents
	Lyric         string
	Articulations Articulation // bitmask

	LeadingGraces  []GraceNote // grace notes before this note
	TrailingGraces []GraceNote // grace notes after this note
}

type GraceNote struct {
	Pitch        int    // MIDI number
	Detune       int    // cents
	Lyric        string
	NotatedType  string // "quarter", "eighth", "16th", etc.
	Acciaccatura bool   // slash grace = halved duration
	Duration     int64  // blicks, computed during buildNotes
}
