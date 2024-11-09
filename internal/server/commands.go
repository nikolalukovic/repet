package server

import (
	"context"
)

type setCommand struct {
	key   string
	value string
	ttl   uint32
}

type subCommand struct {
	key string
}

func ExecuteCommand(ctx context.Context, message RawMessage) error {
	cmd, err := parseCommand(message)
	if err != nil {
		return err
	}

	switch cmd.(type) {
	case setCommand:
		return executeSetCommand(ctx, cmd.(setCommand))
	case subCommand:
		return executeSubCommand(ctx, cmd.(subCommand))
	default:
		return &RepetError{
			Code: CommandNotFound,
		}
	}
}

func executeSetCommand(ctx context.Context, command setCommand) error {
	return nil
}

func executeSubCommand(ctx context.Context, command subCommand) error {
	return nil
}

func parseCommand(message RawMessage) (interface{}, error) {
	return setCommand{}, nil
}
