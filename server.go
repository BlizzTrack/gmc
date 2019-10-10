package gmc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type conn struct {
	conn net.Conn
	RWC  *bufio.ReadWriter
}

const Version = "0.0.1"

func NewServer(address string) error {
	l, e := net.Listen("tcp", address)
	if e != nil {
		return e
	}
	return serve(l)
}

func newConn(rwc net.Conn) (c *conn) {
	c = new(conn)
	c.conn = rwc
	c.RWC = bufio.NewReadWriter(bufio.NewReaderSize(rwc, 1048576), bufio.NewWriter(rwc))
	return c
}

func serve(l net.Listener) error {
	defer l.Close()
	for {
		rw, e := l.Accept()
		if e != nil {
			return e
		}

		go handleClient(newConn(rw))
	}
}

func handleClient(conn *conn) {
	defer conn.conn.Close()
	for {
		netData, err := conn.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		log.Printf("client sent command: %s", temp)

		line := strings.Split(temp, " ")
		var res Response
		command := line[0]
		payload := line[1:]

		switch strings.ToLower(command) {
		case "set":
			set := &SetCommand{}
			res = set.handle(payload, conn)
			break
		case "get":
			get := &GetCommand{}
			res = get.handle(payload)
			break
		case "delete":
			del := &DeleteCommand{}
			res = del.handle(payload)
			break
		case "flush_all":
			flush := &FlushAllCommand{}
			res = flush.handle(payload)
			break
		case "version":
			res = MessageResponse{Message: fmt.Sprintf(StatusVersion, Version)}
		case "touch":
			touch := &TouchCommand{}
			res = touch.handle(payload)
		case "quit":
			return
		}

		if res != nil {
			if err := res.WriteResponse(conn); err != nil {
				log.Printf("write to client failed %+v", err)
				return
			}
		}

		if err := conn.RWC.Flush(); err != nil {
			log.Printf("failed to flush buffer: %v", err)
			return
		}
	}
}

func (c *conn) ReadLine() (line []byte, err error) {
	line, _, err = c.RWC.ReadLine()
	return
}

func (c *conn) Read(p []byte) (n int, err error) {
	return io.ReadFull(c.RWC, p)
}

