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
	number   int     // wedge number for matching start/stop
}

// buildLoudnessCurve takes a chronological list of dynamic events and produces
// SVP loudness curve points [pos, val, pos, val, ...].
//
// Strategy:
//   - At each dynLevel event, emit a quick step transition (two points close
//     together) from the previous level to the new level.
//   - Between a cresc/dim start and the following wedge stop, emit a linear
//     ramp. The ramp target is the loudness of the next dynLevel event after
//     the stop (if one exists at the same position), otherwise we estimate
//     +6 dB for crescendo, -6 dB for diminuendo.
//   - Between other events, hold the current level (no points needed, cubic
//     interpolation holds the last value).
const stepTransitionBlicks = int64(blicksPerQuarter / 8) // 1/8 quarter note

func buildLoudnessCurve(events []dynEvent) []float64 {
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
				return events[j].loudness, true
			}
		}
		return 0, false
	}

	addPoint := func(pos int64, val float64) {
		points = append(points, float64(pos), val)
	}

	for i, ev := range events {
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
				addPoint(ev.position, ev.loudness)
			} else if hasLevel {
				// Step transition: hold old level, then jump to new.
				transitionStart := ev.position - stepTransitionBlicks
				if transitionStart < 0 {
					transitionStart = 0
				}
				addPoint(transitionStart, currentLevel)
				addPoint(ev.position, ev.loudness)
			} else {
				// First dynamic marking — set initial level.
				addPoint(ev.position, ev.loudness)
			}

			currentLevel = ev.loudness
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
							targetLevel = currentLevel + 6
						} else {
							targetLevel = currentLevel - 6
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
					targetLevel = currentLevel + 6
				} else {
					targetLevel = currentLevel - 6
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

// dynamicsToLoudness maps a MusicXML dynamics element to a dB loudness value.
// It uses InnerXML to detect which child elements are present, since empty
// elements like <p/> parse to empty strings and can't be detected via field values.
func dynamicsToLoudness(d *musicxml.Dynamics) (float64, bool) {
	xml := d.InnerXML

	// Check from most specific to least specific to avoid prefix matching issues.
	// E.g. check <pppppp before <ppppp before <pppp etc.
	switch {
	case strings.Contains(xml, "<ffffff"):
		return 15, true
	case strings.Contains(xml, "<fffff"):
		return 12, true
	case strings.Contains(xml, "<ffff"):
		return 9, true
	case strings.Contains(xml, "<fff"):
		return 6, true
	case strings.Contains(xml, "<ff"):
		return 3, true

	case strings.Contains(xml, "<pppppp"):
		return -36, true
	case strings.Contains(xml, "<ppppp"):
		return -33, true
	case strings.Contains(xml, "<pppp"):
		return -30, true
	case strings.Contains(xml, "<ppp"):
		return -27, true
	case strings.Contains(xml, "<pp"):
		return -24, true

	case strings.Contains(xml, "<sffz"):
		return 3, true
	case strings.Contains(xml, "<sfzp"):
		return 0, true
	case strings.Contains(xml, "<sfpp"):
		return 0, true
	case strings.Contains(xml, "<sfz"):
		return 3, true
	case strings.Contains(xml, "<sfp"):
		return 0, true
	case strings.Contains(xml, "<sf"):
		return 3, true

	case strings.Contains(xml, "<mp"):
		return -12, true
	case strings.Contains(xml, "<mf"):
		return -6, true

	case strings.Contains(xml, "<fp"):
		return 0, true
	case strings.Contains(xml, "<fz"):
		return 3, true
	case strings.Contains(xml, "<f"):
		return 0, true

	case strings.Contains(xml, "<rfz"):
		return 0, true
	case strings.Contains(xml, "<rf"):
		return 0, true

	case strings.Contains(xml, "<pf"):
		return -6, true
	case strings.Contains(xml, "<p"):
		return -18, true

	case strings.Contains(xml, "<n"):
		return -48, true
	}
	return 0, false
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
