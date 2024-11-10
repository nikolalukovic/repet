package server

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// RawMessage <version>;<length>;<content>
type RawMessage struct {
	Version int8
	Length  int64
	Content string
}

var EmptyRawMessage = RawMessage{}

func extractMessage(reader *bufio.Reader) (RawMessage, error) {
	strVersion, err := reader.ReadString(';')
	if err != nil {
		return EmptyRawMessage, err
	}

	strVersion = strings.TrimSpace(strVersion)

	if len(strVersion) > 2 {
		return EmptyRawMessage, &RepetError{
			Code: UnsupportedMessageVersion,
		}
	}

	version, err := strconv.Atoi(string(strVersion[0]))
	if err != nil {
		return EmptyRawMessage, err
	}

	if version != 0 {
		return EmptyRawMessage, &RepetError{
			Code: UnsupportedMessageVersion,
		}
	}

	return parseVersion0Message(reader)
}

func parseVersion0Message(reader *bufio.Reader) (RawMessage, error) {
	strLengthDirty, err := reader.ReadString(';')

	if err != nil {
		return EmptyRawMessage, err
	}

	splitStrLength := strings.Split(strLengthDirty, ";")

	if len(splitStrLength) < 2 {
		return EmptyRawMessage, &RepetError{
			Code: UnsupportedMessageVersion,
		}
	}

	length, err := strconv.ParseInt(splitStrLength[0], 10, 64)
	if err != nil {
		return EmptyRawMessage, err
	}

	buff := make([]byte, length)

	_, err = io.ReadFull(reader, buff)
	if err != nil {
		return EmptyRawMessage, err
	}

	return RawMessage{
		Version: 0,
		Length:  length,
		Content: string(buff),
	}, nil
}
