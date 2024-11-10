package server

import (
	"bufio"
	"context"
	"net"
)

type RepetConnectionHandler struct {
	C   net.Conn
	Ctx context.Context
}

func NewRepetConnectionHandler(ctx context.Context, conn net.Conn) RepetConnectionHandler {
	return RepetConnectionHandler{
		C:   conn,
		Ctx: ctx,
	}
}

func (h RepetConnectionHandler) ParseMessage() (RawMessage, error) {
	r := bufio.NewReader(h.C)

	return extractMessage(r)
}

func (h RepetConnectionHandler) ExecuteCommand(msg RawMessage) error {
	return executeCommand(h.Ctx, msg)
}
