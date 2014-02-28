package main

import (
	. "./chat"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(-1)
	}

	server := CreateServer()
	fmt.Printf("Running on %s\n", os.Args[1])
	server.Start(os.Args[1])

}
