package commands

import (
	"github.com/blizztrack/gmc"
	"github.com/blizztrack/gmc/lru"
)

type DeleteCommand struct{}

func (del *DeleteCommand) Handle(payload []string) gmc.Response {
	lru.Delete(payload[0])

	if len(payload) == 2 && isNoReply(payload[1]) {
		return nil
	}

	return gmc.MessageResponse{Message: gmc.StatusDeleted}
}