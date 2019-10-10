package commands

import "bytes"

func isNoReply(payload string) bool {
	return bytes.Equal([]byte(payload), []byte("noreply"))
}
