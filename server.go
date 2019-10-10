package gmc

import (
	"bufio"
	"fmt"
	"github.com/blizztrack/gmc/commands"
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type conn struct {
	conn net.Conn
	rwc  *bufio.ReadWriter
}

const Version = "0.0.1"

func NewServer(address string) error {
	l, e := net.Listen("tcp", address)
	if e != nil {
		return e
	}

	go func() {
		for range time.Tick(time.Second) {
			lru.Clean()
		}
	}()

	return serve(l)
}

func newconn(rwc net.Conn) (c *conn) {
	c = new(conn)
	c.conn = rwc
	c.rwc = bufio.NewReadWriter(bufio.NewReader(rwc), bufio.NewWriter(rwc))
	return c
}

func serve(l net.Listener) error {
	defer l.Close()
	for {
		rw, e := l.Accept()
		if e != nil {
			return e
		}

		go HandleClient(newconn(rw))
	}
}

func HandleClient(conn *conn) {
	defer conn.conn.Close()

	for {
		inLine, err := conn.ReadLine()
		if err != nil || len(inLine) == 0 {
			if err == io.EOF {
				return
			}
			_, _ = conn.rwc.WriteString(err.Error())
			_ = conn.rwc.Flush()
			return
		}

		temp := strings.TrimSpace(string(inLine))
		log.Printf("client sent command: %s", temp)

		line := strings.Split(temp, " ")
		var res responses.Response
		command := line[0]
		payload := line[1:]

		switch strings.ToLower(command) {
		case "set":
			set := &commands.SetCommand{}
			res = set.Handle(payload, conn)
		case "get", "gat":
			get := &commands.GetCommand{}
			res = get.Handle(payload)
		case "gets", "gats":
			gets := &commands.GetsCommand{}
			res = gets.Handle(payload)
		case "delete":
			del := &commands.DeleteCommand{}
			res = del.Handle(payload)
		case "flush_all":
			flush := &commands.FlushAllCommand{}
			res = flush.Handle(payload)
		case "version":
			res = responses.MessageResponse{Message: fmt.Sprintf(responses.StatusVersion, Version)}
		case "touch":
			touch := &commands.TouchCommand{}
			res = touch.Handle(payload)
		case "quit":
			return
		}

		if res != nil {
			if err := res.WriteResponse(conn.rwc); err != nil {
				log.Printf("write to client failed %+v", err)
				return
			}
		}

		if err := conn.rwc.Flush(); err != nil {
			log.Printf("failed to flush buffer: %v", err)
			return
		}
	}
}

func (c *conn) ReadLine() (line []byte, err error) {
	line, _, err = c.rwc.ReadLine()
	return
}

func (c *conn) Read(p []byte) (n int, err error) {
	return io.ReadFull(c.rwc, p)
}
