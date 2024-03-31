package day06

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type race struct {
	time     int
	distance int
}

func parseRaces(f *os.File, partb bool) []race {
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)
	lines.Scan()
	reNum := regexp.MustCompile(`\d+`)
	var times []int
	line := lines.Text()
	if partb {
		line = strings.ReplaceAll(line, " ", "")
	}
	for _, v := range reNum.FindAllString(line, -1) {
		n, _ := strconv.Atoi(v)
		times = append(times, n)
	}
	var distances []int
	lines.Scan()
	line = lines.Text()
	if partb {
		line = strings.ReplaceAll(line, " ", "")
	}
	for _, v := range reNum.FindAllString(line, -1) {
		n, _ := strconv.Atoi(v)
		distances = append(distances, n)
	}
	rv := make([]race, 0, len(times))
	for i := range times {
		rv = append(rv, race{time: times[i], distance: distances[i]})
	}
	return rv
}

func solveRace(race race) []int {
	rv := make([]int, 0, race.time-1)
	for i := 0; i < race.time; i++ {
		remTime := race.time - i
		if remTime*i > race.distance {
			rv = append(rv, i)
		}
	}
	return rv
}

func Day06a(f *os.File) int {
	rv := 1
	races := parseRaces(f, false)
	for _, v := range races {
		rv *= len(solveRace(v))
	}
	return rv
}

func Day06b(f *os.File) int {
	rv := 1
	races := parseRaces(f, true)
	for _, v := range races {
		rv *= len(solveRace(v))
	}
	return rv
}
