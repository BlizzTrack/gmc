package gmc

import (
	"bufio"
	"fmt"
	"github.com/blizztrack/gmc/commands"
	"io"
	"log"
	"net"
	"strings"
)

type Conn struct {
	Conn net.Conn
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

func newConn(rwc net.Conn) (c *Conn) {
	c = new(Conn)
	c.Conn = rwc
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

		go HandleClient(newConn(rw))
	}
}

func HandleClient(Conn *Conn) {
	defer Conn.Conn.Close()
	for {
		netData, err := Conn.ReadLine()
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
			set := &commands.SetCommand{}
			res = set.Handle(payload, Conn)
			break
		case "get":
			get := &commands.GetCommand{}
			res = get.Handle(payload)
			break
		case "delete":
			del := &commands.DeleteCommand{}
			res = del.Handle(payload)
			break
		case "flush_all":
			flush := &commands.FlushAllCommand{}
			res = flush.Handle(payload)
			break
		case "version":
			res = MessageResponse{Message: fmt.Sprintf(StatusVersion, Version)}
		case "touch":
			touch := &commands.TouchCommand{}
			res = touch.Handle(payload)
		case "quit":
			return
		}

		if res != nil {
			if err := res.WriteResponse(Conn); err != nil {
				log.Printf("write to client failed %+v", err)
				return
			}
		}

		if err := Conn.RWC.Flush(); err != nil {
			log.Printf("failed to flush buffer: %v", err)
			return
		}
	}
}

func (c *Conn) ReadLine() (line []byte, err error) {
	line, _, err = c.RWC.ReadLine()
	return
}

func (c *Conn) Read(p []byte) (n int, err error) {
	return io.ReadFull(c.RWC, p)
}

