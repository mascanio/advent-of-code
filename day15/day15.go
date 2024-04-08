package day15

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func hash(s []byte) int {
	rv := 0

	for _, b := range s {
		rv += int(b)
		rv *= 17
		rv = rv % 256
	}

	return rv
}

func A(f *os.File) int {
	rv := 0
	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)
	lines.Scan()
	line := lines.Text()
	for _, l := range strings.Split(line, ",") {
		rv += hash([]byte(l))
	}

	return rv
}

type lens struct {
	id    string
	power int
}

type box []lens

func remove(box *box, a lens) {
	for i := 0; i < len(*box); i++ {
		if (*box)[i].id == a.id {
			*box = append((*box)[:i], (*box)[i+1:]...)
			break
		}
	}
}

func add(box *box, a lens) {
	for i := 0; i < len(*box); i++ {
		if (*box)[i].id == a.id {
			(*box)[i].power = a.power
			return
		}
	}
	*box = append(*box, a)
}

func B(f *os.File) int {
	rv := 0

	lines := bufio.NewScanner(bufio.NewReader(f))
	lines.Split(bufio.ScanLines)
	lines.Scan()
	line := lines.Text()
	boxes := [256]box{}
	for _, l := range strings.Split(line, ",") {
		if strings.Contains(l, "=") {
			ss := strings.Split(l, "=")
			id := ss[0]
			n, _ := strconv.Atoi(ss[1])
			h := hash([]byte(id))
			add(&boxes[h], lens{id: id, power: n})
		} else if strings.Contains(l, "-") {
			id := l[:len(l)-1]
			h := hash([]byte(id))
			remove(&boxes[h], lens{id: id})
		}
	}

	for i := 0; i < len(boxes); i++ {
		box := boxes[i]
		boxId := i + 1
		for j := 0; j < len(box); j++ {
			lenId := j + 1
			rv += boxId * lenId * box[j].power
		}
	}

	return rv
}
