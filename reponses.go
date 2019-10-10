package gmc

import (
	"fmt"
	"github.com/blizztrack/gmc/lru"
	"log"
)

type Response interface {
	WriteResponse(*conn) error
}

type MessageResponse struct {
	Message string
}

func (m MessageResponse) WriteResponse(writer *conn) error {
	n, err := writer.RWC.WriteString(m.Message)
	log.Printf("wrote %d bytes to client", n)
	return err
}

type InvalidParamLengthResponse struct{}

func (m InvalidParamLengthResponse) WriteResponse(writer *conn) error {
	n, err := writer.RWC.WriteString(fmt.Sprintf(StatusClientError, "invalid number of parameters sent"))
	log.Printf("wrote %d bytes to client", n)
	return err
}

type ItemResponse struct {
	Item *lru.Item
}

func (r ItemResponse) WriteResponse(out *conn) error {
	_, _ = fmt.Fprintf(out.RWC, StatusValue, r.Item.Key, r.Item.Flags, len(r.Item.Value))
	_, _ = out.RWC.Write(r.Item.Value)
	_, _ = out.RWC.Write([]byte("\r\n"))
	_, _ = out.RWC.Write([]byte(StatusEnd))

	return nil
}
