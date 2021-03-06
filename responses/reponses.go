package responses

import (
	"fmt"
	"github.com/blizztrack/gmc/lru"
	"io"
)

type Response interface {
	WriteResponse(io.Writer) error
}

type MessageResponse struct {
	Message string
}

func (m MessageResponse) WriteResponse(out io.Writer) error {
	_, err := out.Write([]byte(m.Message))
	return err
}

type InvalidParamLengthResponse struct{}

func (m InvalidParamLengthResponse) WriteResponse(out io.Writer) error {
	_, err := fmt.Fprintf(out, StatusClientError, "invalid number of parameters sent")
	return err
}

type ItemResponse struct {
	Item *lru.Item
}

func (r ItemResponse) WriteResponse(out io.Writer) error {
	_, _ = fmt.Fprintf(out, StatusValue, r.Item.Key, r.Item.Flags, len(r.Item.Value))
	_, _ = out.Write(r.Item.Value)
	_, _ = out.Write([]byte("\r\n"))
	_, _ = out.Write([]byte(StatusEnd))

	return nil
}

type MultiItemResponse struct {
	Items []*lru.Item
}

func (r MultiItemResponse) WriteResponse(out io.Writer) error {
	for _, item := range r.Items {
		if item != nil {
			_, _ = fmt.Fprintf(out, StatusValue, item.Key, item.Flags, len(item.Value))
			_, _ = out.Write(item.Value)
			_, _ = out.Write([]byte("\r\n"))
		}
	}
	_, _ = out.Write([]byte(StatusEnd))

	return nil
}
