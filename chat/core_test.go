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

func startServerClients(N int) (server *Server, clients Clients) {
	server = startServer()

	time.Sleep(50 * time.Microsecond)

	clients = startClients(N)

	time.Sleep(50 * time.Millisecond)

	return
}

func verifyAndStopServer(t *testing.T, server *Server, clients Clients, expected string) {
	N := len(clients)
	tokens := make(chan int, N)

	for i := 0; i < N; i++ {
		msg := <-clients[i].incoming
		tokens <- 0
		if strings.Contains(msg, expected) {
			t.Logf("%d: %s\n", i, msg)
		} else {
			t.Errorf("Message: %s, expected %s\n", msg, expected)
		}
	}

	go func() {
		for i := 0; i < N; i++ {
			<-tokens
		}
		server.Stop()
	}()
}

func TestBroadcast(t *testing.T) {
	N := MAXCLIENTS
	server, clients := startServerClients(N)

	clients[0].PutOutgoing(EXPECTED + "\n")

	verifyAndStopServer(t, server, clients, EXPECTED)
}

func TestJoinLeave(t *testing.T) {
	N := MAXCLIENTS + 1
	M := 10

	server, clients := startServerClients(N)

	if len(server.clients) != MAXCLIENTS {
		t.Errorf("Clients: %d, expected %d", len(server.clients), MAXCLIENTS)
	}

	clients[0].Close()
	time.Sleep(50 * time.Millisecond)

	if len(server.clients) != MAXCLIENTS {
		t.Errorf("Clients: %d, expected %d", len(server.clients), MAXCLIENTS)
	}

	clients[M+1].PutOutgoing(EXPECTED + "\n")

	for i := 1; i < M; i++ {
		log.Printf("Close client %p\n", clients[i])
		clients[i].Close()
	}

	verifyAndStopServer(t, server, clients[M+1:], EXPECTED)

}

func TestChangeName(t *testing.T) {
	N := 2

	server, clients := startServerClients(N)

	clients[0].PutOutgoing(":name Tyr\n")

	verifyAndStopServer(t, server, clients, "changed its name to Tyr")
}
