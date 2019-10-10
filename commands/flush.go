package commands

import (
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
)

type FlushAllCommand struct{}

func (flush *FlushAllCommand) Handle(payload []string) responses.Response {
	lru.Flush()

	return responses.MessageResponse{Message: responses.StatusOK}
}
