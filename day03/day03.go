package day03

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

type plan struct {
	Nrows int
	Ncols int
	p     []string
}

func getPlan(path string) plan {
	var plan plan
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		plan.p = append(plan.p, lines.Text())
		plan.Nrows++
	}
	plan.Ncols = len(plan.p[0])

	return plan
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isSymbol(b byte) bool {
	return !isDigit(b) && b != '.'
}

func isSymbolAdjacent(x, y int, plan *plan) bool {
	for i := -1; i <= 1; i++ {
		row := i + x
		if row < 0 || row >= plan.Nrows {
			continue
		}
		for j := -1; j <= 1; j++ {
			col := j + y
			if (i == 0 && j == 0) || col < 0 || col >= plan.Ncols {
				continue
			}
			if isSymbol(plan.p[row][col]) {
				return true
			}
		}
	}
	return false
}

func Day03a(inputPath string) int {
	plan := getPlan(inputPath)
	rv := 0
	for i := 0; i < plan.Nrows; i++ {
		curNumber := 0
		curNumberAdj := false
		for j := 0; j < plan.Ncols; j++ {
			c := plan.p[i][j]
			if !isDigit(c) {
				if curNumberAdj {
					rv += curNumber
				}
				curNumber = 0
				curNumberAdj = false
			} else {
				curNumber = curNumber*10 + int(c-'0')
				if !curNumberAdj {
					if isSymbolAdjacent(i, j, &plan) {
						curNumberAdj = true
					}
				}
			}
		}
		if curNumberAdj {
			rv += curNumber
		}
		curNumber = 0
		curNumberAdj = false
	}
	return rv
}

var re_num = regexp.MustCompile(`\d+`)

func gearsInOtherRow(y int, s string) []int {
	found := make([]int, 0, 2)
	f := re_num.FindAllStringIndex(s, -1)
	for _, m := range f {
		if (m[0] <= y-1 && (m[1]-1) >= y+1) || m[0] == y || m[0] == y+1 || (m[1]-1) == y || (m[1]-1) == y-1 {
			if len(found) < 2 {
				n, _ := strconv.Atoi(s[m[0]:m[1]])
				found = append(found, n)
			} else {
				return nil
			}
		}
	}
	return found
}

func gearsInSameRow(y int, s string) []int {
	found := make([]int, 0, 2)
	f := re_num.FindAllStringIndex(s, -1)
	for _, m := range f {
		if (m[1]-1) == y-1 || m[0] == y+1 {
			if len(found) < 2 {
				n, _ := strconv.Atoi(s[m[0]:m[1]])
				found = append(found, n)
			} else {
				return nil
			}
		}
	}
	return found
}

func gears(x, y int, plan *plan) int {

	found := make([]int, 0, 2)

	row := x - 1
	if row >= 0 {
		pf := gearsInOtherRow(y, plan.p[row])
		if len(pf) != 0 {
			if len(pf)+len(found) <= 2 {
				found = append(found, pf...)
			} else {
				return 0
			}
		}
	}
	row = x + 1
	if row < plan.Nrows {
		pf := gearsInOtherRow(y, plan.p[row])
		if len(pf) != 0 {
			if len(pf)+len(found) <= 2 {
				found = append(found, pf...)
			} else {
				return 0
			}
		}
	}
	row = x
	pf := gearsInSameRow(y, plan.p[row])
	if len(pf) != 0 {
		if len(pf)+len(found) <= 2 {
			found = append(found, pf...)
		} else {
			return 0
		}
	}
	if len(found) != 2 {
		return 0
	}
	return found[0] * found[1]
}

func Day03b(inputPath string) int {
	plan := getPlan(inputPath)
	rv := 0
	for i := 0; i < plan.Nrows; i++ {
		for j := 0; j < plan.Ncols; j++ {
			c := plan.p[i][j]
			if c != '*' {
				continue
			}
			rv += gears(i, j, &plan)
		}
	}
	return rv
}
