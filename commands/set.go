package commands

import (
	"bufio"
	"fmt"
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
	"io"
	"log"
	"strconv"
)

type SetCommand struct {
	Key     string
	Flags   int
	ExpTime int64
	Length  int64
	NoReply bool
	Payload []byte
}

func (set *SetCommand) Handle(payload []string, client io.Reader) responses.Response {
	if len(payload) < 4 || len(payload) > 5 {
		return responses.InvalidParamLengthResponse{}
	}

	err := set.read(payload)
	if err != nil {
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, err)}
	}

	item := &lru.Item{
		Key:   set.Key,
		Flags: set.Flags,
	}
	item.SetExpires(set.ExpTime)
	n, _, err := bufio.NewReader(client).ReadLine()

	if err != nil {
		log.Printf("failed to read payload: %+v", err)
		if set.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	if int64(len(n)) != set.Length {
		log.Printf("failed to read payload wanted size %d got %d", set.Length, len(n))
		if set.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	item.Value = make([]byte, len(n))
	copy(item.Value, n)

	lru.Set(item)

	if set.NoReply {
		return nil
	}
	return responses.MessageResponse{Message: responses.StatusStored}
}

func (set *SetCommand) read(payload []string) error {
	var err error
	set.Key = payload[0]
	set.Flags, err = strconv.Atoi(payload[1])
	if err != nil {
		return err
	}
	set.ExpTime, err = strconv.ParseInt(payload[2], 10, 64)
	if err != nil {
		return err
	}
	set.Length, err = strconv.ParseInt(payload[3], 10, 64)
	if err != nil {
		return err
	}

	set.NoReply = len(payload) == 5 && isNoReply(payload[4])

	return nil
}
