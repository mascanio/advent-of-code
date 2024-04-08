package day16

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

type pos struct {
	row, col int
}

type direction struct {
	rowInc, colInc int
}

var (
	up    = direction{-1, 0}
	down  = direction{1, 0}
	left  = direction{0, -1}
	right = direction{0, 1}
)

func (d direction) String() string {
	switch d {
	case up:
		return "^"
	case down:
		return "v"
	case right:
		return ">"
	case left:
		return "<"
	}
	panic("dir")
}

type beam struct {
	node *node
	dir  direction
}

func (p pos) move(d direction) pos {
	return pos{row: p.row + d.rowInc, col: p.col + d.colInc}
}

type node struct {
	b     byte
	pos   pos
	beams []beam
}

func (node node) String() string {
	switch len(node.beams) {
	case 0:
		return string(node.b)
	case 1:
		if node.b == '.' {
			return node.beams[0].dir.String()
		}
		return string(node.b)
	default:
		if node.b == '.' {
			return fmt.Sprint(len(node.beams))
		}
		return string(node.b)
	}
}

func (node node) containsBeam(b beam) bool {
	for _, beam := range node.beams {
		if beam.node == b.node {
			return true
		}
	}
	return false
}

func (n node) getMoveDirections(incoming direction) []direction {
	rv := make([]direction, 0, 2)
	switch n.b {
	case '.':
		rv = append(rv, incoming)
	case '|':
		if incoming == up || incoming == down {
			rv = append(rv, incoming)
		} else {
			rv = append(rv, up, down)
		}
	case '-':
		if incoming == right || incoming == left {
			rv = append(rv, incoming)
		} else {
			rv = append(rv, right, left)
		}
	case '/':
		switch incoming {
		case up:
			rv = append(rv, right)
		case down:
			rv = append(rv, left)
		case left:
			rv = append(rv, down)
		case right:
			rv = append(rv, up)
		}
	case '\\':
		switch incoming {
		case up:
			rv = append(rv, left)
		case down:
			rv = append(rv, right)
		case right:
			rv = append(rv, down)
		case left:
			rv = append(rv, up)
		}
	}
	return rv
}

func (n node) move(m *matrix, p pos, incoming direction) []beam {
	rv := make([]beam, 0, 2)
	for _, d := range n.getMoveDirections(incoming) {
		newPos := p.move(d)
		if m.isValid(newPos) {
			rv = append(rv, beam{m.at(newPos), d})
		} else {
			rv = append(rv, beam{nil, d})
		}
	}
	return rv
}

type matrix struct {
	nodes        [][]node
	nRows, nCols int
}

func (m *matrix) Clone() matrix {
	rv := matrix{nodes: make([][]node, 0, m.nRows), nRows: m.nRows, nCols: m.nCols}

	for _, row := range m.nodes {
		newRow := make([]node, 0, m.nCols)
		for _, n := range row {
			newRow = append(newRow, node{b: n.b, pos: n.pos, beams: make([]beam, 0, 2)})
		}
		rv.nodes = append(rv.nodes, newRow)
	}
	return rv
}

func (m *matrix) at(pos pos) *node {
	return &m.nodes[pos.row][pos.col]
}

func (m *matrix) isValid(pos pos) bool {
	return pos.row >= 0 && pos.row < m.nRows && pos.col >= 0 && pos.col < m.nCols
}

func (m matrix) String() string {
	sb := strings.Builder{}
	for _, row := range m.nodes {
		for _, node := range row {
			sb.WriteString(node.String())
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseMatrix(f *os.File) matrix {
	rv := matrix{nodes: make([][]node, 0, 110)}
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)
	row := 0
	for lines.Scan() {
		newRow := make([]node, 0, rv.nCols)
		line := lines.Text()
		for col, b := range []byte(line) {
			newRow = append(newRow, node{b: b, pos: pos{row, col}, beams: make([]beam, 0)})
		}
		if rv.nCols == 0 {
			newRow = slices.Clip(newRow)
			rv.nCols = len(newRow)
		}
		rv.nodes = append(rv.nodes, newRow)
		row++
	}
	rv.nodes = slices.Clip(rv.nodes)
	rv.nRows = len(rv.nodes)

	return rv
}

func (m *matrix) countEnergized() int {
	rv := 0
	for _, row := range m.nodes {
		for _, n := range row {
			if len(n.beams) > 0 {
				rv++
			}
		}
	}
	return rv
}

type queue []beam

func (q *queue) push(n ...beam) {
	*q = append(*q, n...)
}

func (q *queue) pop() beam {
	rv := (*q)[0]
	*q = (*q)[1:]
	return rv
}

func (m *matrix) connectBeams(iniPos pos, initDirection direction) {
	var q = queue{}
	q.push(beam{m.at(iniPos), initDirection})
	for len(q) > 0 {
		beam := q.pop()
		node := beam.node
		if node == nil {
			continue
		}
		newBeams := node.move(m, node.pos, beam.dir)
		for _, newBeam := range newBeams {
			if !node.containsBeam(newBeam) {
				node.beams = append(node.beams, newBeam)
				q.push(newBeam)
			}
		}
	}
}

func A(f *os.File) int {
	matrix := parseMatrix(f)
	matrix.connectBeams(pos{0, 0}, right)
	return matrix.countEnergized()
}

func B(f *os.File) int {
	rv := 0

	matrix := parseMatrix(f)
	c := make(chan int)

	go func() {
		wg := sync.WaitGroup{}
		defer close(c)

		for i := 0; i < matrix.nRows; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				mm := matrix.Clone()
				mm.connectBeams(pos{i, 0}, right)
				c <- mm.countEnergized()
				wg.Done()
			}(&wg)
		}
		for i := 0; i < matrix.nRows; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				mm := matrix.Clone()
				mm.connectBeams(pos{i, matrix.nCols - 1}, left)
				c <- mm.countEnergized()
				wg.Done()
			}(&wg)
		}
		for i := 0; i < matrix.nCols; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				mm := matrix.Clone()
				mm.connectBeams(pos{0, i}, down)
				c <- mm.countEnergized()
				wg.Done()
			}(&wg)
		}
		for i := 0; i < matrix.nCols; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				mm := matrix.Clone()
				mm.connectBeams(pos{matrix.nRows - 1, i}, up)
				c <- mm.countEnergized()
				wg.Done()
			}(&wg)
		}
		wg.Wait()
	}()

	for v := range c {
		rv = max(rv, v)
	}

	return rv
}
