package main

func main() {
	for i := 0; i < 5; i++ {
		defer func(i int) {
			println(i)
		}(i)
	}
}
