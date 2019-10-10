package commands

import (
	"github.com/blizztrack/gmc"
	"github.com/blizztrack/gmc/lru"
)

type FlushAllCommand struct{}

func (flush *FlushAllCommand) Handle(payload []string) gmc.Response {
	lru.Flush()

	return gmc.MessageResponse{Message: gmc.StatusOK}
}
