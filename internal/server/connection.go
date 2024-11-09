package server

import (
	"bufio"
	"net"
)

type RepetConnectionHandler struct {
	C net.Conn
}

func NewRepetConnectionHandler(conn net.Conn) RepetConnectionHandler {
	return RepetConnectionHandler{
		C: conn,
	}
}

func (h RepetConnectionHandler) ParseMessage() (RawMessage, error) {
	_ = bufio.NewReader(h.C)

	return RawMessage{}, nil
}
