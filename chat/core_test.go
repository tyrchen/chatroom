package chat

import (
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

type (
	Clients []*Client
)

const (
	CONNSTR  = ":5555"
	EXPECTED = "Hello world"
)

func startServer() (server *Server) {
	server = CreateServer()
	log.Printf("Server %p created\n", server)
	go server.Start(CONNSTR)
	return
}

func startClient() (client *Client) {
	conn, err := net.Dial("tcp", CONNSTR)

	if err != nil {
		log.Fatal(err)
	}

	client = CreateClient(conn)

	return
}

func startClients(N int) (clients Clients) {
	clients = make(Clients, N)
	for i := 0; i < N; i++ {
		clients[i] = startClient()
	}
	return
}

func TestBroadcast(t *testing.T) {
	server := startServer()
	N := MAXCLIENTS
	tokens := make(chan int, N)

	time.Sleep(50 * time.Microsecond)

	clients := startClients(N)
	clients[0].PutOutgoing(EXPECTED + "\n")

	time.Sleep(50 * time.Millisecond)
	for i := 0; i < N; i++ {
		msg := <-clients[i].incoming
		tokens <- 0
		if strings.Contains(msg, EXPECTED) {
			t.Logf("%d: %s\n", i, msg)
		} else {
			t.Errorf("Message: %s, expected %s\n", msg, EXPECTED)
		}
	}

	go func() {
		for i := 0; i < N; i++ {
			<-tokens
		}
		server.Stop()
	}()

}

func TestJoinLeave(t *testing.T) {
	server := startServer()
	time.Sleep(50 * time.Microsecond)
	N := MAXCLIENTS + 1
	M := 10
	tokens := make(chan int, N)

	clients := startClients(N)
	time.Sleep(50 * time.Millisecond)
	if len(server.clients) != MAXCLIENTS {
		t.Errorf("Clients: %d, expected %d", len(server.clients), MAXCLIENTS)
	}

	clients[0].Close()
	time.Sleep(50 * time.Millisecond)

	if len(server.clients) != MAXCLIENTS {
		t.Errorf("Clients: %d, expected %d", len(server.clients), MAXCLIENTS)
	}

	clients[1].PutOutgoing(EXPECTED + "\n")

	for i := 1; i < M; i++ {
		log.Printf("Close client %p\n", clients[i])
		clients[i].Close()
	}

	for i := N + 1; i < N-M; i++ {
		msg := <-clients[i].incoming
		tokens <- 0
		if strings.Contains(msg, EXPECTED) {
			t.Logf("%d: %s\n", i, msg)
		} else {
			t.Errorf("Message: %s, expected %s\n", msg, EXPECTED)
		}
	}

	go func() {
		for i := 0; i < N; i++ {
			<-tokens
		}
		server.Stop()
	}()
}

func TestChangeName(t *testing.T) {
	server := startServer()
	time.Sleep(50 * time.Microsecond)
	N := 2

	clients := startClients(N)
	time.Sleep(50 * time.Millisecond)

	clients[0].PutOutgoing(":name Tyr\n")
}
