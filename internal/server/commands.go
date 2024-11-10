package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const setDenom = "set"
const getDenom = "get"
const subDenom = "sub"

type setCommand struct {
	key   string
	value string
	ttl   uint64
}

type subCommand struct {
	key string
}

type getCommand struct {
	key string
}

func executeCommand(ctx context.Context, message RawMessage) error {
	cmd, err := parseCommand(message)
	if err != nil {
		return err
	}

	switch cmd.(type) {
	case setCommand:
		return executeSetCommand(ctx, cmd.(setCommand))
	case subCommand:
		return executeSubCommand(ctx, cmd.(subCommand))
	case getCommand:
		return executeGetCommand(ctx, cmd.(getCommand))
	default:
		return &RepetError{
			Code: CommandNotFound,
		}
	}
}

func executeSetCommand(ctx context.Context, command setCommand) error {
	LogInfo(command)
	return nil
}

func executeSubCommand(ctx context.Context, command subCommand) error {
	LogInfo(command)
	return nil
}

func executeGetCommand(ctx context.Context, command getCommand) error {
	LogInfo(command)
	return nil
}

func parseCommand(message RawMessage) (interface{}, error) {
	parts := strings.SplitN(message.Content, " ", 4)
	denom := strings.ToLower(parts[0])

	switch denom {
	case setDenom:
		if len(parts) != 4 {
			return nil, &RepetError{
				Code:    MalformedCommand,
				Details: fmt.Sprintf("%v", parts),
			}
		}

		key := parts[1]
		strTtl := parts[2]
		value := parts[3]

		ttl, err := strconv.ParseUint(strTtl, 10, 32)

		if err != nil {
			return nil, &RepetError{
				Code:    MalformedCommand,
				Details: fmt.Sprintf("%v", parts),
			}
		}

		return setCommand{
			key:   key,
			value: value,
			ttl:   ttl,
		}, nil
	case getDenom:
		if len(parts) != 2 {
			return nil, &RepetError{
				Code:    MalformedCommand,
				Details: fmt.Sprintf("%v", parts),
			}
		}
		key := parts[1]

		return getCommand{
			key: key,
		}, nil
	case subDenom:
		if len(parts) != 2 {
			return nil, &RepetError{
				Code:    MalformedCommand,
				Details: fmt.Sprintf("%v", parts),
			}
		}
		key := parts[1]
		return subCommand{
			key: key,
		}, nil
	default:
		return nil, &RepetError{
			Code:    UnableToParseMessage,
			Details: fmt.Sprintf("Unknown command: %v", denom),
		}
	}
}
