package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var inputFile2 string = "input/day02.txt"

var limits = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

func getGameN(s string) (int, string) {
	RE_GAME := regexp.MustCompile(`Game (\d+): `)
	m := RE_GAME.FindStringSubmatchIndex(s)
	if m == nil {
		log.Fatal("No game")
	}
	gameN, err := strconv.Atoi(string(s)[m[2]:m[3]])
	if err != nil {
		log.Fatal(err)
	}
	return gameN, s[m[1]:]
}

func isDrawPossible(s string) bool {
	RE_NCOLOR := regexp.MustCompile(`(\d+) (.*)`)
	for _, nColor := range strings.Split(s, ",") {
		m := RE_NCOLOR.FindStringSubmatch(nColor)
		if m == nil {
			log.Fatal("Err format")
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			log.Fatal("Err format int")
		}
		color := m[2]
		limit, ok := limits[color]
		if ok {
			if n > limit {
				return false
			}
		}
	}
	return true
}

func isGamePossible(s string) bool {
	for _, draw := range strings.Split(s, ";") {
		if !isDrawPossible(draw) {
			return false
		}
	}
	return true
}

func getDrawRes(s string) map[string]int {
	var rv = map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}
	RE_NCOLOR := regexp.MustCompile(`(\d+) (.*)`)
	for _, nColor := range strings.Split(s, ",") {
		m := RE_NCOLOR.FindStringSubmatch(nColor)
		if m == nil {
			log.Fatal("Err format")
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			log.Fatal("Err format int")
		}
		color := m[2]
		rv[color] = n
	}
	return rv

}

func gamePower(game string) int {
	var mins = map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}
	for _, draw := range strings.Split(game, ";") {
		drawRes := getDrawRes(draw)
		for k, v := range drawRes {
			if v > mins[k] {
				mins[k] = v
			}
		}
	}
	rv := mins["red"] * mins["green"] * mins["blue"]
	return rv
}

func Day02b() int {
	f, err := os.OpenFile(inputFile2, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := 0
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		_, restInput := getGameN(lines.Text())
		r += gamePower(restInput)
	}
	return r

}

func Day02a() int {
	f, err := os.OpenFile(inputFile2, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := 0
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		gameN, restInput := getGameN(lines.Text())
		if isGamePossible(restInput) {
			r += gameN
		}
	}
	return r
}
