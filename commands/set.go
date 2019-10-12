package commands

import (
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
	Length  int
	NoReply bool
	Payload []byte
}

func (set *SetCommand) Handle(payload []string, client io.Reader) responses.Response {
	if len(payload) < 4 || len(payload) > 5 {
		_, err := readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}

		return responses.InvalidParamLengthResponse{}
	}

	err := set.read(payload)
	if err != nil {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}

		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, err)}
	}

	n := make([]byte, set.Length)
	size, err := io.ReadFull(client, n)
	if err != nil {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}

		log.Printf("failed to read payload: %+v", err)
		if set.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	if size != set.Length {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}

		log.Printf("failed to read payload wanted size %d got %d", set.Length, len(n))
		if set.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	item := &lru.Item{
		Key:   set.Key,
		Flags: set.Flags,
	}
	item.SetExpires(set.ExpTime)
	item.Value = make([]byte, set.Length)
	copy(item.Value, n)
	lru.LRU.Add(item.Key, item)

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
	set.Length, err = strconv.Atoi(payload[3])
	if err != nil {
		return err
	}

	set.NoReply = len(payload) == 5 && isNoReply(payload[4])

	return nil
}
