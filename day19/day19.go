package day19

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type operator byte

func (op operator) operate(val, lim int) bool {
	if op == '<' {
		return val < lim
	}
	return val > lim
}

func (op operator) interval(val int) validInterval {
	if op == '<' {
		return validInterval{0, val}
	}
	return validInterval{val, 4000}
}

type partRating int

const (
	x partRating = iota
	m
	a
	s
)

func (p partRating) String() string {
	ar := [...]string{"x", "m", "a", "s"}
	return ar[p]
}

func parsePartRating(b string) partRating {
	switch b {
	case "x":
		return x
	case "m":
		return m
	case "a":
		return a
	case "s":
		return s
	}
	panic("part " + b)
}

type stepType int

const (
	invalid stepType = iota
	wf
	rule
	stepNotApplies
)

func (s stepType) String() string {
	switch s {
	case stepNotApplies:
		return "stepNotApplies"
	case wf:
		return "wf"
	case rule:
		return "rule"
	}
	panic(s)
}

type step struct {
	operator
	t       stepType
	part    partRating
	lim     int
	next    acceptor
	nextStr string
}

func (s step) String() string {
	if s.t == rule {
		return fmt.Sprintf("%v%v%v:%v", s.part, string(s.operator), s.lim, s.nextStr)
	}
	return s.nextStr
}

func (s step) Step(p part) (stepType, acceptor) {
	switch s.t {
	case invalid, stepNotApplies:
		panic("Invalid step")
	case rule:
		if s.operate(p.ratings[s.part], s.lim) {
			return wf, s.next
		}
		return stepNotApplies, nil
	case wf:
		return wf, s.next
	}
	panic("invalid step")
}

func (s step) ValidIntervals() intervals {
	switch s.t {
	case invalid, stepNotApplies:
		panic("Invalid step")
	case rule:
		next := s.next.ValidIntervals()
		next[s.part] = s.operator.interval(s.lim).intersec(next[s.part])
		return next
	case wf:
		return s.next.ValidIntervals()
	}
	panic("invalid step")
}

func parseStepper(s string) step {
	reStep := regexp.MustCompile(`([xmas])([<>])(\d+)\:(\w+)`)
	m := reStep.FindStringSubmatch(s)
	if len(m) > 0 {
		lim, _ := strconv.Atoi(m[3])
		return step{t: rule, part: parsePartRating(m[1]), operator: operator(m[2][0]), lim: lim, nextStr: m[4]}
	}
	return step{t: wf, nextStr: s}
}

type validInterval struct {
	lo, hi int
}

