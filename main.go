package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mascanio/advent-of-code/day03"
)

func main() {
	defer func(t time.Time) {
		fmt.Println(time.Since(t))
	}(time.Now())
	fmt.Println(day03.Day03b(os.Args[1]))
}
