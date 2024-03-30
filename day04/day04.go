package day04

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

func Day04a(f *os.File) int {
	reLine := regexp.MustCompile(`Card\s+(\d+)\:([\s\d]*)\|([\s\d]*)`)
	reNum := regexp.MustCompile(`\d+`)
	rv := 0
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		line := lines.Text()
		m := reLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		winning := make(map[string](struct{}))
		rowWin := 0
		for _, n := range reNum.FindAllString(m[2], -1) {
			winning[n] = struct{}{}
		}
		for _, v := range reNum.FindAllString(m[3], -1) {
			if _, exists := winning[v]; exists {
				if rowWin == 0 {
					rowWin = 1
				} else {
					rowWin = rowWin * 2
				}
			}
		}
		rv += rowWin
	}
	return rv
}

func Day04b(f *os.File) int {
	reLine := regexp.MustCompile(`Card\s+(\d+)\:([\s\d]*)\|([\s\d]*)`)
	reNum := regexp.MustCompile(`\d+`)
	rv := 0
	copiesOfCards := make(map[int](int))
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		line := lines.Text()
		m := reLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		cardN, err := strconv.Atoi(m[1])
		if err != nil {
			log.Fatal(err)
		}
		cardN--
		copiesOfCards[cardN] = copiesOfCards[cardN] + 1
		winning := make(map[string](struct{}))
		rowWin := 0
		for _, n := range reNum.FindAllString(m[2], -1) {
			winning[n] = struct{}{}
		}
		for _, v := range reNum.FindAllString(m[3], -1) {
			if _, exists := winning[v]; exists {
				rowWin++
			}
		}
		for i := cardN + 1; i < cardN+1+rowWin; i++ {
			copiesOfCards[i] = copiesOfCards[i] + copiesOfCards[cardN]
		}
	}
	for _, v := range copiesOfCards {
		rv += v
	}
	return rv
}
