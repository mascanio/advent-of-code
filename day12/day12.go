package day12

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
)

type spring byte

func (s spring) String() string {
	return string(s)
}

var re_springs = regexp.MustCompile(`([\?\.#]+)`)
var re_num = regexp.MustCompile(`\d+`)

func parseLine(line string) ([]spring, []int) {
	springs := make([]spring, 0, len(line))
	prev := byte(' ')
	for _, b := range re_springs.Find([]byte(line)) {
		if b == '.' && prev == '.' {
			continue
		}
		springs = append(springs, spring(b))
		prev = b
	}

	correct := make([]int, 0)
	for _, v := range re_num.FindAllString(line, -1) {
		n, _ := strconv.Atoi(v)
		correct = append(correct, n)
	}
	return springs, correct
}

func remainingAreNotBroken(springs []spring) bool {
	for _, v := range springs {
		if v == '#' {
			return false
		}
	}
	return true
}

type sol [32]byte

func calculateSolId(springs []spring, correct []int, prev spring) sol {
	var buffer bytes.Buffer
	for _, spring := range springs {
		buffer.WriteByte(byte(spring))
	}
	buffer.WriteByte(byte(prev))
	for _, v := range correct {
		buffer.Write([]byte(strconv.Itoa(v)))
	}
	return sha256.Sum256(buffer.Bytes())
}

func proc(springs []spring, prev spring, correct []int, sols *map[sol](int)) int {
	//fmt.Println(springs, correct)
	solId := calculateSolId(springs, correct, prev)
	if sol, ok := (*sols)[solId]; ok {
		return sol
	}
	if len(springs) == 0 && (len(correct) == 0 || len(correct) == 1 && correct[0] == 0) {
		(*sols)[solId] = 1
		return 1
	} else if len(springs) != 0 && len(correct) == 0 && remainingAreNotBroken(springs) {
		(*sols)[solId] = 1
		return 1
	} else if len(springs) == 0 || len(correct) == 0 {
		(*sols)[solId] = 0
		return 0
	}
	switch springs[0] {
	case '.':
		if prev == '#' && correct[0] != 0 {
			(*sols)[solId] = 0
			return 0
		}
		if correct[0] == 0 {
			correct = correct[1:]
		}
		rv := proc(springs[1:], '.', correct, sols)
		(*sols)[solId] = rv
		return rv
	case '#':
		if correct[0] == 0 {
			(*sols)[solId] = 0
			return 0
		}
		rv := proc(springs[1:], '#', append([]int{correct[0] - 1}, correct[1:]...), sols)
		(*sols)[solId] = rv
		return rv
	case '?':
		// Try fixed
		springs[0] = '.'
		fixedTry := proc(springs, prev, correct, sols)
		// Try broken
		springs[0] = '#'
		brokenTry := proc(springs, prev, correct, sols)
		springs[0] = '?'
		rv := fixedTry + brokenTry
		(*sols)[solId] = rv
		return rv
	}
	panic(springs[0])
}

func Day12a(f *os.File) int {
	rv := 0

	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	for lines.Scan() {
		springs, correct := parseLine(lines.Text())
		sols := make(map[sol](int), 0)
		rv += proc(springs, ' ', correct, &sols)
	}

	return rv
}

func Day12b(f *os.File) int {
	rv := 0

	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	i := 0
	c := make(chan int)
	go func() {
		wg := sync.WaitGroup{}
		defer close(c)
		for lines.Scan() {
			wg.Add(1)
			line := lines.Text()
			go func(wg *sync.WaitGroup, line string) {
				springs, correct := parseLine(line)
				newSprings := make([]spring, 0, len(springs)*5)
				for i := 0; i < 5; i++ {
					newSprings = append(newSprings, springs...)
					if i != 4 {
						newSprings = append(newSprings, '?')
					}
				}
				newCorrect := make([]int, 0, len(correct)*5)
				for i := 0; i < 5; i++ {
					newCorrect = append(newCorrect, correct...)
				}
				sols := make(map[sol](int), 0)
				c <- (proc(newSprings, ' ', newCorrect, &sols))
				wg.Done()
			}(&wg, line)
			i++
		}
		wg.Wait()
	}()

	finished := 0
	for v := range c {
		finished++
		fmt.Println(finished, "/", i)
		rv += v
	}

	return rv
}
