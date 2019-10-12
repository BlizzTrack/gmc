package commands

import (
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
)

type GetsCommand struct{}

func (gets *GetsCommand) Handle(payload []string) responses.Response {
	if len(payload) < 1 {
		return responses.InvalidParamLengthResponse{}
	}

	items := make([]*lru.Item, 0)

	for _, key := range payload {
		if item, ok := lru.LRU.Get(key); ok {
			items = append(items, item)
		}
	}

	return responses.MultiItemResponse{Items: items}
}
