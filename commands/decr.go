package commands

import (
	"fmt"
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
	"strconv"
)

type DecrCommand struct{}

func (*DecrCommand) Handle(payload []string) responses.Response {
	if len(payload) < 2 || len(payload) > 3 {
		return responses.InvalidParamLengthResponse{}
	}

	key := payload[0]

	amount, err := strconv.Atoi(payload[1])
	if err != nil || amount < 0{
		return responses.MessageResponse{Message: responses.StatusError}
	}

	value, ok := lru.LRU.Get(key)
	if !ok {
		return responses.MessageResponse{Message: responses.StatusNotFound}
	}

	current, err := strconv.Atoi(string(value.Value))
	if err != nil {
		return responses.MessageResponse{Message: responses.StatusError}
	}

	if current == 0 {
		return responses.MessageResponse{Message: fmt.Sprintf("%d\r\n", current)}
	}

	if current - amount < 0   {
		return responses.MessageResponse{Message: fmt.Sprintf("%d\r\n", current)}
	}

	current = current - amount
	value.Value = []byte(strconv.Itoa(current))

	lru.LRU.Add(value.Key, value)

	if len(payload) == 3 && isNoReply(payload[2]) {
		return nil
	}

	return responses.MessageResponse{Message: fmt.Sprintf("%d\r\n", current)}
}
