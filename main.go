package main

import (
	"fmt"

	"halo0201.live/code/pattern/factory"
)

func main() {
	var h factory.Human
	h = factory.NewAsia()

	fmt.Println(h.Language())
}
