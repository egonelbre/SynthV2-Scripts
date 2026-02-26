package main

import (
	"strings"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

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
				if wr.stopPos == ev.position || (ev.position-wr.stopPos >= 0 && ev.position-wr.stopPos < stepTransitionBlicks*2) {
					if wr.stopIdx == i-1 || (i > 0 && events[i-1].kind == dynWedgeStop) {
						isWedgeTarget = true
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

type dynLevel2 struct {
	loudness float64
	tension  float64
}

// dynamicsToLevel maps a MusicXML dynamics element to loudness (dB) and tension values.
//
// Loudness is kept within -12 to 12 dB. For extreme dynamics (pp and softer,
// ff and louder), loudness stays close to p/f range while tension is adjusted
// to convey the additional intensity difference.
func dynamicsToLevel(d *musicxml.Dynamics) (dynLevel2, bool) {
	xml := d.InnerXML

	// Check from most specific to least specific to avoid prefix matching issues.
	switch {
	// Fortissimo variants: loudness near f, tension increases.
	case strings.Contains(xml, "<ffffff"):
		return dynLevel2{12, 0.8}, true
	case strings.Contains(xml, "<fffff"):
		return dynLevel2{11, 0.7}, true
	case strings.Contains(xml, "<ffff"):
		return dynLevel2{10, 0.5}, true
	case strings.Contains(xml, "<fff"):
		return dynLevel2{9, 0.4}, true
	case strings.Contains(xml, "<ff"):
		return dynLevel2{8, 0.2}, true

	// Pianissimo variants: loudness near p, tension decreases.
	case strings.Contains(xml, "<pppppp"):
		return dynLevel2{-12, -0.8}, true
	case strings.Contains(xml, "<ppppp"):
		return dynLevel2{-11, -0.7}, true
	case strings.Contains(xml, "<pppp"):
		return dynLevel2{-10, -0.5}, true
	case strings.Contains(xml, "<ppp"):
		return dynLevel2{-9, -0.4}, true
	case strings.Contains(xml, "<pp"):
		return dynLevel2{-8, -0.2}, true

	// Sforzando variants.
	case strings.Contains(xml, "<sffz"):
		return dynLevel2{6, 0.3}, true
	case strings.Contains(xml, "<sfzp"):
		return dynLevel2{3, 0}, true
	case strings.Contains(xml, "<sfpp"):
		return dynLevel2{3, 0}, true
	case strings.Contains(xml, "<sfz"):
		return dynLevel2{6, 0.3}, true
	case strings.Contains(xml, "<sfp"):
		return dynLevel2{3, 0}, true
	case strings.Contains(xml, "<sf"):
		return dynLevel2{6, 0.3}, true

	// Core dynamics.
	case strings.Contains(xml, "<mp"):
		return dynLevel2{-3, 0}, true
	case strings.Contains(xml, "<mf"):
		return dynLevel2{3, 0}, true

	case strings.Contains(xml, "<fp"):
		return dynLevel2{0, 0}, true
	case strings.Contains(xml, "<fz"):
		return dynLevel2{6, 0.3}, true
	case strings.Contains(xml, "<f"):
		return dynLevel2{6, 0}, true

	case strings.Contains(xml, "<rfz"):
		return dynLevel2{3, 0.2}, true
	case strings.Contains(xml, "<rf"):
		return dynLevel2{3, 0.2}, true

	case strings.Contains(xml, "<pf"):
		return dynLevel2{0, 0}, true
	case strings.Contains(xml, "<p"):
		return dynLevel2{-6, 0}, true

	case strings.Contains(xml, "<n"):
		return dynLevel2{-12, -0.8}, true
	}
	return dynLevel2{}, false
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
