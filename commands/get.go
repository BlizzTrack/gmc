package commands

import (
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
)

type GetCommand struct{}

func (get *GetCommand) Handle(payload []string) responses.Response {
	if len(payload) < 1 || len(payload) > 1 {
		return responses.InvalidParamLengthResponse{}
	}

	item, ok := lru.LRU.Get(payload[0])
	if !ok {
		return responses.MessageResponse{Message: responses.StatusEnd}
	}

	if item.IsExpired() {
		lru.LRU.Remove(item.Key)

		return responses.MessageResponse{Message: responses.StatusEnd}
	}

	return responses.ItemResponse{Item: item}
}