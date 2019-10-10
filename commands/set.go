package commands

import (
	"fmt"
	"github.com/blizztrack/gmc"
	"github.com/blizztrack/gmc/lru"
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

func (set *SetCommand) Handle(payload []string, client *gmc.Conn) gmc.Response {
	if len(payload) < 4 || len(payload) > 5 {
		return gmc.InvalidParamLengthResponse{}
	}

	err := set.read(payload)
	if err != nil {
		return gmc.MessageResponse{Message: fmt.Sprintf(gmc.StatusClientError, err)}
	}

	item := &lru.Item{
		Key:   set.Key,
		Flags: set.Flags,
	}
	item.SetExpires(set.ExpTime)
	log.Println(item.String())

	n, err := client.ReadLine()

	if err != nil {
		log.Printf("failed to read payload: %+v", err)
		if set.NoReply {
			return nil
		}
		return gmc.MessageResponse{Message: fmt.Sprintf(gmc.StatusClientError, "bad data chunk")}
	}

	if int64(len(n)) != set.Length {
		log.Printf("failed to read payload wanted size %d got %d", set.Length, len(n))
		if set.NoReply {
			return nil
		}
		return gmc.MessageResponse{Message: fmt.Sprintf(gmc.StatusClientError, "bad data chunk")}
	}

	item.Value = make([]byte, len(n))
	copy(item.Value, n)

	lru.Set(item)
	log.Println(item.String())

	if set.NoReply {
		return nil
	}
	return gmc.MessageResponse{Message: gmc.StatusStored}
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