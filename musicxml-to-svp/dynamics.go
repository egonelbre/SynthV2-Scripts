package main

import (
	"strings"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

// buildDynamics collects dynamic events from directions in a part.
func buildDynamics(part *musicxml.Part, unrolled []playedMeasure) []dynEvent {
	var events []dynEvent

	walkPartElements(part, unrolled, func(cursor int64, divisions int, pm playedMeasure, value any) {
		dir, ok := value.(*musicxml.Direction)
		if !ok {
			return
		}
		for _, dt := range dir.DirectionType {
			for _, dyn := range dt.Dynamics {
				if lvl, ok := dynamicsToLevel(dyn); ok {
					events = append(events, dynEvent{
						position: cursor,
						kind:     dynLevel,
						loudness: lvl.loudness,
						tension:  lvl.tension,
					})
				}
			}
			if dt.Wedge != nil {
				num := dt.Wedge.Number
				if num == 0 {
					num = 1
				}
				switch dt.Wedge.Type {
				case "crescendo":
					events = append(events, dynEvent{
						position: cursor,
						kind:     dynCrescStart,
						number:   num,
					})
				case "diminuendo":
					events = append(events, dynEvent{
						position: cursor,
						kind:     dynDimStart,
						number:   num,
					})
				case "stop":
					events = append(events, dynEvent{
						position: cursor,
						kind:     dynWedgeStop,
						number:   num,
					})
				}
			}
			for _, w := range dt.Words {
				if isTextCresc(w.EnclosedText) {
					events = append(events, dynEvent{
						position: cursor,
						kind:     dynCrescStart,
						number:   -1,
					})
				} else if isTextDim(w.EnclosedText) {
					events = append(events, dynEvent{
						position: cursor,
						kind:     dynDimStart,
						number:   -1,
					})
				}
			}
		}
	})

	return events
}

// Dynamic event types for building loudness curves.
type dynEventKind int

const (
	dynLevel      dynEventKind = iota // instant dynamic level (p, mf, f, ...)
	dynCrescStart                     // start of crescendo hairpin or text cresc.
	dynDimStart                       // start of diminuendo hairpin or text dim.
	dynWedgeStop                      // end of hairpin
)

type dynEvent struct {
	position int64
	kind     dynEventKind
	loudness float64 // only meaningful for dynLevel
	tension  float64 // only meaningful for dynLevel
	number   int     // wedge number for matching start/stop
}

const stepTransitionBlicks = int64(blicksPerQuarter / 8) // 1/8 quarter note

// buildCurve builds an SVP parameter curve from dynamic events using the given
// value extractor. defaultDelta is used to estimate cresc/dim targets when no
// following dynLevel event exists.
func buildCurve(events []dynEvent, getValue func(dynEvent) float64, defaultDelta float64) []float64 {
	if len(events) == 0 {
		return nil
	}

	var points []float64
	currentLevel := 0.0
	hasLevel := false

	// Pair up wedge starts and stops by number.
	type wedgeInfo struct {
		startIdx int
		stopIdx  int
		kind     dynEventKind // dynCrescStart or dynDimStart
	}
	openWedges := map[int]int{} // number -> index in events
	var wedges []wedgeInfo

	for i, ev := range events {
		switch ev.kind {
		case dynCrescStart, dynDimStart:
			openWedges[ev.number] = i
		case dynWedgeStop:
			if startIdx, ok := openWedges[ev.number]; ok {
				wedges = append(wedges, wedgeInfo{
					startIdx: startIdx,
					stopIdx:  i,
					kind:     events[startIdx].kind,
				})
				delete(openWedges, ev.number)
			}
		}
	}

	// Build a set of wedge ranges for quick lookup.
	type wedgeRange struct {
		startPos  int64
		stopPos   int64
		kind      dynEventKind
		stopIdx   int
	}
	var ranges []wedgeRange
	for _, w := range wedges {
		ranges = append(ranges, wedgeRange{
			startPos: events[w.startIdx].position,
			stopPos:  events[w.stopIdx].position,
			kind:     w.kind,
			stopIdx:  w.stopIdx,
		})
	}

	// Find the next dynLevel event at or after a given index.
	findNextLevel := func(fromIdx int) (float64, bool) {
		for j := fromIdx; j < len(events); j++ {
			if events[j].kind == dynLevel {
				return getValue(events[j]), true
			}
		}
		return 0, false
	}

	addPoint := func(pos int64, val float64) {
		points = append(points, float64(pos), val)
	}

	for i, ev := range events {
		evVal := getValue(ev)

		switch ev.kind {
		case dynLevel:
			// Check if this level is the target of a just-ended wedge.
			// If so, the ramp already brought us here; just update currentLevel.
			isWedgeTarget := false
			for _, wr := range ranges {
				if ev.position-wr.stopPos >= 0 && ev.position-wr.stopPos < stepTransitionBlicks*2 {
					// Scan backwards to find if a wedge stop is among
					// the recent events at or near the same position.
					for k := i - 1; k >= 0 && events[k].position >= wr.stopPos; k-- {
						if events[k].kind == dynWedgeStop {
							isWedgeTarget = true
							break
						}
					}
				}
			}

			if isWedgeTarget {
				// End of ramp — just place the final point.
				addPoint(ev.position, evVal)
			} else if hasLevel {
				// Step transition: hold old level, then jump to new.
				transitionStart := ev.position - stepTransitionBlicks
				if transitionStart < 0 {
					transitionStart = 0
				}
				addPoint(transitionStart, currentLevel)
				addPoint(ev.position, evVal)
			} else {
				// First dynamic marking — set initial level.
				addPoint(ev.position, evVal)
			}

			currentLevel = evVal
			hasLevel = true

		case dynCrescStart, dynDimStart:
			// Find the matching stop and determine target level.
			var stopPos int64
			var targetLevel float64
			found := false

			for _, wr := range ranges {
				if wr.startPos == ev.position && wr.kind == ev.kind {
					stopPos = wr.stopPos
					// Look for a dynLevel right after the stop.
					if lvl, ok := findNextLevel(wr.stopIdx + 1); ok {
						targetLevel = lvl
					} else {
						// Estimate.
						if ev.kind == dynCrescStart {
							targetLevel = currentLevel + defaultDelta
						} else {
							targetLevel = currentLevel - defaultDelta
						}
					}
					found = true
					break
				}
			}

			if !found {
				// Unpaired cresc/dim text — estimate over 2 measures.
				stopPos = ev.position + 2*4*blicksPerQuarter
				if ev.kind == dynCrescStart {
					targetLevel = currentLevel + defaultDelta
				} else {
					targetLevel = currentLevel - defaultDelta
				}
			}

			// Emit ramp: start point at current level, end point at target.
			if hasLevel {
				addPoint(ev.position, currentLevel)
			}
			addPoint(stopPos, targetLevel)
			currentLevel = targetLevel

		case dynWedgeStop:
			// Already handled via ranges above.
		}
	}

	return points
}

type dynamicLevel struct {
	loudness float64
	tension  float64
}

// dynamicLevels maps MusicXML dynamics element names to loudness (dB) and tension values.
//
// Loudness is kept within -12 to 12 dB. For extreme dynamics (pp and softer,
// ff and louder), loudness stays close to p/f range while tension is adjusted
// to convey the additional intensity difference.
var dynamicLevels = map[string]dynamicLevel{
	// Fortissimo variants: loudness near f, tension increases.
	"ffffff": {12, 0.8},
	"fffff":  {11, 0.7},
	"ffff":   {10, 0.5},
	"fff":    {9, 0.4},
	"ff":     {8, 0.2},

	// Pianissimo variants: loudness near p, tension decreases.
	"pppppp": {-12, -0.8},
	"ppppp":  {-11, -0.7},
	"pppp":   {-10, -0.5},
	"ppp":    {-9, -0.4},
	"pp":     {-8, -0.2},

	// Sforzando variants.
	"sffz": {6, 0.3},
	"sfzp": {3, 0},
	"sfpp": {3, 0},
	"sfz":  {6, 0.3},
	"sfp":  {3, 0},
	"sf":   {6, 0.3},

	// Core dynamics.
	"mp": {-3, 0},
	"mf": {3, 0},
	"fp": {0, 0},
	"fz": {6, 0.3},
	"f":  {6, 0},

	"rfz": {3, 0.2},
	"rf":  {3, 0.2},

	"pf": {0, 0},
	"p":  {-6, 0},

	"n": {-12, -0.8},
}

// dynamicsToLevel maps a MusicXML dynamics element to loudness (dB) and tension values.
func dynamicsToLevel(d *musicxml.Dynamics) (dynamicLevel, bool) {
	name := firstXMLElementName(d.InnerXML)
	if lvl, ok := dynamicLevels[name]; ok {
		return lvl, true
	}
	return dynamicLevel{}, false
}

// firstXMLElementName extracts the tag name of the first XML element in s.
// For example, "<ff/>" returns "ff", "<p default-x=\"10\"/>" returns "p".
// Skips XML comments (<!-- ... -->) and processing instructions (<? ... ?>).
func firstXMLElementName(s string) string {
	for {
		start := strings.Index(s, "<")
		if start < 0 {
			return ""
		}
		start++ // skip '<'
		if start >= len(s) {
			return ""
		}
		// Skip comments.
		if strings.HasPrefix(s[start:], "!--") {
			end := strings.Index(s[start:], "-->")
			if end < 0 {
				return ""
			}
			s = s[start+end+3:]
			continue
		}
		// Skip processing instructions.
		if s[start] == '?' {
			end := strings.Index(s[start:], "?>")
			if end < 0 {
				return ""
			}
			s = s[start+end+2:]
			continue
		}
		end := start
		for end < len(s) && s[end] != ' ' && s[end] != '>' && s[end] != '/' {
			end++
		}
		return s[start:end]
	}
}

type accentEvent struct {
	position int64
	duration int64 // full note duration, used to scale spike width
	strong   bool  // strong accent (marcato) = bigger bump
}

// applyAccents overlays accent spikes onto an existing parameter curve.
// Each accent inserts a brief spike: a sharp rise at the note onset that
// decays over 1/4 of the note duration. normalBump is used for regular
// accents, strongBump for strong accents (marcato).
func applyAccents(points []float64, accents []accentEvent, normalBump, strongBump float64) []float64 {
	for _, acc := range accents {
		bump := normalBump
		if acc.strong {
			bump = strongBump
		}

		// Spike decays over a fraction of the note duration.
		spikeWidth := acc.duration / accentSpikeWidthFraction
		if spikeWidth < minAccentSpikeWidth {
			spikeWidth = minAccentSpikeWidth
		}

		// Find the current curve value at the accent position.
		baseVal := curveValueAt(points, acc.position)

		// Insert spike: peak at onset, decay back to base.
		points = insertCurvePoints(points, acc.position, baseVal+bump, acc.position+spikeWidth, baseVal)
	}
	return points
}

// curveValueAt returns the interpolated value of a curve at a given position.
// Uses Catmull-Rom cubic interpolation to match SVP's "cubic" curve mode.
func curveValueAt(points []float64, pos int64) float64 {
	if len(points) < 2 {
		return 0
	}
	fpos := float64(pos)

	// Before first point: use first value.
	if fpos <= points[0] {
		return points[1]
	}
	// After last point: use last value.
	if fpos >= points[len(points)-2] {
		return points[len(points)-1]
	}

	// Find the segment index (i) where points[i] <= fpos <= points[i+2].
	segIdx := 0
	for i := 0; i < len(points)-2; i += 2 {
		if fpos >= points[i] && fpos <= points[i+2] {
			segIdx = i
			break
		}
	}

	p1, v1 := points[segIdx], points[segIdx+1]
	p2, v2 := points[segIdx+2], points[segIdx+3]

	if p2 == p1 {
		return v2
	}

	// Catmull-Rom: get the neighboring points for tangent computation.
	// Clamp to endpoints if at the boundary.
	var p0, v0, p3, v3 float64
	if segIdx >= 2 {
		p0, v0 = points[segIdx-2], points[segIdx-1]
	} else {
		p0, v0 = p1, v1
	}
	if segIdx+4 < len(points) {
		p3, v3 = points[segIdx+4], points[segIdx+5]
	} else {
		p3, v3 = p2, v2
	}

	// Compute tangents at p1 and p2 using finite differences.
	dt := p2 - p1
	var m1, m2 float64
	if p2-p0 != 0 {
		m1 = (v2 - v0) / (p2 - p0) * dt
	}
	if p3-p1 != 0 {
		m2 = (v3 - v1) / (p3 - p1) * dt
	}

	// Hermite interpolation.
	t := (fpos - p1) / dt
	t2 := t * t
	t3 := t2 * t
	return (2*t3-3*t2+1)*v1 + (t3-2*t2+t)*m1 + (-2*t3+3*t2)*v2 + (t3-t2)*m2
}

// insertCurvePoints inserts two points (pos1,val1) and (pos2,val2) into a
// sorted curve point array, maintaining position order. Existing points at
// the same positions are replaced to avoid duplicates.
func insertCurvePoints(points []float64, pos1 int64, val1 float64, pos2 int64, val2 float64) []float64 {
	newPts := []float64{float64(pos1), val1, float64(pos2), val2}

	if len(points) == 0 {
		return newPts
	}

	fpos1 := float64(pos1)
	fpos2 := float64(pos2)

	// Find insertion index (before the first point >= pos1).
	idx := len(points)
	for i := 0; i < len(points); i += 2 {
		if points[i] >= fpos1 {
			idx = i
			break
		}
	}

	// Find end index: skip existing points within [pos1, pos2] to replace them.
	end := idx
	for end < len(points) && points[end] <= fpos2 {
		end += 2
	}

	result := make([]float64, 0, len(points)+4)
	result = append(result, points[:idx]...)
	result = append(result, newPts...)
	result = append(result, points[end:]...)
	return result
}

// isTextCresc checks if a words element indicates crescendo.
func isTextCresc(text string) bool {
	t := strings.ToLower(strings.TrimSpace(text))
	return strings.HasPrefix(t, "cresc") || strings.HasPrefix(t, "crésc")
}

// isTextDim checks if a words element indicates diminuendo/decrescendo.
func isTextDim(text string) bool {
	t := strings.ToLower(strings.TrimSpace(text))
	return strings.HasPrefix(t, "dim") || strings.HasPrefix(t, "decresc")
}
