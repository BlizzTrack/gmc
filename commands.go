package gmc

import (
	"bytes"
	"fmt"
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

func (set *SetCommand) Handle(payload []string, client *conn) Response {
	set.read(payload)

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
		return MessageResponse{Message: fmt.Sprintf(StatusClientError, "bad data chunk")}
	}

	if int64(len(n)) != set.Length {
		log.Printf("failed to read payload wanted size %d got %d", set.Length, len(n))
		if set.NoReply {
			return nil
		}
		return MessageResponse{Message: fmt.Sprintf(StatusClientError, "bad data chunk")}
	}

	item.Value = make([]byte, len(n))
	copy(item.Value, n)

	lru.Set(item)
	log.Println(item.String())

	if set.NoReply {
		return nil
	}
	return MessageResponse{Message: StatusStored}
}

func (set *SetCommand) read(pieces []string) {
	set.Key = pieces[0]
	set.Flags, _ = strconv.Atoi(pieces[1])
	set.ExpTime, _ = strconv.ParseInt(pieces[2], 10, 64)
	set.Length, _ = strconv.ParseInt(pieces[3], 10, 64)

	set.NoReply = len(pieces) == 5 && isNoReply(pieces[4])
}

type GetCommand struct{}

func (get *GetCommand) Handle(payload []string) Response {
	item, err := lru.Get(payload[0])
	if err != nil {
		return MessageResponse{Message: StatusNotFound}
	}

	if item.IsExpired() {
		lru.Delete(item.Key)

		return MessageResponse{Message: StatusNotFound}
	}

	return ItemResponse{Item: item}
}

type DeleteCommand struct{}

func (del *DeleteCommand) Handle(payload []string) Response {
	lru.Delete(payload[0])

	if len(payload) == 2 && isNoReply(payload[1]) {
		return nil
	}

	return MessageResponse{Message: StatusDeleted}
}

type FlushAllCommand struct{}

func (flush *FlushAllCommand) Handle(payload []string) Response {
	lru.Flush()

	return MessageResponse{Message: StatusOK}
}

func isNoReply(payload string) bool {
	return bytes.Equal([]byte(payload), []byte("noreply"))
}