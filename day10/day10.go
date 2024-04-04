package day10

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type pipeInput byte

var pipeBytes = [...]pipeInput{'|', '-', 'L', 'J', '7', 'F', '.', 'S', 'X', ' ', 'Q'}

type pipe int

const (
	ns pipe = iota
	ew
	ne
	nw
	sw
	se
	ground
	start
	blocked
	empty
	inPath
)

func (p pipe) String() string {
	return string(pipeBytes[p])
}

func parsePipe(b byte) pipe {
	for i, v := range pipeBytes {
		if byte(v) == b {
			return pipe(i)
		}
	}
	log.Fatal("Bad input ", b)
	return '0'
}

type move int

const (
	n move = iota
	s
	e
	w
	sqn
	sqs
	sqe
	sqw
)

type pos struct {
	x, y int
}

type node struct {
	pos
	pipe       pipe
	prev, next *node
	inLoop     bool
	explored   bool
}

func (node node) String() string {
	return fmt.Sprint(node.pipe.String(), node.pos, " inside ", node.inLoop)
}

type matrix struct {
	m            [][]*node
	nRows, nCols int
}

func (m *matrix) at(pos pos) *node {
	return m.m[pos.y][pos.x]
}

func (m *matrix) atWithSqueeze(x, y int) *node {
	return m.m[y][x]
}

