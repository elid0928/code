package main

import "fmt"

func main() {

	foo := []int{0, 0, 0, 42, 100}
	bar := foo[1:4:4]
	fmt.Printf("foo address, %p\n", &foo)
	fmt.Printf("bar address, %p\n", &bar)
	bar = append(bar, 99)
	bar[0] = 42
	fmt.Println("foo:", foo) // foo: [0 0 0 42 100]
	fmt.Println("bar:", bar) // bar: [0 0 42 99]

	fmt.Printf("foo address, %p\n", &foo)
	fmt.Printf("bar address, %p\n", &bar)
}
