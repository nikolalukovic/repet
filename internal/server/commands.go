package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const setDenom = "set"
const getDenom = "get"

const pubDenom = "pub"
const subDenom = "sub"

const idDenom = "id"

type setCommand struct {
	key   string
	value string
	ttl   time.Duration
}

type getCommand struct {
	key string
}

type subCommand struct {
	channel string
}

type pubCommand struct {
	channel string
	value   string
}

type nameCommand struct {
	name string
}

func executeSetCommand(cmd setCommand) error {
	setValue(cmd)
	return nil
}

func executeGetCommand(c Client, cmd getCommand) error {
	value, ok := getValue(cmd.key)
	if ok {
		fmtMsg := fmt.Sprintf("0;%d;%s", len(value), value)
		_, err := c.conn.Write([]byte(fmtMsg))
		if err != nil {
			return err
		}
	} else {
		_, err := c.conn.Write([]byte("0;0"))
		if err != nil {
			return nil
		}
	}
	return nil
}

func executeSubCommand(c *Client, command subCommand) error {
	c.channels = append(c.channels, command.channel)
	LogInfo(fmt.Sprintf("Connection %s subscribed to %s", c.id, command.channel))
	return nil
}

func executePubCommand(c Client, s Server, cmd pubCommand) error {
	LogInfo(fmt.Sprintf("Connection %s wants to publish \"%s\" to channel %s", c.id, cmd.value, cmd.channel))
	for _, client := range s.conns {
		for _, channel := range client.channels {
			if cmd.channel == channel {
				msg := fmt.Sprintf("chan %s %s", channel, cmd.value)
				fmtMsg := fmt.Sprintf("0;%d;%s", len(msg), msg)
				_, err := client.conn.Write([]byte(fmtMsg))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func executeNameCommand(c *Client, cmd nameCommand) error {
	c.name = &cmd.name
	LogInfo(fmt.Sprintf("Connection %s has been named to \"%s\"", c.id, *c.name))
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
		channel := parts[1]
		return subCommand{
			channel: channel,
		}, nil
	case pubDenom:
		if len(parts) != 3 {
			return nil, &RepetError{
				Code:    MalformedCommand,
				Details: fmt.Sprintf("%v", parts),
			}
		}
		channel := parts[1]
		value := parts[2]

		return pubCommand{
			channel: channel,
			value:   value,
		}, nil
	case idDenom:
		if len(parts) != 2 {
			return nil, &RepetError{
				Code:    MalformedCommand,
				Details: fmt.Sprintf("%v", parts),
			}
		}
		id := parts[1]
		return nameCommand{
			name: id,
		}, nil
	default:
		return nil, &RepetError{
			Code:    UnableToParseMessage,
			Details: fmt.Sprintf("Unknown command: %v", denom),
		}
	}
}
