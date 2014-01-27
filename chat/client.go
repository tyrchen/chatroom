package chat

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	conn     net.Conn
	incoming Message
	outgoing Message
	reader   *bufio.Reader
	writer   *bufio.Writer
	quiting  chan net.Conn
}

func (self *Client) GetIncoming() string {
	return <-self.incoming
}

func (self *Client) PutOutgoing(message string) {
	self.outgoing <- message
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

func (self *Client) Listen() {
	go self.Read()
	go self.Write()
}

func (self *Client) Read() {
	for {
		if line, _, err := self.reader.ReadLine(); err == nil {
			self.incoming <- string(line)
		} else {
			log.Printf("Read error: %s\n", err)
			self.quiting <- self.conn
			return
		}
	}

}

func (self *Client) Write() {
	for data := range self.outgoing {
		if _, err := self.writer.WriteString(data + "\n"); err != nil {
			self.quiting <- self.conn
			return
		}

		if err := self.writer.Flush(); err != nil {
			log.Printf("Write error: %s\n", err)
			self.quiting <- self.conn
			return
		}
	}

}

func (self *Client) Close() {
	self.conn.Close()
}
