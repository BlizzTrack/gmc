package commands

import (
	"fmt"
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
	"io"
	"log"
	"strconv"
)

type AddCommand struct {
	Key     string
	Flags   int
	ExpTime int64
	Length  int
	NoReply bool
	Payload []byte
}

func (add *AddCommand) Handle(payload []string, client io.Reader) responses.Response {
	if len(payload) < 4 || len(payload) > 5 {
		_, err := readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}

		return responses.InvalidParamLengthResponse{}
	}

	err := add.read(payload)
	if err != nil {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}

		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, err)}
	}

	if lru.Has(add.Key) {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}
		return &responses.MessageResponse{Message: responses.StatusNotStored}
	}

	n := make([]byte, add.Length)
	size, err := io.ReadFull(client, n)

	if err != nil {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}
		log.Printf("failed to read payload: %+v", err)
		if add.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	if size != add.Length {
		_, err = readLine(client)
		if err != nil {
			log.Printf("failed to read line to the end")
		}

		log.Printf("failed to read payload wanted size %d got %d", add.Length, len(n))
		if add.NoReply {
			return nil
		}
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, "bad data chunk")}
	}

	item := &lru.Item{
		Key:   add.Key,
		Flags: add.Flags,
	}
	item.SetExpires(add.ExpTime)
	item.Value = make([]byte, add.Length)
	copy(item.Value, n)

	lru.Set(item)

	if add.NoReply {
		return nil
	}

	return responses.MessageResponse{Message: responses.StatusStored}
}

func (add *AddCommand) read(payload []string) error {
	var err error
	add.Key = payload[0]
	add.Flags, err = strconv.Atoi(payload[1])
	if err != nil {
		return err
	}
	add.ExpTime, err = strconv.ParseInt(payload[2], 10, 64)
	if err != nil {
		return err
	}
	add.Length, err = strconv.Atoi(payload[3])
	if err != nil {
		return err
	}

	add.NoReply = len(payload) == 5 && isNoReply(payload[4])

	return nil
}
