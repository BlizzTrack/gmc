package commands

import (
	"bufio"
	"bytes"
	"io"
)

func isNoReply(payload string) bool {
	return bytes.Equal([]byte(payload), []byte("noreply"))
}

func readLine(reader io.Reader) (line []byte, err error) {
	line, err = bufio.NewReader(reader).ReadBytes('\n')
	return
}