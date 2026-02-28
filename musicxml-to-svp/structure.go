package main

import (
	"slices"
	"strconv"
	"strings"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

type playedMeasure struct {
	measureIdx  int   // index into part.Measure
	verse       int   // 1-based lyric number for this pass
	startBlicks int64 // absolute position in blicks
	divisions   int   // time units per quarter note
}

func measureHasForwardRepeat(m *musicxml.Measure) bool {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Repeat != nil && bl.Repeat.Direction == "forward" {
				return true
			}
		}
	}
	return false
}

func measureBackwardRepeat(m *musicxml.Measure) *musicxml.Repeat {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Repeat != nil && bl.Repeat.Direction == "backward" {
				return bl.Repeat
			}
		}
	}
	return nil
}

func measureEndingStart(m *musicxml.Measure) *musicxml.Ending {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Ending != nil && bl.Ending.Type == "start" {
				return bl.Ending
			}
		}
	}
	return nil
}

func measureHasEndingStop(m *musicxml.Measure) bool {
	for _, el := range m.Element {
		if bl, ok := el.Value.(*musicxml.Barline); ok {
			if bl.Ending != nil && (bl.Ending.Type == "stop" || bl.Ending.Type == "discontinue") {
				return true
			}
		}
	}
	return false
}

func parseEndingNumbers(s string) []int {
	var nums []int
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if n, err := strconv.Atoi(part); err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

func unrollMeasures(measures []*musicxml.Measure) []playedMeasure {
	type endingInfo struct {
		inEnding bool
		numbers  []int
	}
	type repeatRegion struct {
		start, end   int
		times        int
		hasEndings   bool
		maxEndingNum int
		endingInfos  []endingInfo
	}

	// Pre-scan using a stack to properly handle nested repeats.
	var regions []repeatRegion
	var stack []int

	for i, m := range measures {
		if measureHasForwardRepeat(m) {
			stack = append(stack, i)
		}
		if br := measureBackwardRepeat(m); br != nil {
			times := br.Times
			if times == 0 {
				times = 2
			}
			start := 0
			if len(stack) > 0 {
				start = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

			r := repeatRegion{
				start:       start,
				end:         i,
				times:       times,
				endingInfos: make([]endingInfo, i-start+1),
			}

			var currentEndingNums []int
			inEnding := false
			for j := start; j <= i; j++ {
				if se := measureEndingStart(measures[j]); se != nil {
					inEnding = true
					currentEndingNums = parseEndingNumbers(se.Number)
					r.hasEndings = true
					for _, n := range currentEndingNums {
						if n > r.maxEndingNum {
							r.maxEndingNum = n
						}
					}
				}
				if inEnding {
					r.endingInfos[j-start] = endingInfo{inEnding: true, numbers: currentEndingNums}
				}
				if measureHasEndingStop(measures[j]) {
					inEnding = false
				}
			}

			regions = append(regions, r)
		}
	}

	// findOuterAt returns the outermost (largest span) region starting at pos.
	findOuterAt := func(pos int) *repeatRegion {
		var best *repeatRegion
		for i := range regions {
			if regions[i].start == pos {
				if best == nil || regions[i].end > best.end {
					best = &regions[i]
				}
			}
		}
		return best
	}

	// findNestedAt returns the outermost nested region starting at pos
	// that fits strictly within parent.
	findNestedAt := func(pos int, parent *repeatRegion) *repeatRegion {
		var best *repeatRegion
		for i := range regions {
			r := &regions[i]
			if r.start == pos && r.end <= parent.end &&
				!(r.start == parent.start && r.end == parent.end) {
				if best == nil || r.end > best.end {
					best = r
				}
			}
		}
		return best
	}

	// emitPass recursively emits measures for one pass of a region,
	// handling nested repeats.
	var emitPass func(r *repeatRegion, pass int, result *[]playedMeasure)
	emitPass = func(r *repeatRegion, pass int, result *[]playedMeasure) {
		j := r.start
		for j <= r.end {
			ei := r.endingInfos[j-r.start]
			if ei.inEnding && !slices.Contains(ei.numbers, pass) {
				j++
				continue
			}
			if nested := findNestedAt(j, r); nested != nil {
				for p := 1; p <= nested.times; p++ {
					emitPass(nested, p, result)
				}
				j = nested.end + 1
				continue
			}
			*result = append(*result, playedMeasure{measureIdx: j, verse: pass})
			j++
		}
	}

	// Top-level emit.
	var result []playedMeasure
	i := 0
	for i < len(measures) {
		if r := findOuterAt(i); r != nil {
			for pass := 1; pass <= r.times; pass++ {
				emitPass(r, pass, &result)
			}
			i = r.end + 1

			// If there are endings and the final pass exceeds all ending numbers,
			// consume continuation measures (implicit final volta) with the final verse.
			if r.hasEndings && r.times > r.maxEndingNum {
				finalVerse := r.times
				for i < len(measures) {
					if findOuterAt(i) != nil {
						break
					}
					result = append(result, playedMeasure{measureIdx: i, verse: finalVerse})
					i++
				}
			}
			continue
		}

		result = append(result, playedMeasure{measureIdx: i, verse: 1})
		i++
	}

	return result
}

// navigationMarkers holds positions of segno, coda, fine, and jump instructions
// found in a sequence of measures.
type navigationMarkers struct {
	// segnoIdx maps segno ID -> measure index
	segnoIdx map[string]int
	// codaIdx maps coda ID -> measure index
	codaIdx map[string]int
	// fineIdx is the measure index of the Fine marker, or -1
	fineIdx int
	// jumpIdx is the measure index where a D.C. or D.S. jump occurs, or -1
	jumpIdx int
	// jumpType is "dacapo" or "dalsegno"
	jumpType string
	// jumpTarget is the segno ID for D.S. jumps (empty for D.C.)
	jumpTarget string
	// tocodaIdx is the measure index of the To Coda instruction, or -1
	tocodaIdx int
	// tocodaTarget is the coda ID referenced by the To Coda instruction
	tocodaTarget string
}

func scanNavigation(measures []*musicxml.Measure) navigationMarkers {
	nav := navigationMarkers{
		segnoIdx: make(map[string]int),
		codaIdx:  make(map[string]int),
		fineIdx:  -1,
		jumpIdx:  -1,
		tocodaIdx: -1,
	}

	for i, m := range measures {
		for _, el := range m.Element {
			switch v := el.Value.(type) {
			case *musicxml.Barline:
				if v.Segno != "" {
					nav.segnoIdx[v.Segno] = i
				}
				if v.Coda != "" {
					nav.codaIdx[v.Coda] = i
				}
			case *musicxml.Direction:
				if v.Sound == nil {
					continue
				}
				s := v.Sound
				if s.Segno != "" {
					nav.segnoIdx[s.Segno] = i
				}
				if s.Coda != "" {
					nav.codaIdx[s.Coda] = i
				}
				if s.Fine != "" {
					nav.fineIdx = i
				}
				if s.Dacapo == "yes" {
					nav.jumpIdx = i
					nav.jumpType = "dacapo"
				}
				if s.Dalsegno != "" {
					nav.jumpIdx = i
					nav.jumpType = "dalsegno"
					nav.jumpTarget = s.Dalsegno
				}
				if s.Tocoda != "" {
					nav.tocodaIdx = i
					nav.tocodaTarget = s.Tocoda
				}
			}
		}
	}
	return nav
}

// measureHasAfterJumpRepeat checks if a measure has a backward repeat
// with after-jump="yes".
func measureHasAfterJumpRepeat(m *musicxml.Measure) bool {
	if br := measureBackwardRepeat(m); br != nil {
		return br.AfterJump == "yes"
	}
	return false
}

// unrollMeasureRange unrolls measures[start:end+1] using the repeat logic.
// When afterJump is true, repeats and voltas are skipped unless the backward
// repeat has after-jump="yes".
func unrollMeasureRange(measures []*musicxml.Measure, start, end int, afterJump bool) []playedMeasure {
	sub := measures[start : end+1]
	if !afterJump {
		result := unrollMeasures(sub)
		// Remap indices back to the original measure array.
		for i := range result {
			result[i].measureIdx += start
		}
		return result
	}

	// After a jump: play each measure once, except for repeats with after-jump="yes".
	var result []playedMeasure
	i := 0
	for i < len(sub) {
		// Check for after-jump repeat regions.
		if measureHasForwardRepeat(sub[i]) {
			// Find the matching backward repeat with after-jump="yes".
			regionStart := i
			foundAfterJump := false
			for j := i; j < len(sub); j++ {
				if measureHasAfterJumpRepeat(sub[j]) {
					// Unroll this region normally.
					regionMeasures := sub[regionStart : j+1]
					regionResult := unrollMeasures(regionMeasures)
					for k := range regionResult {
						regionResult[k].measureIdx += start + regionStart
					}
					result = append(result, regionResult...)
					i = j + 1
					foundAfterJump = true
					break
				}
			}
			if foundAfterJump {
				continue
			}
		}
		result = append(result, playedMeasure{measureIdx: start + i, verse: 1})
		i++
	}
	return result
}

// unrollWithNavigation handles D.C., D.S., Coda, and Fine navigation on top
// of barline repeats. If no navigation markers are found, it falls back to
// plain unrollMeasures.
func unrollWithNavigation(measures []*musicxml.Measure) []playedMeasure {
	nav := scanNavigation(measures)

	if nav.jumpIdx < 0 {
		return unrollMeasures(measures)
	}

	last := len(measures) - 1

	// Determine where the jump goes back to.
	jumpBackTo := 0
	if nav.jumpType == "dalsegno" {
		if idx, ok := nav.segnoIdx[nav.jumpTarget]; ok {
			jumpBackTo = idx
		}
	}

	// Phase 1: play from beginning through the jump measure with normal repeats.
	result := unrollMeasureRange(measures, 0, nav.jumpIdx, false)

	// Determine the end of phase 2.
	phase2End := last
	if nav.tocodaIdx >= 0 && nav.tocodaIdx >= jumpBackTo {
		phase2End = nav.tocodaIdx
	} else if nav.fineIdx >= 0 && nav.fineIdx >= jumpBackTo {
		phase2End = nav.fineIdx
	}

	// Phase 2: jump back, play without repeats (unless after-jump).
	result = append(result, unrollMeasureRange(measures, jumpBackTo, phase2End, true)...)

	// Phase 3: if there's a coda, jump to coda and play to end with normal repeats.
	if nav.tocodaIdx >= 0 {
		codaIdx := -1
		if idx, ok := nav.codaIdx[nav.tocodaTarget]; ok {
			codaIdx = idx
		}
		if codaIdx >= 0 {
			result = append(result, unrollMeasureRange(measures, codaIdx, last, false)...)
		}
	}

	return result
}

// buildStructure unrolls repeats and computes measure start positions, meters, and tempos
// from the first part.
func buildStructure(firstPart *musicxml.Part) ([]playedMeasure, []MeterChange, []TempoChange) {
	unrolled := unrollWithNavigation(firstPart.Measure)

	var meters []MeterChange
	var tempos []TempoChange

	divisions := 4
	cursor := int64(0)
	meterNum := 4
	meterDen := 4

	for measureIdx := range unrolled {
		pm := &unrolled[measureIdx]
		measure := firstPart.Measure[pm.measureIdx]

		pm.startBlicks = cursor
		pm.divisions = divisions

		// Track the maximum cursor extent within this measure to correctly
		// handle pickup measures, incomplete final measures, and multi-voice parts.
		measureCursor := int64(0)
		maxCursor := int64(0)
		for _, el := range measure.Element {
			switch value := el.Value.(type) {
			case *musicxml.Attributes:
				if value.Divisions != 0 {
					divisions = value.Divisions
					pm.divisions = divisions
				}
				for _, t := range value.Time {
					meterNum = parseBeats(t.Beats)
					meterDen = t.BeatType
					meters = append(meters, MeterChange{
						MeasureIndex: measureIdx,
						Numerator:    meterNum,
						Denominator:  meterDen,
					})
				}
			case *musicxml.Direction:
				addedTempo := false
				if value.Sound != nil && value.Sound.Tempo != "" {
					bpm, err := strconv.ParseFloat(value.Sound.Tempo, 64)
					if err == nil {
						tempos = append(tempos, TempoChange{
							Position: cursor + measureCursor,
							BPM:      bpm,
						})
						addedTempo = true
					}
				}
				if !addedTempo {
					for _, dt := range value.DirectionType {
						if dt.Metronome != nil && dt.Metronome.PerMinute != nil {
							bpm, err := strconv.ParseFloat(dt.Metronome.PerMinute.EnclosedText, 64)
							if err == nil {
								q := beatUnitToQuarters(dt.Metronome.BeatUnit, dt.Metronome.BeatUnitDot != "")
								bpm = bpm * q
								tempos = append(tempos, TempoChange{
									Position: cursor + measureCursor,
									BPM:      bpm,
								})
							}
						}
					}
				}
			case *musicxml.Note:
				if value.Grace != nil {
					continue
				}
				dur := parseDuration(value.Duration)
				if dur > 0 {
					if value.Chord == "" {
						blicks := durationToBlicks(dur, divisions)
						measureCursor += blicks
						if measureCursor > maxCursor {
							maxCursor = measureCursor
						}
					}
				}
			case *musicxml.Backup:
				dur := parseDuration(value.Duration)
				if dur > 0 {
					measureCursor -= durationToBlicks(dur, divisions)
					if measureCursor < 0 {
						measureCursor = 0
					}
				}
			case *musicxml.Forward:
				dur := parseDuration(value.Duration)
				if dur > 0 {
					measureCursor += durationToBlicks(dur, divisions)
					if measureCursor > maxCursor {
						maxCursor = measureCursor
					}
				}
			}
		}

		expectedDuration := int64(meterNum) * blicksPerQuarter * 4 / int64(meterDen)

		// Use actual measure duration when available (handles pickup measures,
		// incomplete final measures, and cadenzas). Fall back to expected
		// duration for empty measures.
		if maxCursor > 0 && maxCursor != expectedDuration {
			cursor += maxCursor
		} else {
			cursor += expectedDuration
		}
	}

	return unrolled, meters, tempos
}
