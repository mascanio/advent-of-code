package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mascanio/advent-of-code/day05"
)

func main() {
	defer func(t time.Time) {
		fmt.Println(time.Since(t))
	}(time.Now())
	path := os.Args[1]
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Println(day05.Day05b(f))
}
