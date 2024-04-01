package day08

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

type node struct {
	name string
	r    *node
	l    *node
}

type nodeParse struct {
	l, r   string
	parsed *node
}

type m struct {
	root       node
	directions string
}

func parseTree(f *os.File) (rv m) {
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)

	lines.Scan()
	rv.directions = lines.Text()

	parseTree := parseNodes(lines)
	rv.root = *createNode(&parseTree, "AAA")

	return
}

func parseTreeB(f *os.File) ([]*node, string) {
	rv := make([]*node, 0)
	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)

	lines.Scan()
	directions := lines.Text()
	parseTree := parseNodes(lines)

	for k := range parseTree {
		if k[len(k)-1] == 'A' {
			rv = append(rv, createNode(&parseTree, k))
		}
	}

	return rv, directions
}

func createNode(parseTree *map[string](*nodeParse), name string) *node {
	parseNode := (*parseTree)[name]
	node := node{name: name}
	if parseNode.parsed == nil {
		parseNode.parsed = &node
		if (*parseTree)[parseNode.l].parsed != nil {
			node.l = (*parseTree)[parseNode.l].parsed
		} else {
			node.l = createNode(parseTree, parseNode.l)
		}
		if (*parseTree)[parseNode.r].parsed != nil {
			node.r = (*parseTree)[parseNode.r].parsed
		} else {
			node.r = createNode(parseTree, parseNode.r)
		}
	}
	return &node
}

func parseNodes(lines *bufio.Scanner) map[string](*nodeParse) {
	nodes := make(map[string](*nodeParse))
	reNode := regexp.MustCompile(`(...) = \((...), (...)\)`)
	for lines.Scan() {
		line := lines.Text()
		m := reNode.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		nodes[m[1]] = &nodeParse{l: m[2], r: m[3]}
	}
	return nodes
}

func solve(node *node, directions string, depth int) int {
	if node.name == "ZZZ" {
		return 0
	}
	direction := directions[depth%len(directions)]
	switch direction {
	case 'L':
		return solve(node.l, directions, depth+1) + 1
	case 'R':
		return solve(node.r, directions, depth+1) + 1
	}
	log.Fatal("WTF")
	return 33
}

func solveCycle(node *node, directions string, depth int) int {
	for {
		if strings.HasSuffix(node.name, "Z") {
			return depth
		}
		direction := directions[depth%len(directions)]
		switch direction {
		case 'L':
			node = node.l
		case 'R':
			node = node.r
		}
		depth++
	}
}

func Day08a(f *os.File) int {
	rv := 0
	m := parseTree(f)
	rv = solve(&m.root, m.directions, 0)
	return rv
}

func gdc(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lmc(a, b int) int {
	return (a * b) / gdc(a, b)
}

func Day08b(f *os.File) int {
	rv := 0
	parsedTrees, directions := parseTreeB(f)

	for _, v := range parsedTrees {
		tmp := solveCycle(v, directions, 0)
		if rv == 0 {
			rv = tmp
		} else {
			rv = lmc(rv, tmp)
		}
	}

	return rv
}
