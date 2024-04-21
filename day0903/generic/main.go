package main

import (
	"fmt"
)

func main() {

	var a, b int = 1, 2
	fmt.Println(Sum([]int{a, b}))
	fmt.Println(Sum([]float64{1.2, 1.3}))
	// math.MaxInt
}

type Number interface {
	int | int64 | float64
}

func Sum[T Number](numbers []T) T {
	var total T
	for _, x := range numbers {
		total = total + x
	}
	return total
}
