package main

import "fmt"

type User struct {
	Name string
	Age  int32
}

func main() {
	var m, n int
	m = 12
	n = 34
	sum := 0
	m = n + 1
	for i := 0; i < 5; i++ {
		sum += m
	}
	u := User{
		Name: "liujiadong",
		Age:  32,
	}

	fmt.Println(u)

}
