package day11

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
)

type nodeMap[T any] struct {
	nodes        [][]*T
	nCols, nRows int
}

func createNodeMap[T any](nCols, nRows int) nodeMap[T] {
	rv := nodeMap[T]{nCols: nCols, nRows: nRows}

	rv.nodes = make([][]*T, nRows)
	for i := range rv.nodes {
		rv.nodes[i] = make([]*T, nCols)
	}
	return rv
}

type graph struct {
	nodeMap[node]
	galaxies []*node
	rowDist  []int
	colDist  []int
}

type pos struct {
	x, y int
}

type node struct {
	pos
	isGalaxy  bool
	galaxyNum int
}

func (node node) String() string {
	if node.isGalaxy {
		return fmt.Sprint(node.galaxyNum)
	}
	return "."
}

func (graph graph) String() string {
	sb := strings.Builder{}
	for _, v := range graph.nodes {
		for _, node := range v {
			sb.WriteString(node.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (nodes *nodeMap[T]) isValid(pos pos) bool {
	return pos.x >= 0 && pos.x < nodes.nCols && pos.y >= 0 && pos.y < nodes.nRows
}

func (nodes *nodeMap[T]) at(pos pos) *T {
	if !nodes.isValid(pos) {
		return nil
	}
	return nodes.nodes[pos.y][pos.x]
}

func (nodes *nodeMap[T]) set(pos pos, elem *T) {
	nodes.nodes[pos.y][pos.x] = elem
}

func (a *node) distance(b *node, g *graph) int {
	if a.x == b.x {
		return g.rowDist[min(a.y, b.y)]
	}
	return g.colDist[min(a.x, b.x)]
}

func (galaxy *graph) expand(n int) {
	// expand rows
	for row := 0; row < galaxy.nRows; row++ {
		noGalaxies := true
		for _, node := range galaxy.nodes[row] {
			if node.isGalaxy {
				noGalaxies = false
				break
			}
		}
		if noGalaxies {
			galaxy.rowDist[row] = n
		} else {
			galaxy.rowDist[row] = 1
		}
	}
	// expand cols
	for col := 0; col < galaxy.nCols; col++ {
		noGalaxies := true
		for row := 0; row < galaxy.nRows; row++ {
			node := galaxy.nodes[row][col]
			if node.isGalaxy {
				noGalaxies = false
				break
			}
		}
		if noGalaxies {
			galaxy.colDist[col] = n
		} else {
			galaxy.colDist[col] = 1
		}
	}
}

func (g *graph) setPosAndGalaxyNum() {
	galaxyNum := 1
	for row := 0; row < g.nRows; row++ {
		for col := 0; col < g.nCols; col++ {
			node := g.at(pos{x: col, y: row})
			node.pos = pos{x: col, y: row}
			if node.isGalaxy {
				node.galaxyNum = galaxyNum
				g.galaxies = append(g.galaxies, node)
				galaxyNum++
			}
		}
	}
}

func parseGalaxy(f *os.File) graph {
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)
	rv := graph{}
	rv.nodes = make([][]*node, 0)
	nRow := 0
	for lines.Scan() {
		newRow := make([]*node, 0)
		for _, v := range lines.Bytes() {
			newRow = append(newRow, &node{isGalaxy: v == '#'})
		}
		rv.nodes = append(rv.nodes, newRow)
		nRow++
	}
	rv.nRows = nRow
	rv.nCols = len(rv.nodes[0])
	rv.colDist = make([]int, rv.nCols*2)
	rv.rowDist = make([]int, rv.nRows*2)
	rv.setPosAndGalaxyNum()
	return rv
}

type queueItem struct {
	node        *node
	index, prio int
	visited     bool
}

func (n *queueItem) getNodesAdjacentForwardUnvisited(nodeMap *nodeMap[queueItem]) []*node {
	// Dont go row up. Also, don't return visited nodes
	rv := make([]*node, 0, 3)
	{
		newPos := n.node.pos
		newPos.x++
		if newNode := nodeMap.at(newPos); newNode != nil && !newNode.visited {
			rv = append(rv, newNode.node)
		}
	}
	{
		newPos := n.node.pos
		newPos.x--
		if newNode := nodeMap.at(newPos); newNode != nil && !newNode.visited {
			rv = append(rv, newNode.node)
		}
	}
	{
		newPos := n.node.pos
		newPos.y++
		if newNode := nodeMap.at(newPos); newNode != nil && !newNode.visited {
			rv = append(rv, newNode.node)
		}
	}
	return rv
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

func dijkstra(g *graph, s *node) []int {
	pq := make(PriorityQueue, 0, g.nCols*g.nRows)
	nodeMap := createNodeMap[queueItem](g.nCols, g.nRows)

	start := &queueItem{node: s}
	start.prio = 0
	start.index = 0
	pq = append(pq, start)
	idx := 1
	for i, row := range g.nodes {
		if i < s.y {
			// Start looking forward galaxies only
			continue
		}
		for j, n := range row {
			if i == s.y && j < s.x {
				// Start looking forward galaxies only
				continue
			}
			if n != s {
				newNode := &queueItem{node: n}
				nodeMap.set(n.pos, newNode)
				newNode.prio = math.MaxInt - 1
				newNode.visited = false
				newNode.index = idx
				idx++
				pq = append(pq, newNode)
			}
		}
	}
	heap.Init(&pq)
	for len(pq) != 0 {
		u := heap.Pop(&pq).(*queueItem)
		for _, v := range u.getNodesAdjacentForwardUnvisited(&nodeMap) {
			d := u.node.distance(v, g)
			alt := u.prio + d
			queueNode := nodeMap.at(v.pos)
			if alt < queueNode.prio {
				queueNode.prio = alt
				pq.update(queueNode, alt)
			}
		}
	}
	rv := make([]int, len(g.galaxies))
	for _, nn := range nodeMap.nodes {
		for _, v := range nn {
			if v != nil && v.node.isGalaxy {
				rv[g.at(v.node.pos).galaxyNum-1] = v.prio
			}
		}
	}
	return rv
}

func proc(f *os.File, n int) int {
	rv := 0
	graph := parseGalaxy(f)
	graph.expand(n)

	c := make(chan int)

	go func() {
		wg := sync.WaitGroup{}
		defer close(c)
		for i, galaxy := range graph.galaxies {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				partial := 0
				distances := dijkstra(&graph, galaxy)
				for j := i + 1; j < len(distances); j++ {
					partial += distances[j]
				}
				c <- partial
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

func Day11a(f *os.File) int {
	return proc(f, 2)
}

func Day11b(f *os.File) int {
	return proc(f, 1000000)
}
