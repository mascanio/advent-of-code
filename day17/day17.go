package day17

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type nodeMap[T any] struct {
	nodes        [][]*T
	nCols, nRows int
}

const maxRem = int(10)

type pos struct {
	row, col int
}

type direction struct {
	rowInc, colInc int
}

var (
	none  = direction{0, 0}
	up    = direction{-1, 0}
	down  = direction{1, 0}
	left  = direction{0, -1}
	right = direction{0, 1}
)

var directions = [4]direction{up, down, left, right}

func (d direction) Int() int {
	switch d {
	case up:
		return 0
	case down:
		return 1
	case right:
		return 2
	case left:
		return 3
	}
	panic("dir")
}

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

func (p pos) move(d direction) pos {
	return pos{row: p.row + d.rowInc, col: p.col + d.colInc}
}

type node struct {
	heatDisp int
	pos      pos
	solDir   direction
}

func (nodes *nodeMap[T]) isValid(pos pos) bool {
	return pos.col >= 0 && pos.col < nodes.nCols && pos.row >= 0 && pos.row < nodes.nRows
}

func (nodes *nodeMap[T]) at(pos pos) *T {
	if !nodes.isValid(pos) {
		return nil
	}
	return nodes.nodes[pos.row][pos.col]
}

type movement struct {
	dir direction
	pre int
}

func (n node) getValidMovements(incoming direction, pre int) []movement {
	rv := make([]movement, 0)
	if incoming == none {
		rv = append(rv, movement{right, 0}, movement{down, 0})
		return rv
	}

	if pre < 3 {
		rv = append(rv, movement{incoming, pre + 1})
		return rv
	} else if pre+1 < maxRem {
		rv = append(rv, movement{incoming, pre + 1})
	}
	switch incoming {
	case up, down:
		rv = append(rv, movement{left, 0}, movement{right, 0})
	case left, right:
		rv = append(rv, movement{down, 0}, movement{up, 0})
	}

	return rv
}

type matrix struct {
	nodeMap[node]
}

func (matrix matrix) String() string {
	sb := strings.Builder{}

	for _, row := range matrix.nodes {
		for _, v := range row {
			a := direction{}
			if v.solDir != a {
				sb.WriteString(fmt.Sprint(v.solDir))
			} else {
				sb.WriteString(fmt.Sprint(v.heatDisp))
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func parseInput(f *os.File) matrix {
	var rv = matrix{}
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)
	row := 0
	for lines.Scan() {
		l := lines.Text()
		newRow := make([]*node, 0)
		for col, b := range l {
			n, _ := strconv.Atoi(string(b))
			newRow = append(newRow, &node{heatDisp: n, pos: pos{row, col}})
		}
		rv.nodes = append(rv.nodes, newRow)
		row++
	}
	rv.nRows = len(rv.nodes)
	rv.nCols = len(rv.nodes[0])
	return rv
}

type queueItem struct {
	node        *node
	index, prio int
	incomingDir direction
	preMoves    int
	prev        *queueItem
}

type PriorityQueue []*queueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].prio < pq[j].prio
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*queueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *queueItem, distance int) {
	item.prio = distance
	heap.Fix(pq, item.index)
}

func getMinOf(m map[mapKey]*queueItem, pos pos) *queueItem {
	curMin := math.MaxInt
	var curNode *queueItem
	for _, dir := range directions {
		for pre := 0; pre < maxRem; pre++ {
			qItem, exists := m[mapKey{pos, dir, pre}]
			if !exists {
				continue
			}
			fmt.Println(qItem, " ", qItem.prio)
			if qItem.prio < curMin {
				curMin = qItem.prio
				curNode = qItem
			}
		}
	}
	return curNode
}

type mapKey struct {
	pos pos
	dir direction
	pre int
}

func dijkstra(g *matrix, s *node) int {
	pq := make(PriorityQueue, 0, g.nCols*g.nRows)
	endPos := pos{g.nRows - 1, g.nCols - 1}
	zeroPos := pos{0, 0}
	nm := make(map[mapKey]*queueItem)

	qItem := &queueItem{node: s}
	qItem.prio = 0
	qItem.incomingDir = none
	qItem.preMoves = 0
	nm[mapKey{s.pos, none, 0}] = qItem
	heap.Push(&pq, qItem)

	for len(pq) != 0 {
		u := heap.Pop(&pq).(*queueItem)
		for _, movement := range u.node.getValidMovements(u.incomingDir, u.preMoves) {
			newPos := u.node.pos.move(movement.dir)
			if !g.isValid(newPos) || newPos == zeroPos || (newPos == endPos && movement.pre < 3) {
				continue
			}
			var v *queueItem
			var exists bool
			if v, exists = nm[mapKey{newPos, movement.dir, movement.pre}]; !exists {
				v = &queueItem{node: g.at(newPos)}
				v.prio = u.prio + v.node.heatDisp
				v.incomingDir = movement.dir
				v.preMoves = movement.pre
				v.prev = u
				nm[mapKey{newPos, movement.dir, movement.pre}] = v
				heap.Push(&pq, v)
			} else {
				alt := u.prio + v.node.heatDisp
				if alt < v.prio {
					v.prio = alt
					v.prev = u
					pq.update(v, alt)
				}
			}
		}
	}
	rv := getMinOf(nm, endPos)

	prevN := rv
	for {
		fmt.Println(prevN.node.pos)
		if prevN.prev == nil {
			break
		}
		prevN = prevN.prev
	}

	return rv.prio
}

func A(f *os.File) int {
	return 0
}

func B(f *os.File) int {
	matrix := parseInput(f)
	return dijkstra(&matrix, matrix.at(pos{0, 0}))
}
