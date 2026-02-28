package main

import (
	"strconv"
	"strings"

	"github.com/egonelbre/synthv2-scripts/musicxml-to-svp/internal/musicxml"
)

type playedMeasure struct {
	measureIdx int // index into part.Measure
	verse      int // 1-based lyric number for this pass
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

func intSliceContains(s []int, v int) bool {
	for _, n := range s {
		if n == v {
			return true
		}
	}
	return false
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
			if ei.inEnding && !intSliceContains(ei.numbers, pass) {
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

type measureInfo struct {
	startBlicks int64
	divisions   int
}

// buildStructure unrolls repeats and computes measure start positions, meters, and tempos
// from the first part.
func buildStructure(firstPart *musicxml.Part) ([]playedMeasure, []measureInfo, []MeterChange, []TempoChange) {
	unrolled := unrollMeasures(firstPart.Measure)

	var meters []MeterChange
	var tempos []TempoChange
	var infos []measureInfo

	divisions := 4
	cursor := int64(0)

	for measureIdx, pm := range unrolled {
		measure := firstPart.Measure[pm.measureIdx]

		infos = append(infos, measureInfo{
			startBlicks: cursor,
			divisions:   divisions,
		})

		// Track the maximum cursor extent within this measure to correctly
		// handle pickup measures, incomplete final measures, and multi-voice parts.
		measureCursor := int64(0)
		maxCursor := int64(0)
		for _, el := range measure.Element {
			switch value := el.Value.(type) {
			case *musicxml.Attributes:
				if value.Divisions != 0 {
					divisions = value.Divisions
					infos[len(infos)-1].divisions = divisions
				}
				for _, t := range value.Time {
					meters = append(meters, MeterChange{
						MeasureIndex: measureIdx,
						Numerator:    parseBeats(t.Beats),
						Denominator:  t.BeatType,
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

		meterNum := 4
		meterDen := 4
		for _, m := range meters {
			if m.MeasureIndex <= measureIdx {
				meterNum = m.Numerator
				meterDen = m.Denominator
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

	return unrolled, infos, meters, tempos
}
