package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mascanio/advent-of-code/day12"
)

func main() {
	defer func(t time.Time) {
		fmt.Println(time.Since(t))
	}(time.Now())
	path := os.Args[2]
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	switch os.Args[1] {
	case "a":
		fmt.Println(day12.Day12a(f))
	case "b":
		fmt.Println(day12.Day12b(f))
	}
}
