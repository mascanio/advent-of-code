package day14

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

type matrix [][]byte

func (m matrix) String() string {
	sb := strings.Builder{}
	for _, v := range m {
		for _, e := range v {
			sb.WriteByte(e)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (a *matrix) eq(b *matrix) bool {
	for i := 0; i < len((*a)); i++ {
		for j := 0; j < len((*a)[0]); j++ {
			if (*a)[i][j] != (*b)[i][j] {
				return false
			}
		}
	}
	return true
}

func (a *matrix) copy() matrix {
	rv := make(matrix, 0, len(*a))
	for i := 0; i < len((*a)); i++ {
		rv = append(rv, slices.Clone((*a)[i]))
	}
	return rv
}

func parseInput(f *os.File) matrix {
	rv := make([][]byte, 0)
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		rv = append(rv, []byte(lines.Text()))
	}
	return rv
}

func move(m *matrix, is, ie, js, je int, accessRowsFirst bool) {
	iinc, jinc := 1, 1
	if is > ie {
		iinc = -1
	}
	if js > je {
		jinc = -1
	}
	if accessRowsFirst {
		for i := is; i != ie; i += iinc {
			for j := js; j != je; j += jinc {
				if (*m)[i][j] != '.' {
					continue
				}
				for k := i + iinc; k != ie; k += iinc {
					if (*m)[k][j] == '#' {
						break
					}
					if (*m)[k][j] == 'O' {
						(*m)[i][j] = 'O'
						(*m)[k][j] = '.'
						break
					}
				}
			}
		}
	} else {
		for j := js; j != je; j += jinc {
			for i := is; i != ie; i += iinc {
				if (*m)[i][j] != '.' {
					continue
				}
				for k := j + jinc; k != je; k += jinc {
					if (*m)[i][k] == '#' {
						break
					}
					if (*m)[i][k] == 'O' {
						(*m)[i][j] = 'O'
						(*m)[i][k] = '.'
						break
					}
				}
			}
		}
	}
}

func moveUp(m *matrix) {
	move(m, 0, len(*m), 0, len((*m)[0]), true)
}
func moveDown(m *matrix) {
	move(m, len(*m)-1, -1, len((*m)[0])-1, -1, true)
}

func moveRight(m *matrix) {
	move(m, 0, len(*m), len((*m)[0])-1, -1, false)
}

func moveLeft(m *matrix) {
	move(m, 0, len(*m), 0, len((*m)[0]), false)
}

func count(m *matrix) int {
	rv := 0
	for i := 0; i < len(*m); i++ {
		count := 0
		for j := 0; j < len((*m)[0]); j++ {
			if (*m)[i][j] == 'O' {
				count++
			}
		}
		rv += count * (len(*m) - i)
	}
	return rv
}

func A(f *os.File) int {
	rv := 0
	input := parseInput(f)
	moveUp(&input)
	rv = count(&input)
	return rv
}

type result struct {
	result, pos int
	m           matrix
}

func B(f *os.File) int {
	input := parseInput(f)
	results := make([]result, 0)
	loopLen := 0
	loopStart := 0
	for i := 0; i < 1000; i++ {
		moveUp(&input)
		moveLeft(&input)
		moveDown(&input)
		moveRight(&input)
		count := count(&input)
		loopFound := false
		for _, v := range results {
			if count == v.result && v.m.eq(&input) {
				loopFound = true
				loopLen = i - v.pos
				loopStart = v.pos
				break
			}
		}
		if loopFound {
			break
		}
		results = append(results, result{m: input.copy(), result: count, pos: i})
	}

	x := (1000000000 - loopStart - 1) % loopLen
	return results[loopStart+x].result
}
