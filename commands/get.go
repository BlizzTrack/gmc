package commands

import (
	"github.com/blizztrack/gmc"
	"github.com/blizztrack/gmc/lru"
)

type GetCommand struct{}

func (get *GetCommand) Handle(payload []string) gmc.Response {
	if len(payload) < 1 || len(payload) > 1 {
		return gmc.InvalidParamLengthResponse{}
	}

	item, err := lru.Get(payload[0])
	if err != nil {
		return gmc.MessageResponse{Message: gmc.StatusNotFound}
	}

	if item.IsExpired() {
		lru.Delete(item.Key)

		return gmc.MessageResponse{Message: gmc.StatusNotFound}
	}

	return gmc.ItemResponse{Item: item}
}