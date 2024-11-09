package server

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// <version>;<length>;<content>
type RawMessage struct {
	Version int8
	Length  int64
	Content string
}

var EmptyRawMessage = RawMessage{}

func ExtractMessage(reader *bufio.Reader) (RawMessage, error) {
	strVersion, err := reader.ReadString(';')
	if err != nil {
		return EmptyRawMessage, err
	}

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
	strLength, err := reader.ReadString(';')

	if err != nil {
		return EmptyRawMessage, err
	}

	l := strings.Split(strLength, ";")

	length, err := strconv.ParseInt(l[0], 10, 64)
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
