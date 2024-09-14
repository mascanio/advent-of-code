package main

import (
	"fmt"
	"log"
	"os"
	"time"

	d "github.com/mascanio/advent-of-code/day19"
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
		fmt.Println(d.A(f))
	case "b":
		fmt.Println(d.B(f))
	}
}
