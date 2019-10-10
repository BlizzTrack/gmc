package commands

import (
	"github.com/blizztrack/gmc/lru"
	"github.com/blizztrack/gmc/responses"
)

type HasCommand struct{}

func (*HasCommand) Handle(payload []string) responses.Response {
	if len(payload) < 1 || len(payload) > 1 {
		return responses.InvalidParamLengthResponse{}
	}

	if lru.Has(payload[0]) {
		return responses.MessageResponse{Message: responses.StatusExists}

	}

	return responses.MessageResponse{Message: responses.StatusNotFound}
}
