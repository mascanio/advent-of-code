package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
)

var inputFile string = "input/day01.txt"

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func getDigitsA(reader []byte) (int, int) {
	firstDigit, lastDigit := 0, 0
	for _, b := range reader {
		if isDigit(b) {
			if firstDigit == 0 {
				firstDigit = int(b - '0')
			} else {
				lastDigit = int(b - '0')
			}
		}
	}
	return firstDigit, lastDigit
}

var numberStr = []string{
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

func getDigitsB(bs []byte) (int, int) {
	firstDigit, lastDigit := 0, 0
	scanner := bufio.NewScanner(bytes.NewReader(bs))
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		if err == nil && token != nil {
			for i, b := range token {
				if isDigit(b) {
					advance = i + 1
					token = make([]byte, 1)
					token[0] = b
					return
				}
				str := string(token[i:])
				for n, pre := range numberStr {
					if strings.HasPrefix(str, pre) {
						advance = i + 1
						token = make([]byte, 1)
						token[0] = byte(n + 1 + '0')
						return
					}
				}
				advance++
			}
		}
		return
	})
	scanner.Scan()
	var err error
	firstDigit, err = strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		lastDigit, err = strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

	}
	return firstDigit, lastDigit
}

func day01(f *os.File, getDigitsFun func([]byte) (int, int)) int {
	reader := bufio.NewReader(f)

	acu := 0
	lineScanner := bufio.NewScanner(reader)
	lineScanner.Split(bufio.ScanLines)
	for lineScanner.Scan() {
		firstDigit, lastDigit := getDigitsFun(lineScanner.Bytes())
		if lastDigit == 0 {
			lastDigit = firstDigit
		}
		acu += 10*firstDigit + lastDigit
	}
	return acu
}

func Day01a() int {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	return day01(f, getDigitsA)
}

func Day01b() int {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	return day01(f, getDigitsB)
}
