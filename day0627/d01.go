package main

import "fmt"

func main() {
	number := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Printf("%p\n", number)
	number = append(number, 8)
	fmt.Printf("%p\n", number)
}
