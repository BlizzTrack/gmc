package commands

import (
	"fmt"
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
	"strconv"
)

type TouchCommand struct{}

func (get *TouchCommand) Handle(payload []string) responses.Response {
	if len(payload) < 2 || len(payload) > 3 {
		return responses.InvalidParamLengthResponse{}
	}

	item, ok := lru.LRU.Get(payload[0])
	if ok {
		return responses.MessageResponse{Message: responses.StatusNotFound}
	}
	if item.IsExpired() {
		lru.LRU.Remove(item.Key)

		return responses.MessageResponse{Message: responses.StatusNotFound}
	}
	ExpTime, err := strconv.ParseInt(payload[1], 10, 64)
	if err != nil {
		return responses.MessageResponse{Message: fmt.Sprintf(responses.StatusClientError, err)}
	}
	item.SetExpires(ExpTime)

	// Update the item in the cache so we move it to the top
	lru.LRU.Add(item.Key, item)

	if len(payload) == 3 && isNoReply(payload[2]) {
		return nil
	}

	return responses.MessageResponse{Message: responses.StatusTouched}
}