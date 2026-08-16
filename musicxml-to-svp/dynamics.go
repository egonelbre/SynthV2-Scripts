package main

import (
	"slices"
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
					position := cursor
					// fp and friends: a brief attack, then the settled level.
					if lvl.attack != nil {
						events = append(events, dynEvent{
							position: position,
							kind:     dynLevel,
							loudness: lvl.attack.loudness,
							tension:  lvl.attack.tension,
						})
						position += attackBlicks + attackDecayBlicks
					}
					events = append(events, dynEvent{
						position:   position,
						kind:       dynLevel,
						loudness:   lvl.loudness,
						tension:    lvl.tension,
						transition: attackDecay(lvl),
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
	// transition overrides how long the ramp into this level takes;
	// 0 means the default step transition.
	transition int64
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
	// Default to mf level if no explicit level precedes the first hairpin.
	currentLevel := getValue(dynEvent{loudness: 1.5, tension: 0})
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
		startPos int64
		stopPos  int64
		kind     dynEventKind
		stopIdx  int
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

	usedRanges := make([]bool, len(ranges))

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
				transition := stepTransitionBlicks
				if ev.transition > 0 {
					transition = ev.transition
				}
				transitionStart := ev.position - transition
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

			for ri, wr := range ranges {
				if usedRanges[ri] {
					continue
				}
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
					usedRanges[ri] = true
					found = true
					break
				}
			}

			if !found {
				// Unpaired cresc/dim text — use next dynLevel as endpoint
				// if it's within 4 measures, otherwise estimate over 2 measures.
				maxRange := int64(4 * 4 * blicksPerQuarter) // 4 measures of 4/4
				foundNearby := false
				for j := i + 1; j < len(events); j++ {
					if events[j].kind == dynLevel {
						if events[j].position-ev.position <= maxRange {
							targetLevel = getValue(events[j])
							stopPos = events[j].position
							foundNearby = true
						}
						break
					}
				}
				if !foundNearby {
					stopPos = ev.position + 2*4*blicksPerQuarter
					if ev.kind == dynCrescStart {
						targetLevel = currentLevel + defaultDelta
					} else {
						targetLevel = currentLevel - defaultDelta
					}
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
	// attack is the brief louder level struck before settling to this one,
	// as in fp (forte-piano) and its sforzando variants.
	attack *dynamicLevel
}

// How long the attack of an fp-style dynamic is held before dropping.
const (
	attackBlicks      = int64(blicksPerQuarter) // hold the attack for a quarter note
	attackDecayBlicks = int64(blicksPerQuarter) // then ease down over another quarter
)

// attackDecay returns the ramp length into a level that follows an attack.
func attackDecay(lvl dynamicLevel) int64 {
	if lvl.attack == nil {
		return 0
	}
	return attackDecayBlicks
}

var (
	attackF  = &dynamicLevel{loudness: 3, tension: 0}
	attackSF = &dynamicLevel{loudness: 3, tension: 0.3}
)

// dynamicLevels maps MusicXML dynamics element names to loudness (dB) and tension values.
//
// Loudness is kept within -6 to 6 dB. For extreme dynamics (pp and softer,
// ff and louder), loudness stays close to p/f range while tension is adjusted
// to convey the additional intensity difference.
var dynamicLevels = map[string]dynamicLevel{
	// Fortissimo variants: loudness near f, tension increases.
	"ffffff": {loudness: 6, tension: 0.8},
	"fffff":  {loudness: 5.5, tension: 0.7},
	"ffff":   {loudness: 5, tension: 0.5},
	"fff":    {loudness: 4.5, tension: 0.4},
	"ff":     {loudness: 4, tension: 0.2},

	// Pianissimo variants: loudness near p, tension decreases.
	"pppppp": {loudness: -6, tension: -0.8},
	"ppppp":  {loudness: -5.5, tension: -0.7},
	"pppp":   {loudness: -5, tension: -0.5},
	"ppp":    {loudness: -4.5, tension: -0.4},
	"pp":     {loudness: -4, tension: -0.2},

	// Sforzando variants. The *p forms attack, then drop to piano.
	"sffz": {loudness: 3, tension: 0.3},
	"sfzp": {loudness: -3, tension: 0, attack: attackSF},
	"sfpp": {loudness: -4, tension: -0.2, attack: attackSF},
	"sfz":  {loudness: 3, tension: 0.3},
	"sfp":  {loudness: -3, tension: 0, attack: attackSF},
	"sf":   {loudness: 3, tension: 0.3},

	// Core dynamics.
	"mp": {loudness: -1.5, tension: 0},
	"mf": {loudness: 1.5, tension: 0},
	"fp": {loudness: -3, tension: 0, attack: attackF},
	"fz": {loudness: 3, tension: 0.3},
	"f":  {loudness: 3, tension: 0},

	"rfz": {loudness: 1.5, tension: 0.2},
	"rf":  {loudness: 1.5, tension: 0.2},

	"pf": {loudness: 0, tension: 0},
	"p":  {loudness: -3, tension: 0},

	"n": {loudness: -6, tension: -0.8},
}

// dynamicsToLevel maps a MusicXML dynamics element to loudness (dB) and tension values.
// It skips non-dynamic elements like other-dynamics (e.g., "sempre") to find the
// actual dynamic marking.
func dynamicsToLevel(d *musicxml.Dynamics) (dynamicLevel, bool) {
	s := d.InnerXML
	for {
		name, rest := nextXMLElementName(s)
		if name == "" {
			return dynamicLevel{}, false
		}
		if lvl, ok := dynamicLevels[name]; ok {
			return lvl, true
		}
		s = rest
	}
}

// nextXMLElementName extracts the tag name of the next XML element in s,
// returning the name and the remaining string after the element.
// For example, "<ff/>" returns ("ff", ""), "<p default-x=\"10\"/>" returns ("p", "").
// Skips XML comments (<!-- ... -->) and processing instructions (<? ... ?>).
func nextXMLElementName(s string) (string, string) {
	for {
		start := strings.Index(s, "<")
		if start < 0 {
			return "", ""
		}
		start++ // skip '<'
		if start >= len(s) {
			return "", ""
		}
		// Skip comments.
		if strings.HasPrefix(s[start:], "!--") {
			end := strings.Index(s[start:], "-->")
			if end < 0 {
				return "", ""
			}
			s = s[start+end+3:]
			continue
		}
		// Skip processing instructions.
		if s[start] == '?' {
			end := strings.Index(s[start:], "?>")
			if end < 0 {
				return "", ""
			}
			s = s[start+end+2:]
			continue
		}
		// Skip closing tags.
		if s[start] == '/' {
			end := strings.Index(s[start:], ">")
			if end < 0 {
				return "", ""
			}
			s = s[start+end+1:]
			continue
		}
		end := start
		for end < len(s) && s[end] != ' ' && s[end] != '>' && s[end] != '/' {
			end++
		}
		// Find the end of this element (closing '>') to compute the rest.
		closeIdx := strings.Index(s[end:], ">")
		rest := ""
		if closeIdx >= 0 {
			rest = s[end+closeIdx+1:]
		}
		return s[start:end], rest
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
		// Don't let the spike extend past the note end.
		if spikeWidth > acc.duration {
			spikeWidth = acc.duration
		}

		// Find the current curve value at the accent position.
		baseVal := curveValueAt(points, acc.position)

		// Insert spike: peak at onset, decay back to base.
		points = insertCurvePoints(points,
			float64(acc.position), baseVal+bump,
			float64(acc.position+spikeWidth), baseVal)
	}
	return points
}

// applyStresses layers a short bump at the start of each stressed note on top
// of the dynamics curve. Each shape starts and ends at zero offset, so the
// curve is untouched between words.
func applyStresses(points []float64, stresses []accentEvent, bump float64) []float64 {
	shapes := make([][]float64, 0, len(stresses))
	for _, s := range stresses {
		width := min(s.duration/2, maxStressWidth)
		width = max(width, minAccentSpikeWidth)
		width = min(width, s.duration)

		shapes = append(shapes, []float64{
			float64(s.position - stressAnchorGap), 0,
			float64(s.position), bump,
			float64(s.position + width), 0,
		})
	}
	return addStressShapes(points, shapes)
}

// applyXyloStresses layers a struck-instrument shape on each stressed note: a
// strong attack at the word start, a fast decay, then a fade below the
// dynamics level that is held until the note ends.
func applyXyloStresses(points []float64, stresses []accentEvent, bump float64) []float64 {
	shapes := make([][]float64, 0, len(stresses))
	for _, s := range stresses {
		attack := min(s.duration/accentSpikeWidthFraction, maxStressWidth)
		attack = max(attack, minAccentSpikeWidth)
		attack = min(attack, s.duration)

		shape := []float64{
			float64(s.position - stressAnchorGap), 0,
			float64(s.position), bump,
			float64(s.position + attack), bump * xyloDecayFraction,
		}
		// Fade to the floor within a quarter note, hold it until the note ends,
		// then return to the dynamics level for the next word.
		if tailEnd := s.duration - 2*stressAnchorGap; tailEnd > 2*attack {
			floor := bump * xyloTailFraction
			decayEnd := min(attack+xyloDecaySpan, tailEnd)
			shape = append(shape, float64(s.position+decayEnd), floor)
			if decayEnd < tailEnd {
				shape = append(shape, float64(s.position+tailEnd), floor)
			}
		}
		shape = append(shape, float64(s.position+s.duration), 0)
		shapes = append(shapes, shape)
	}
	return addStressShapes(points, shapes)
}

// addStressShapes adds stress offset shapes to a curve. Every position from
// either the curve or the shapes gets a point valued base+offset, so dynamics
// under a stress (a diminuendo, say) survive instead of being overwritten.
func addStressShapes(points []float64, shapes [][]float64) []float64 {
	var offsets []float64
	for _, shape := range shapes {
		// Clip the tail of the previous shape when notes are close enough that
		// the shapes overlap; the earlier trail-off runs into the new attack.
		for len(offsets) >= 2 && offsets[len(offsets)-2] >= shape[0] {
			offsets = offsets[:len(offsets)-2]
		}
		offsets = append(offsets, shape...)
	}
	if len(offsets) == 0 {
		return points
	}

	positions := make([]float64, 0, len(points)/2+len(offsets)/2)
	for i := 0; i < len(points); i += 2 {
		positions = append(positions, points[i])
	}
	for i := 0; i < len(offsets); i += 2 {
		positions = append(positions, offsets[i])
	}
	slices.Sort(positions)
	positions = slices.Compact(positions)

	// ponytail: linear scan per position, fine for a few thousand points.
	result := make([]float64, 0, 2*len(positions))
	for _, pos := range positions {
		result = append(result, pos, curveValueAt(points, int64(pos))+offsetAt(offsets, pos))
	}
	return result
}

// offsetAt linearly interpolates a stress offset curve, returning 0 outside it.
func offsetAt(offsets []float64, pos float64) float64 {
	if pos <= offsets[0] || pos >= offsets[len(offsets)-2] {
		return 0
	}
	for i := 0; i+3 < len(offsets); i += 2 {
		p1, v1, p2, v2 := offsets[i], offsets[i+1], offsets[i+2], offsets[i+3]
		if pos < p1 || pos > p2 {
			continue
		}
		if p2 == p1 {
			return v2
		}
		return v1 + (v2-v1)*(pos-p1)/(p2-p1)
	}
	return 0
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

// insertCurvePoints inserts points, given as pos/value pairs in ascending
// position order, into a sorted curve point array. Existing points within the
// inserted range are replaced to avoid duplicates.
func insertCurvePoints(points []float64, newPts ...float64) []float64 {
	if len(points) == 0 {
		return newPts
	}

	fpos1 := newPts[0]
	fpos2 := newPts[len(newPts)-2]

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
