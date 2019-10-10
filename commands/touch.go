package commands

import (
	"fmt"
	"github.com/blizztrack/gmc"
	"github.com/blizztrack/gmc/lru"
	"strconv"
)

type TouchCommand struct{}

func (get *TouchCommand) Handle(payload []string) gmc.Response {
	if len(payload) < 2 || len(payload) > 3 {
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
	ExpTime, err := strconv.ParseInt(payload[1], 10, 64)
	if err != nil {
		return gmc.MessageResponse{Message: fmt.Sprintf(gmc.StatusClientError, err)}
	}

	item.SetExpires(ExpTime)

	if len(payload) == 3 && isNoReply(payload[2]) {
		return nil
	}

	return gmc.MessageResponse{Message: gmc.StatusTouched}
}