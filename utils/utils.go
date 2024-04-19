package utils

import (
	"fmt"
	"strings"
)

type Pos struct {
	Row, Col int
}

type Direction struct {
	RowInc, ColInc int
}

var (
	None  = Direction{0, 0}
	Up    = Direction{-1, 0}
	Down  = Direction{1, 0}
	Left  = Direction{0, -1}
	Right = Direction{0, 1}
)

var Directions = [4]Direction{Up, Down, Left, Right}

func (d Direction) Int() int {
	switch d {
	case Up:
		return 0
	case Down:
		return 1
	case Right:
		return 2
	case Left:
		return 3
	}
	panic("dir")
}

func (d Direction) String() string {
	switch d {
	case Up:
		return "^"
	case Down:
		return "v"
	case Right:
		return ">"
	case Left:
		return "<"
	}
	panic("dir")
}

func (p Pos) Move(d Direction) Pos {
	return Pos{Row: p.Row + d.RowInc, Col: p.Col + d.ColInc}
}

func (p Pos) MoveN(d Direction, n int) Pos {
	return Pos{Row: p.Row + d.RowInc*n, Col: p.Col + d.ColInc*n}
}

type NodeMap[T any] struct {
	Nodes      [][]T
	Cols, Rows int
}

func (nodes *NodeMap[T]) IsValid(pos Pos) bool {
	return pos.Col >= 0 && pos.Col < nodes.Cols && pos.Row >= 0 && pos.Row < nodes.Rows
}

func (nodes *NodeMap[T]) At(pos Pos) T {
	return nodes.Nodes[pos.Row][pos.Col]
}

func (nodes *NodeMap[T]) AtP(pos Pos) *T {
	return &nodes.Nodes[pos.Row][pos.Col]
}

func (nodes *NodeMap[T]) Set(pos Pos, val T) {
	nodes.Nodes[pos.Row][pos.Col] = val
}

func (nodes *NodeMap[T]) GetAdjacentNodes(pos Pos, d []Direction, validAdjacent func(candidate, current T, candidateP, currentP Pos) bool) []*T {
	if !nodes.IsValid(pos) {
		return nil
	}
	rv := make([]*T, 0, len(d))
	curNode := nodes.At(pos)
	for _, dir := range d {
		newPos := pos.Move(dir)
		if !nodes.IsValid(newPos) {
			continue
		}
		newNode := nodes.AtP(newPos)
		if !validAdjacent(*newNode, curNode, newPos, pos) {
			continue
		}
		rv = append(rv, newNode)
	}

	return rv
}

func (nodes *NodeMap[T]) String() string {
	sb := strings.Builder{}
	for i := 0; i < nodes.Rows; i++ {
		for j := 0; j < nodes.Cols; j++ {
			sb.WriteString(fmt.Sprint(nodes.Nodes[i][j]))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func CreateNodeMap[T any](rows, cols int) NodeMap[T] {
	rv := NodeMap[T]{Rows: rows, Cols: cols}
	rv.Nodes = make([][]T, rows)
	for i := 0; i < rows; i++ {
		rv.Nodes[i] = make([]T, cols)
	}
	return rv
}

type Queue[T any] []T

func (q *Queue[T]) Pop() T {
	rv := (*q)[0]
	*q = (*q)[1:]
	return rv
}

func (q *Queue[T]) Push(elem T) {
	*q = append(*q, elem)
}