func (m matrix) String() string {
	sb := strings.Builder{}
	for _, row := range m.m {
		for _, elem := range row {
			sb.WriteString(elem.pipe.String())
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m matrix) inOut() string {
	sb := strings.Builder{}
	for i, row := range m.m {
		if i%2 == 1 {
			continue
		}
		for j, elem := range row {
			if j%2 == 1 {
				continue
			}
			if elem.inLoop || elem.pipe == blocked {
				sb.WriteString(elem.pipe.String())
			} else if elem.explored && elem.pipe != empty {
				sb.WriteString("O")
			} else {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseMatrix(f *os.File) *matrix {
	rv := matrix{}

	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	rv.m = make([][]*node, 0)
	y := 0
	for lines.Scan() {
		row := make([]*node, 0)
		emptyRow := make([]*node, 0)
		for x, b := range lines.Bytes() {
			row = append(row, &node{pipe: parsePipe(b), pos: pos{x: x * 2, y: y * 2}})
			row = append(row, &node{pipe: empty, pos: pos{x: x*2 + 1, y: y * 2}})
			emptyRow = append(emptyRow, &node{pipe: empty, pos: pos{x: x * 2, y: y*2 + 1}})
			emptyRow = append(emptyRow, &node{pipe: empty, pos: pos{x: x*2 + 1, y: y*2 + 1}})
		}
		rv.m = append(rv.m, row)
		rv.m = append(rv.m, emptyRow)
		y++
	}

	rv.nRows = len(rv.m)
	rv.nCols = len(rv.m[0])

	return &rv
}

func findStart(matrix *matrix) *node {
	for _, rows := range matrix.m {
		for _, elem := range rows {
			if elem.pipe == start {
				return elem
			}
		}
	}
	log.Fatal("Not found start")
	return nil
}

func (matrix *matrix) getSqueezeNodeInBetween(a *node, b *node) *node {
	if a.x > b.x && a.y == b.y {
		return matrix.atWithSqueeze(a.x-1, a.y)
	}
	if a.x < b.x && a.y == b.y {
		return matrix.atWithSqueeze(a.x+1, a.y)
	}
	if a.y > b.y && a.x == b.x {
		return matrix.atWithSqueeze(a.x, a.y-1)
	}
	if a.y < b.y && a.x == b.x {
		return matrix.atWithSqueeze(a.x, a.y+1)
	}
	log.Fatal("Error")
	return nil
}

func (matrix *matrix) blockSqueezeBetweenNodes(a *node, b *node) {
	matrix.getSqueezeNodeInBetween(a, b).pipe = blocked
}

func getMovesFromPipe(p pipe) []move {
	switch p {
	case ns:
		return []move{n, s}
	case ew:
		return []move{e, w}
	case ne:
		return []move{n, e}
	case nw:
		return []move{n, w}
	case sw:
		return []move{s, w}
	case se:
		return []move{s, e}
	case start:
		return []move{n, s, e, w}
	default:
		return []move{}
	}
}

func (pos pos) move(m move) pos {
	switch m {
	case n:
		pos.y -= 2
	case s:
		pos.y += 2
	case e:
		pos.x += 2
	case w:
		pos.x -= 2
	case sqn:
		pos.y--
	case sqs:
		pos.y++
	case sqe:
		pos.x++
	case sqw:
		pos.x--
	}
	return pos
}

func (m *matrix) isInLimits(pos pos) bool {
	return pos.x >= 0 && pos.x < m.nRows && pos.y >= 0 && pos.y < m.nCols
}

func (m *matrix) isValidMove(initialPos pos, move move) bool {
	movedPos := initialPos.move(move)
	if !m.isInLimits(movedPos) {
		return false
	}
	dest := m.at(movedPos)
	if dest.pipe == start {
		return true
	}
	p := dest.pipe
	switch move {
	case n:
		return p == ns || p == se || p == sw
	case s:
		return p == ns || p == ne || p == nw
	case e:
		return p == ew || p == sw || p == nw
	case w:
		return p == ew || p == ne || p == se
	}
	return false
}

func getConnectedNodes(matrix *matrix, pos pos) []*node {
	n := matrix.at(pos)
	moves := getMovesFromPipe(n.pipe)
	if n.pipe != start {
		return []*node{matrix.at(pos.move(moves[0])), matrix.at(pos.move(moves[1]))}
	}
	rv := make([]*node, 0)
	for _, move := range moves {
		if matrix.isValidMove(pos, move) {
			rv = append(rv, matrix.at(pos.move(move)))
		}
	}
	return rv
}

func lenToFinal(matrix *matrix, prev *node, current *node, last *node) int {
	nextNodes := getConnectedNodes(matrix, current.pos)
	if len(nextNodes) != 2 {
		log.Fatal("Err")
	}
	current.inLoop = true
	current.prev = prev
	if nextNodes[0] == prev {
		current.next = nextNodes[1]
	} else {
		current.next = nextNodes[0]
	}
	matrix.blockSqueezeBetweenNodes(current, current.next)
	if current.next != last {
		return lenToFinal(matrix, current, current.next, last) + 1
	}
	return 1
}

func initProc(f *os.File) (*matrix, int) {
	matrix := parseMatrix(f)
	start := findStart(matrix)
	start.inLoop = true
	connected := getConnectedNodes(matrix, start.pos)
	start.next = connected[0]
	start.prev = connected[1]
	start.next.prev = start
	matrix.blockSqueezeBetweenNodes(start, start.next)
	matrix.blockSqueezeBetweenNodes(start, start.prev)

	len := lenToFinal(matrix, start, start.next, start)

	fmt.Println(matrix)
	return matrix, len
}

func Day10a(f *os.File) int {
	_, len := initProc(f)

	return len/2 + 1
}

func (matrix *matrix) getAdjacentNodesNotProcessed(n *node) []*node {
	rv := make([]*node, 0)
	moves := []move{sqn, sqs, sqe, sqw}

	for _, move := range moves {
		newPos := n.pos.move(move)
		if !matrix.isInLimits(newPos) {
			continue
		}
		newNode := matrix.at(newPos)
		if !newNode.explored {
			rv = append(rv, newNode)
		}
	}
	return rv
}

type queue []*node

func (queue *queue) empty() bool {
	return len(*queue) == 0
}

func (queue *queue) push(node *node) {
	*queue = append(*queue, node)
}

func (queue *queue) pop() *node {
	rv := (*queue)[len(*queue)-1]
	*queue = (*queue)[:len(*queue)-1]
	return rv
}

func outside(matrix *matrix) int {
	rv := 0
	q := queue{}
	start := matrix.at(pos{0, 0})
	q.push(start)
	rv++

	for !q.empty() {
		node := q.pop()
		if node.explored {
			continue
		}
		for _, adjNode := range matrix.getAdjacentNodesNotProcessed(node) {
			if !adjNode.inLoop && adjNode.pipe != blocked {
				q.push(adjNode)
			} else {
				adjNode.explored = true
			}
		}
		if node.pipe != empty {
			rv++
		}
		node.explored = true
	}

	return rv
}

func Day10b(f *os.File) int {
	matrix, len := initProc(f)

	out := outside(matrix)

	fmt.Println(matrix.inOut())
	inside := (matrix.nCols/2)*(matrix.nRows/2) - len - out

	return inside
}
