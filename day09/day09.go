package day09

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func parseLine(reNums *regexp.Regexp, line string) []int {
	var nums []int
	for _, ns := range reNums.FindAllString(line, -1) {
		n, _ := strconv.Atoi(ns)
		nums = append(nums, n)
	}
	return nums
}

func prependInt(x []int, y int) []int {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}

func processLineR(line []int) []int {
	allZeroes := true
	newLine := make([]int, 0, len(line)+1)
	for i := 0; i < len(line)-1; i++ {
		x := line[i+1] - line[i]
		if x != 0 {
			allZeroes = false
		}
		newLine = append(newLine, x)
	}
	if !allZeroes {
		newLine = processLineR(newLine)
	}
	extrapolatedF := line[len(line)-1] + newLine[len(newLine)-1]
	extrapolatedB := line[0] - newLine[0]
	return prependInt(append(line, extrapolatedF), extrapolatedB)
}

func extrapolateForward(line []int) int {
	r := processLineR(line)
	return r[len(r)-1]
}

func extrapolateBack(line []int) int {
	r := processLineR(line)
	return r[0]
}

func Day09a(f *os.File) int {
	rv := 0
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	reNums := regexp.MustCompile(`\-?\d+`)
	for lines.Scan() {
		line := lines.Text()
		nums := parseLine(reNums, line)
		rv += extrapolateForward(nums)
	}

	return rv
}

func Day09b(f *os.File) int {
	defer func(t time.Time) {
		fmt.Println("SEQ", time.Since(t))
	}(time.Now())
	rv := 0
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	reNums := regexp.MustCompile(`\-?\d+`)
	for lines.Scan() {
		line := lines.Text()
		nums := parseLine(reNums, line)
		rv += extrapolateBack(nums)
	}

	return rv
}

func Day09bparallel(f *os.File) int {
	defer func(t time.Time) {
		fmt.Println("PAR", time.Since(t))
	}(time.Now())
	rv := 0
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	reNums := regexp.MustCompile(`\-?\d+`)
	c := make(chan int)
	go func() {
		wg := sync.WaitGroup{}
		defer close(c)
		for lines.Scan() {
			wg.Add(1)
			line := lines.Text()
			go func(wg *sync.WaitGroup) {
				nums := parseLine(reNums, line)
				c <- extrapolateBack(nums)
				wg.Done()
			}(&wg)
		}
		wg.Wait()
	}()

	for v := range c {
		rv += v
	}

	return rv
}
