package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	MAXCLIENTS = 50
)

type Message chan string
type Queue []chan int
type ClientTable map[net.Conn]*Client

type Client struct {
	conn     net.Conn
	incoming Message
	outgoing Message
	reader   *bufio.Reader
	writer   *bufio.Writer
	quiting  chan net.Conn
}

func (self *Client) Listen() {
	//fmt.Printf("Client %v is listening.\n", self.conn)
	go self.Read()
	go self.Write()
}

func (self *Client) Read() {
	for {
		if line, _, err := self.reader.ReadLine(); err == nil {
			//fmt.Println("read:", string(line))
			self.incoming <- string(line)
		} else {
			fmt.Printf("Error: %s\n", err)
			self.quiting <- self.conn
			return
		}
	}

}

func (self *Client) Write() {
	for data := range self.outgoing {
		//fmt.Println("write:", data)
		if _, err := self.writer.WriteString(data + "\n"); err != nil {
			self.quiting <- self.conn
			return
		}

		if err := self.writer.Flush(); err != nil {
			fmt.Printf("Error: %s\n", err)
			self.quiting <- self.conn
			return
		}
	}

}

func (self *Client) GetIncoming() string {
	return <-self.incoming
}

func (self *Client) PutOutgoing(message string) {
	self.outgoing <- message
}

type Server struct {
	clients  ClientTable
	queue    Queue
	pending  chan net.Conn
	quiting  chan net.Conn
	incoming Message
	outgoing Message
}

func CreateClient(conn net.Conn) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	client := &Client{
		conn:     conn,
		incoming: make(Message),
		outgoing: make(Message),
		quiting:  make(chan net.Conn),
		reader:   reader,
		writer:   writer,
	}
	client.Listen()
	return client
}

func CreateServer() Server {
	server := Server{
		clients:  make(ClientTable),
		queue:    make(Queue, MAXCLIENTS),
		pending:  make(chan net.Conn),
		quiting:  make(chan net.Conn),
		incoming: make(Message),
		outgoing: make(Message),
	}
	server.Listen()
	return server
}

func (self *Server) Start(connString string) {
	listener, _ := net.Listen("tcp", connString)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("A new connection %v kicks\n", conn)

		self.pending <- conn
	}
}

func (self *Server) Join(conn net.Conn) {
	client := CreateClient(conn)
	self.clients[conn] = client

	go func() {
		for {
			msg := <-client.incoming
			fmt.Printf("Got message: %s\n", msg)
			self.incoming <- msg
		}
	}()

	go func() {
		for {
			conn := <-client.quiting
			fmt.Printf("Conn %v is quiting\n", conn)
			self.quiting <- conn
		}
	}()
}

func (self *Server) Leave(conn net.Conn) {
	fmt.Printf("Client %v is leaving\n", conn)
	conn.Close()
	delete(self.clients, conn)
}

func (self *Server) Broadcast(message string) {
	fmt.Printf("Broadcasting message: %s\n", message)
	for _, client := range self.clients {
		client.outgoing <- message
	}
}

func (self *Server) Listen() {
	go func() {
		for {
			select {
			case message := <-self.incoming:
				self.Broadcast(message)
			case conn := <-self.pending:
				self.Join(conn)
			case conn := <-self.quiting:
				self.Leave(conn)
			}
		}
	}()
}
