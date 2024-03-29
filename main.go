package main

import (
	"fmt"
	"time"
)

func main() {
	defer func(t time.Time) {
		fmt.Println(time.Since(t))
	}(time.Now())
	fmt.Println(Day01a())
}
