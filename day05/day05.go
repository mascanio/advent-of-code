package day05

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type mapLine struct {
	destRangeStart   int
	sourceRangeStart int
	rangeLen         int
}

type mapping struct {
	lines []mapLine
}

type almanach struct {
	seeds    [][]int
	mappings []mapping
}

func parseAlmanach(f *os.File) almanach {
	rv := almanach{}
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)

	rv.seeds = parseSeeds(lines)
	rv.mappings = parseMappings(lines)
	return rv
}

func parseMappings(lines *bufio.Scanner) []mapping {
	var rv []mapping
	reMap := regexp.MustCompile(`(.*) map\:`)
	for lines.Scan() {
		line := lines.Text()
		if lines.Text() == "" {
			continue
		}
		if m := reMap.FindStringSubmatch(line); m != nil {
			rv = append(rv, mapping{})
			continue
		}
		addLine(&rv[len(rv)-1], parseMapLine(line))
	}
	return rv
}

func parseSeeds(lines *bufio.Scanner) [][]int {
	lines.Scan()
	var seeds [][]int
	reSeeds := regexp.MustCompile(`\d+ \d+`)
	for _, seed := range reSeeds.FindAllString(lines.Text(), -1) {
		ss := strings.Split(seed, " ")
		n, _ := strconv.Atoi(ss[0])
		m, _ := strconv.Atoi(ss[1])
		seeds = append(seeds, []int{n, m + n})
	}
	return seeds
}

func addLine(m *mapping, line mapLine) {
	m.lines = append(m.lines, line)
}

func parseMapLine(s string) mapLine {
	reLine := regexp.MustCompile(`(\d+) (\d+) (\d+)`)
	m := reLine.FindStringSubmatch(s)
	destRangeStart, _ := strconv.Atoi(m[1])
	sourceRangeStart, _ := strconv.Atoi(m[2])
	rangeLen, _ := strconv.Atoi(m[3])

	return mapLine{destRangeStart: destRangeStart, sourceRangeStart: sourceRangeStart, rangeLen: rangeLen}
}

func intersecs(start, end int, line *mapLine) bool {
	return end > line.sourceRangeStart && start < line.sourceRangeStart+line.rangeLen
}

func intersections(start, end int, mapping *mapping) (rv [][]int) {
	startRem := []int{start, end}
	endRem := []int{start, end}
	for _, v := range mapping.lines {
		if !intersecs(start, end, &v) {
			continue
		}
		offset := v.destRangeStart - v.sourceRangeStart
		intersStart := max(v.sourceRangeStart, start)
		intersEnd := min(v.sourceRangeStart+v.rangeLen, end)
		rStart := offset + intersStart
		rEnd := offset + intersEnd
		rv = append(rv, []int{rStart, rEnd})
		if intersStart < startRem[1] {
			startRem[1] = intersStart
		}
		if intersEnd > endRem[0] {
			endRem[0] = intersEnd
		}
	}
	// Remaining
	if startRem[0] != startRem[1] {
		rv = append(rv, startRem)
	}
	if endRem[0] != endRem[1] {
		rv = append(rv, endRem)
	}
	return
}

func foo(input [][]int, mappings []mapping) int {
	rv := math.MaxInt
	if len(mappings) == 0 {
		for _, v := range input {
			rv = min(rv, v[0])
		}
		return rv
	}
	for _, v := range input {
		start := v[0]
		end := v[1]
		mapping := &mappings[0]
		inters := intersections(start, end, mapping)
		step := foo(inters, mappings[1:])
		rv = min(rv, step)
	}
	return rv
}

func Day05b(f *os.File) int {
	almanach := parseAlmanach(f)
	return foo(almanach.seeds, almanach.mappings)
}
