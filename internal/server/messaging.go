package server

import (
	"bufio"
)

// <version>;<length>;<content>
type RawMessage struct {
	version int8
	length  int64
	content string
}

func ExtractMessage(reader *bufio.Reader) (RawMessage, error) {
	return RawMessage{}, nil
}
