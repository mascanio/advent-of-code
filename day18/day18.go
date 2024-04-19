package day18

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/mascanio/advent-of-code/utils"
)

type inputLine struct {
	dir   utils.Direction
	n     int
	color string
}

func parseDir(s string) utils.Direction {
	switch s {
	case "U":
		return utils.Up
	case "D":
		return utils.Down
	case "R":
		return utils.Right
	case "L":
		return utils.Left
	}
	panic("Invalid direction " + s)
}

func parseDirB(s string) utils.Direction {
	switch s {
	case "3":
		return utils.Up
	case "1":
		return utils.Down
	case "0":
		return utils.Right
	case "2":
		return utils.Left
	}
	panic("Invalid direction " + s)
}

func parseInputLine(s string) inputLine {
	reLine := regexp.MustCompile(`([UDRL]) (\d+) (.*)`)
	m := reLine.FindStringSubmatch(s)
	n, _ := strconv.Atoi(m[2])
	return inputLine{parseDir(m[1]), n, m[3]}
}

func parseInputLineB(s string) inputLine {
	reLine := regexp.MustCompile(`[UDRL] \d+ \(#([\da-f]+)(\d)\)`)
	m := reLine.FindStringSubmatch(s)
	n, _ := strconv.ParseInt(m[1], 16, 0)
	return inputLine{parseDirB(m[2]), int(n), m[0]}
}

type input struct {
	lines []inputLine
}

func parseInput(f *os.File, parseL func(s string) inputLine) input {
	rv := input{lines: make([]inputLine, 0)}
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	for lines.Scan() {
		rv.lines = append(rv.lines, parseL(lines.Text()))
	}
	return rv
}

type dig int

const (
	undef dig = iota
	outside
	digged
)

type node struct {
	n   dig
	pos utils.Pos
}

func (n node) String() string {
	if n.n == digged {
		return "#"
	}
	return "."
}

type matrix struct {
	utils.NodeMap[node]
}

func getSize(input input) (rows, cols int, initPos utils.Pos) {
	var x, y, maxx, maxy, minx, miny int

	for _, line := range input.lines {
		switch line.dir {
		case utils.Up:
			y -= line.n
			if y < miny {
				miny = y
			}
		case utils.Down:
			y += line.n
			if y > maxy {
				maxy = y
			}
		case utils.Left:
			x -= line.n
			if x < minx {
				minx = x
			}
		case utils.Right:
			x += line.n
			if x > maxx {
				maxx = x
			}
		}
	}
	maxx++
	maxy++
	rows = maxy - miny
	cols = maxx - minx
	initPos = utils.Pos{Row: -miny, Col: -minx}
	return
}

func createMatrix(input input) matrix {
	nRows, nCols, iniPos := getSize(input)
	rv := matrix{utils.CreateNodeMap[node](nRows, nCols)}

	pos := iniPos
	for _, line := range input.lines {
		newPos := pos.MoveN(line.dir, line.n+1)
		for i := pos; i != newPos; i = i.Move(line.dir) {
			rv.Set(i, node{digged, utils.Pos{}})
		}
		pos = pos.MoveN(line.dir, line.n)
	}

	for row := 0; row < nRows; row++ {
		for col := 0; col < nCols; col++ {
			rv.Nodes[row][col].pos = utils.Pos{Row: row, Col: col}
		}
	}

	return rv
}

func fillMatrix(m *matrix) int {
	q := make(utils.Queue[*node], 0)
	rv := 0
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			if i != 0 && j != 0 && i != m.Rows-1 && j != m.Cols-1 {
				continue
			}
			node := m.AtP(utils.Pos{Row: i, Col: j})
			if node.n == undef {
				node.n = outside
				q.Push(node)
			}
		}
	}
	for len(q) != 0 {
		n := q.Pop()
		for _, adjNode := range m.GetAdjacentNodes(n.pos, utils.Directions[:], func(candidate, current node, candidateP, currentP utils.Pos) bool {
			return candidate.n == undef
		}) {
			adjNode.n = outside
			q.Push(adjNode)
		}
	}
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			pos := utils.Pos{Row: i, Col: j}
			if m.At(pos).n != outside {
				m.AtP(pos).n = digged
				rv++
			}
		}
	}

	return rv
}

func getVertex(input input) (loopLen int, rv []utils.Pos) {
	rv = make([]utils.Pos, 0, len(input.lines))

	curPos := utils.Pos{Row: 0, Col: 0}
	rv = append(rv, curPos)
	for _, line := range input.lines {
		loopLen += line.n
		curPos = curPos.MoveN(line.dir, line.n)
		rv = append(rv, curPos)
	}
	return
}

func area(poss []utils.Pos) int {
	// Shoelace
	rv := 0
	for i := 0; i < len(poss)-1; i++ {
		a, b := poss[i], poss[i+1]
		tmp := a.Row*b.Col - a.Col*b.Row
		rv += tmp
	}
	if rv < 0 {
		rv = -rv
	}
	return rv / 2
}

func numInside(area, looplen int) int {
	// Use picks theorem
	return area - looplen/2 + 1
}

func A(f *os.File) int {
	input := parseInput(f, parseInputLine)
	m := createMatrix(input)
	rv := fillMatrix(&m)
	return rv
}

func B(f *os.File) int {
	input := parseInput(f, parseInputLineB)
	iniTime := time.Now()
	defer func() { fmt.Println(time.Since(iniTime)) }()
	loopLen, loop := getVertex(input)
	area := area(loop)
	return numInside(area, loopLen) + loopLen
}