func (a validInterval) intersecs(b validInterval) bool {
	if a.lo > b.lo {
		a, b = b, a
	}
	return b.lo <= a.hi
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (a validInterval) intersec(intervals []validInterval) []validInterval {
	// Precond intervals sorted and dont intersec
	rv := make([]validInterval, 0)
	for i := 0; i < len(intervals)-1; i++ {
		if intervals[i].lo >= intervals[i+i].lo || intervals[i].intersecs(intervals[i+1]) {
			panic("Invalid pre intersec")
		}
	}
	for i := 0; i < len(intervals); i++ {
		b := intervals[i]
		if !a.intersecs(b) {
			if a.hi < b.lo {
				break
			}
			continue
		}
		rv = append(rv, validInterval{max(a.lo, b.lo), min(a.hi, b.hi)})
	}
	return rv
}

func (a validInterval) union(intervals []validInterval) []validInterval {
	// Precond intervals sorted and dont intersec
	rv := make([]validInterval, 0)
	for i := 0; i < len(intervals)-1; i++ {
		if intervals[i].lo >= intervals[i+i].lo || intervals[i].intersecs(intervals[i+1]) {
			panic("Invalid pre intersec")
		}
	}
	last := intervals[len(intervals)-1]
	if last.hi < a.lo {
		return append(rv, a)
	}
	for i := 0; i < len(intervals); i++ {
		b := intervals[i]
		if a.intersecs(b) {
			a = validInterval{min(a.lo, b.lo), max(a.hi, b.hi)}
		} else if b.hi < a.lo {
			rv = append(rv, a)
			rv = append(rv, b)
		} else {
			rv = append(rv, b)
		}
	}
	return rv
}

/*
func (a validInterval) intersec(intervals []validInterval) validInterval {
	b := intervals[0]
	if a.lo > b.lo {
		a, b = b, a
	}
	if a.hi < b.lo {
		return validInterval{0, 0}
	}
	return validInterval{a.hi, b.lo}
}*/

/*
	func (a validInterval) union(b validInterval) []validInterval {
		if a.lo > b.lo {
			a, b = b, a
		}
		// a.lo <= b.lo
		if a.hi >= b.hi {
			return []validInterval{a}
		}
		if b.lo <= a.hi {
			return []validInterval{validInterval{a.lo, b.hi}}
		}
		return []validInterval{a, b}
	}
*/
type intervals [4][]validInterval

type acceptor interface {
	Accept(part) bool
	ValidIntervals() intervals
}

type workflow struct {
	name  string
	steps []step
}

func (w workflow) Accept(p part) bool {
	for _, step := range w.steps {
		t, wf := step.Step(p)
		if t == stepNotApplies {
			continue
		}
		return wf.Accept(p)
	}
	panic("A")
}

func (w workflow) ValidIntervals() intervals {
	for _, step := range w.steps {

	}
}

type WAccept struct{}

func (w WAccept) Accept(_ part) bool { return true }
func (w WAccept) ValidIntervals() intervals {
	return intervals{{validInterval{0, 4000}}, {validInterval{0, 4000}}, {validInterval{0, 4000}}, {validInterval{0, 4000}}}
}
func (w WAccept) String() string { return "A" }

type WReject struct{}

func (w WReject) Accept(_ part) bool { return false }
func (w WReject) String() string     { return "R" }
func (w WReject) ValidIntervals() intervals {
	return intervals{{validInterval{0, 0}}, {validInterval{0, 0}}, {validInterval{0, 0}}, {validInterval{0, 0}}}
}

func (w workflow) String() string {
	sb := strings.Builder{}

	sb.WriteString(w.name + " [")
	for _, v := range w.steps {
		sb.WriteString(fmt.Sprint(v) + " ")
	}
	sb.WriteString("]\n")
	return sb.String()
}

type part struct {
	ratings [4]int
}

func parsePart(s string) part {
	rv := part{}
	s = strings.TrimSpace(s)[1 : len(s)-1]
	for _, v := range strings.Split(s, ",") {
		partS := strings.Split(v, "=")
		value, _ := strconv.Atoi(partS[1])
		rv.ratings[parsePartRating(partS[0])] = value
	}
	return rv
}

type input struct {
	workflows   []workflow
	parts       []part
	workflowMap map[string]acceptor
}

func (i input) String() string {
	sb := strings.Builder{}

	for _, v := range i.workflows {
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintln(i.parts))
	sb.WriteString("\n")
	sb.WriteString("\n")

	return sb.String()
}

func parseWorkflow(s string) workflow {
	rv := workflow{steps: make([]step, 0)}
	reName := regexp.MustCompile(`(.+){(.+)}`)
	m := reName.FindStringSubmatch(s)
	rv.name = m[1]
	for _, v := range strings.Split(m[2], ",") {
		rv.steps = append(rv.steps, parseStepper(v))
	}
	return rv
}

func parseInput(f *os.File) input {
	rv := input{workflows: make([]workflow, 0), parts: make([]part, 0), workflowMap: make(map[string]acceptor)}

	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)

	for lines.Scan() {
		line := lines.Text()
		if strings.TrimSpace(line) == "" {
			break
		}
		workflow := parseWorkflow(line)
		rv.workflows = append(rv.workflows, workflow)
		rv.workflowMap[workflow.name] = workflow
	}
	rv.workflowMap["A"] = WAccept{}
	rv.workflowMap["R"] = WReject{}
	for _, workflow := range rv.workflows {
		for i := 0; i < len(workflow.steps); i++ {
			workflow.steps[i].next = rv.workflowMap[workflow.steps[i].nextStr]
		}
	}
	// Parse parts
	for lines.Scan() {
		line := lines.Text()
		rv.parts = append(rv.parts, parsePart(line))
	}

	return rv
}

func A(f *os.File) int {
	rv := 0
	input := parseInput(f)
	for _, part := range input.parts {
		if input.workflowMap["in"].Accept(part) {
			for _, v := range part.ratings {
				rv += v
			}
		}
	}
	return rv
}

func B(f *os.File) int {
	return 0
}
