package day13

import (
	"bufio"
	"os"
	"slices"
	"strings"
	"sync"
)

type node bool
type line []node

type pattern struct {
	n            []line
	nCols, nRows int
}

func (n node) String() string {
	if n {
		return "#"
	}
	return "."
}

func (p pattern) String() string {
	sb := strings.Builder{}

	for _, row := range p.n {
		for _, v := range row {
			sb.WriteString(v.String())
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func (l line) String() string {
	sb := strings.Builder{}

	for _, e := range l {
		sb.WriteString(e.String())
	}

	return sb.String()
}

func parseNode(b byte) node {
	return node(b == '#')
}

func (p *pattern) getRow(y int) line {
	return p.n[y]
}

func (p *pattern) getCol(x int) line {
	rv := make([]node, 0, p.nRows)

	for _, row := range p.n {
		rv = append(rv, row[x])
	}
	return rv
}

func (a *line) eq(b line, diff int) bool {
	for i := range len(*a) {
		if (*a)[i] != (b)[i] {
			if diff == 0 {
				return false
			}
			diff--
		}
	}
	return true
}

func parsePatterns(f *os.File) []pattern {
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)
	rrv := make([]pattern, 0)
	rv := pattern{}
	rv.n = make([]line, 0)
	for lines.Scan() {
		line := lines.Text()
		if line == "" {
			rrv = append(rrv, rv)
			rv = pattern{}
			continue
		}
		newRow := make([]node, 0, rv.nCols)
		for _, b := range []byte(line) {
			if rv.nRows == 0 {
				rv.nCols++
			}
			newRow = append(newRow, parseNode(b))
		}
		newRow = slices.Clip(newRow)
		rv.n = append(rv.n, newRow)
		rv.nRows++
	}
	rrv = append(rrv, rv)
	return rrv
}

func rMirror(a, b []line, tolerance, diff int) bool {
	if len(a) == 0 || len(b) == 0 {
		return tolerance == 0
	}
	if !a[0].eq(b[0], 0) {
		if tolerance > 0 && a[0].eq(b[0], diff) {
			tolerance--
		} else {
			return false
		}
	}
	return rMirror(a[1:], b[1:], tolerance, diff)
}

func isMirrorRow(p *pattern, y, tolerance, diff int) bool {
	a := make([]line, 0, p.nRows)
	b := make([]line, 0, p.nRows)

	for i := 0; i <= y; i++ {
		a = append(a, p.getRow(i))
	}
	for i := y + 1; i < p.nRows; i++ {
		b = append(b, p.getRow(i))
	}
	slices.Reverse(a)

	return rMirror(a, b, tolerance, diff)
}

func isMirrorCol(p *pattern, y, tolerance, diff int) bool {
	a := make([]line, 0, p.nCols)
	b := make([]line, 0, p.nCols)

	for i := 0; i <= y; i++ {
		a = append(a, p.getCol(i))
	}
	for i := y + 1; i < p.nCols; i++ {
		b = append(b, p.getCol(i))
	}
	slices.Reverse(a)

	return rMirror(a, b, tolerance, diff)
}

func A(f *os.File) int {
	rv := 0
	for _, pattern := range parsePatterns(f) {
		for row := 0; row < pattern.nRows-1; row++ {
			if isMirrorRow(&pattern, row, 0, 0) {
				rv += (row + 1) * 100
				break
			}
		}
		for col := 0; col < pattern.nCols-1; col++ {
			if isMirrorCol(&pattern, col, 0, 0) {
				rv += col + 1
				break
			}
		}
	}
	return rv
}

func B(f *os.File) int {
	rv := 0
	c := make(chan int)

	go func() {
		wg := sync.WaitGroup{}
		defer close(c)
		for _, p := range parsePatterns(f) {
			wg.Add(1)
			go func(pattern pattern) {
				rv := 0
				for row := 0; row < pattern.nRows-1; row++ {
					if isMirrorRow(&pattern, row, 1, 1) {
						rv += (row + 1) * 100
						break
					}
				}
				for col := 0; col < pattern.nCols-1; col++ {
					if isMirrorCol(&pattern, col, 1, 1) {
						rv += col + 1
						break
					}
				}
				c <- rv
				wg.Done()
			}(p)
		}
		wg.Wait()
	}()
	for r := range c {
		rv += r
	}
	return rv
}
