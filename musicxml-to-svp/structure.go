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
	// Pre-scan to find repeat regions.
	type repeatRegion struct {
		start, end int
		times      int
	}
	var regions []repeatRegion
	repeatStart := 0
	for i, m := range measures {
		if measureHasForwardRepeat(m) {
			repeatStart = i
		}
		if br := measureBackwardRepeat(m); br != nil {
			times := br.Times
			if times == 0 {
				times = 2
			}
			regions = append(regions, repeatRegion{start: repeatStart, end: i, times: times})
			repeatStart = i + 1
		}
	}

	// Emit the unrolled sequence.
	var result []playedMeasure
	i := 0
	regionIdx := 0

	for i < len(measures) {
		if regionIdx < len(regions) && i == regions[regionIdx].start {
			r := regions[regionIdx]
			regionIdx++

			// Find endings and track which measures are inside which ending.
			type endingInfo struct {
				inEnding bool
				numbers  []int
			}
			endingInfos := make([]endingInfo, r.end-r.start+1)
			var currentEndingNums []int
			inEnding := false
			maxEndingNum := 0
			hasEndings := false

			for j := r.start; j <= r.end; j++ {
				if startEnding := measureEndingStart(measures[j]); startEnding != nil {
					inEnding = true
					currentEndingNums = parseEndingNumbers(startEnding.Number)
					hasEndings = true
					for _, n := range currentEndingNums {
						if n > maxEndingNum {
							maxEndingNum = n
						}
					}
				}
				if inEnding {
					endingInfos[j-r.start] = endingInfo{inEnding: true, numbers: currentEndingNums}
				}
				if measureHasEndingStop(measures[j]) {
					inEnding = false
				}
			}

			// Emit passes.
			for pass := 1; pass <= r.times; pass++ {
				for j := r.start; j <= r.end; j++ {
					ei := endingInfos[j-r.start]
					if ei.inEnding && !intSliceContains(ei.numbers, pass) {
						continue
					}
					result = append(result, playedMeasure{measureIdx: j, verse: pass})
				}
			}

			i = r.end + 1

			// If there are endings and the final pass exceeds all ending numbers,
			// consume continuation measures (implicit final volta) with the final verse.
			if hasEndings && r.times > maxEndingNum {
				finalVerse := r.times
				for i < len(measures) {
					if regionIdx < len(regions) && i == regions[regionIdx].start {
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

		measureDuration := int64(0)
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
						Numerator:    t.Beats,
						Denominator:  t.BeatType,
					})
				}
			case *musicxml.Direction:
				if value.Sound != nil && value.Sound.Tempo != "" {
					bpm, err := strconv.ParseFloat(value.Sound.Tempo, 64)
					if err == nil {
						tempos = append(tempos, TempoChange{
							Position: cursor + measureDuration,
							BPM:      bpm,
						})
					}
				}
				if len(tempos) == 0 || value.Sound == nil {
					for _, dt := range value.DirectionType {
						if dt.Metronome != nil && dt.Metronome.PerMinute != nil {
							bpm, err := strconv.ParseFloat(dt.Metronome.PerMinute.EnclosedText, 64)
							if err == nil {
								q := beatUnitToQuarters(dt.Metronome.BeatUnit, dt.Metronome.BeatUnitDot != "")
								bpm = bpm * q
								tempos = append(tempos, TempoChange{
									Position: cursor + measureDuration,
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
						measureDuration += blicks
					}
				}
			case *musicxml.Backup:
				dur := parseDuration(value.Duration)
				if dur > 0 {
					measureDuration -= durationToBlicks(dur, divisions)
				}
			case *musicxml.Forward:
				dur := parseDuration(value.Duration)
				if dur > 0 {
					measureDuration += durationToBlicks(dur, divisions)
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
		cursor += expectedDuration
	}

	return unrolled, infos, meters, tempos
}
