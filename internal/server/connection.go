package server

import (
	"bufio"
	"context"
)

type RepetConnectionHandler struct {
	Reader  *bufio.Reader
	Client  *Client
	Context context.Context
}

func NewRepetConnectionHandler(ctx context.Context, c *Client) *RepetConnectionHandler {
	r := bufio.NewReader(c.conn)
	return &RepetConnectionHandler{
		Reader:  r,
		Client:  c,
		Context: ctx,
	}
}

func (h *RepetConnectionHandler) ParseMessage() (RawMessage, error) {
	return extractMessage(h.Reader)
}

func (h *RepetConnectionHandler) ExecuteCommand(msg RawMessage, server Server) error {
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
		_, err = h.Client.conn.Write([]byte("0;2;OK"))
		if err != nil {
			return err
		}
		return nil
	case getCommand:
		return executeGetCommand(*h.Client, cmd.(getCommand))
	case subChanCommand:
		return executeSubChanCommand(h.Client, cmd.(subChanCommand))
	case subKeyCommand:
		return executeSubKeyCommand(h.Client, cmd.(subKeyCommand))
	case pubCommand:
		return executePubCommand(*h.Client, server, cmd.(pubCommand))
	case nameCommand:
		return executeNameCommand(h.Client, cmd.(nameCommand))
	default:
		return &RepetError{
			Code: CommandNotFound,
		}
	}
}
