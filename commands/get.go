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

	item, err := lru.Get(payload[0])
	if err != nil {
		return responses.MessageResponse{Message: responses.StatusNotFound}
	}

	if item.IsExpired() {
		lru.Delete(item.Key)

		return responses.MessageResponse{Message: responses.StatusNotFound}
	}

	return responses.ItemResponse{Item: item}
}