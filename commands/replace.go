package commands

import (
	"fmt"
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
	"io"
	"log"
	"strconv"
)

type ReplaceCommand struct {
	Key     string
	Flags   int
	ExpTime int64
	Length  int
	NoReply bool
	Payload []byte
}

func (replace *ReplaceCommand) Handle(payload []string, client io.Reader) responses.Response {
	if len(payload) < 4 || len(payload) > 5 {
		_, err := readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}
		return responses.InvalidParamLengthResponse{}
	}

	err := replace.read(payload)
	if err != nil {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, err)}
	}

	if !lru.LRU.Has(replace.Key) {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}
		return &responses.MessageResponse{Message: responses.StatusNotStored}
	}

	n := make([]byte, replace.Length)
	size, err := io.ReadFull(client, n)

	if err != nil {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}
		log.Printf("failed to read payload: %+v", err)
		if replace.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	if size != replace.Length {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}
		log.Printf("failed to read payload wanted size %d got %d", replace.Length, len(n))
		if replace.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	item := &lru.Item{
		Key:   replace.Key,
		Flags: replace.Flags,
	}
	item.SetExpires(replace.ExpTime)
	item.Value = make([]byte, replace.Length)
	copy(item.Value, n)

	lru.LRU.Add(item.Key, item)

	if replace.NoReply {
		return nil
	}
	return responses.MessageResponse{Message: responses.StatusStored}
}

func (replace *ReplaceCommand) read(payload []string) error {
	var err error
	replace.Key = payload[0]
	replace.Flags, err = strconv.Atoi(payload[1])
	if err != nil {
		return err
	}
	replace.ExpTime, err = strconv.ParseInt(payload[2], 10, 64)
	if err != nil {
		return err
	}
	replace.Length, err = strconv.Atoi(payload[3])
	if err != nil {
		return err
	}

	replace.NoReply = len(payload) == 5 && isNoReply(payload[4])

	return nil
}
