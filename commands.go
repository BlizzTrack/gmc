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

func (set *SetCommand) handle(payload []string, client *conn) Response {
	if len(payload) < 4 || len(payload) > 5 {
		return InvalidParamLengthResponse{}
	}

	err := set.read(payload)
	if err != nil {
		return MessageResponse{fmt.Sprintf(StatusClientError, err)}
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

type GetCommand struct{}

func (get *GetCommand) handle(payload []string) Response {
	if len(payload) < 1 || len(payload) > 1 {
		return InvalidParamLengthResponse{}
	}

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

type TouchCommand struct{}

func (get *TouchCommand) handle(payload []string) Response {
	if len(payload) < 2 || len(payload) > 3 {
		return InvalidParamLengthResponse{}
	}

	item, err := lru.Get(payload[0])
	if err != nil {
		return MessageResponse{Message: StatusNotFound}
	}
	if item.IsExpired() {
		lru.Delete(item.Key)

		return MessageResponse{Message: StatusNotFound}
	}
	ExpTime, err := strconv.ParseInt(payload[1], 10, 64)
	if err != nil {
		return MessageResponse{fmt.Sprintf(StatusClientError, err)}
	}

	item.SetExpires(ExpTime)

	if len(payload) == 3 && isNoReply(payload[2]) {
		return nil
	}

	return MessageResponse{Message: StatusTouched}
}

type DeleteCommand struct{}

func (del *DeleteCommand) handle(payload []string) Response {
	lru.Delete(payload[0])

	if len(payload) == 2 && isNoReply(payload[1]) {
		return nil
	}

	return MessageResponse{Message: StatusDeleted}
}

type FlushAllCommand struct{}

func (flush *FlushAllCommand) handle(payload []string) Response {
	lru.Flush()

	return MessageResponse{Message: StatusOK}
}

func isNoReply(payload string) bool {
	return bytes.Equal([]byte(payload), []byte("noreply"))
}
