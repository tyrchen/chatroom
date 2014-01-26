package main

import (
	"bufio"
	. "chatroom/chat"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5555")
	defer conn.Close()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	if err != nil {
		log.Fatal(err)
	}

	client := CreateClient(conn)

	go func() {
		for {
			out.WriteString(client.GetIncoming() + "\n")
			out.Flush()
		}
	}()
	for {
		line, _, _ := in.ReadLine()
		client.PutOutgoing(string(line))
	}

}
