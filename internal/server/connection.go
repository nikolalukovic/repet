package server

import (
	"bufio"
	"context"
	"net"
)

type RepetConnectionHandler struct {
	R   *bufio.Reader
	C   net.Conn
	Ctx context.Context
}

func NewRepetConnectionHandler(ctx context.Context, conn net.Conn) RepetConnectionHandler {
	r := bufio.NewReader(conn)
	return RepetConnectionHandler{
		R:   r,
		C:   conn,
		Ctx: ctx,
	}
}

func (h RepetConnectionHandler) ParseMessage() (RawMessage, error) {
	return extractMessage(h.R)
}

func (h RepetConnectionHandler) ExecuteCommand(msg RawMessage) error {
	cmd, err := parseCommand(msg)
	if err != nil {
		return err
	}

	switch cmd.(type) {
	case setCommand:
		err := executeSetCommand(cmd.(setCommand))
		if err != nil {
			return err
		}
		_, err = h.C.Write([]byte("0;2;OK"))
		if err != nil {
			return err
		}
		return nil
	case subCommand:
		return executeSubCommand(cmd.(subCommand))
	case getCommand:
		return executeGetCommand(h.C, cmd.(getCommand))
	default:
		return &RepetError{
			Code: CommandNotFound,
		}
	}
}
