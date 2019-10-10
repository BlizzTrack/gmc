package commands

import (
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
)

type DeleteCommand struct{}

func (del *DeleteCommand) Handle(payload []string) responses.Response {
	lru.Delete(payload[0])

	if len(payload) == 2 && isNoReply(payload[1]) {
		return nil
	}

	return responses.MessageResponse{Message: responses.StatusDeleted}
}