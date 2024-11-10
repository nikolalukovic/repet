package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const setDenom = "set"
const getDenom = "get"
const subDenom = "sub"

type setCommand struct {
	key   string
	value string
	ttl   time.Duration
}

type subCommand struct {
	key string
}

type getCommand struct {
	key string
}

func executeSetCommand(cmd setCommand) error {
	setValue(cmd)
	return nil
}

func executeSubCommand(command subCommand) error {
	return nil
}

func executeGetCommand(conn net.Conn, cmd getCommand) error {
	value, ok := getValue(cmd.key)
	if ok {
		_, err := conn.Write([]byte("0;" + strconv.Itoa(len(value)) + ";" + value))
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Write([]byte("0;0"))
		if err != nil {
			return nil
		}
	}
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
			ttl:   time.Duration(ttl * uint64(time.Millisecond)),
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
