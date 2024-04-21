package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		fmt.Fprintf(conn, text+"\n")
		reply, err := bufio.NewReader(conn).ReadString('%')
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Server:", reply)
	}
}
