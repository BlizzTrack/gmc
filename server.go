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

func newConn(rwc net.Conn) (c *conn) {
	c = new(conn)
	c.conn = rwc
	c.rwc = bufio.NewReadWriter(bufio.NewReaderSize(rwc, 1048576), bufio.NewWriter(rwc))
	return c
}

func serve(l net.Listener) error {
	defer l.Close()
	for {
		rw, e := l.Accept()
		if e != nil {
			return e
		}

		log.Printf("handling client: %+v", rw.RemoteAddr())

		conn := newConn(rw)
		go conn.serve()
	}
}

func (c *conn) serve() {
	defer c.conn.Close()

	for {
		err := c.handle()
		if err != nil {
			if err == io.EOF {
				return
			}
			_, _ = c.rwc.WriteString(err.Error())
		}
		c.end()
	}
}

func (c *conn) handle() error {
	inLine, err := c.ReadLine()
	if err != nil || len(inLine) == 0 {
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			_, _ = c.rwc.WriteString(err.Error())
		}

		return responses.Error
	}

	temp := strings.TrimSpace(string(inLine))
	if temp == "" {
		return nil
	}

	log.Printf("client sent command: %s", temp)

	line := strings.Split(temp, " ")
	var res responses.Response
	command := line[0]
	payload := line[1:]

	switch strings.ToLower(command) {
	case "set":
		set := &commands.SetCommand{}
		res = set.Handle(payload, c.rwc)
	case "add":
		set := &commands.AddCommand{}
		res = set.Handle(payload, c.rwc)
	case "replace":
		set := &commands.ReplaceCommand{}
		res = set.Handle(payload, c.rwc)
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
	case "incr":
		touch := &commands.IncrCommand{}
		res = touch.Handle(payload)
	case "decr":
		touch := &commands.DecrCommand{}
		res = touch.Handle(payload)
	case "has":
		touch := &commands.HasCommand{}
		res = touch.Handle(payload)
	case "quit":
		if len(line) == 4 {
			return io.EOF
		}
		return responses.Error
	}

	if res != nil {
		if err := res.WriteResponse(c.rwc); err != nil {
			return fmt.Errorf("write to client failed %+v", err)
		}
	}

	return nil
}

func (c *conn) end() {
	c.rwc.Flush()
}

func (c *conn) ReadLine() (line []byte, err error) {
	line, err = c.rwc.ReadBytes('\n')
	return
}

func (c *conn) Read(p []byte) (n int, err error) {
	return io.ReadFull(c.rwc, p)
}
