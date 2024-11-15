package server

import (
	"context"
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"io"
	"net"
	"os"
)

type Client struct {
	conn     net.Conn
	id       string
	name     *string
	channels []string
}

type Server struct {
	ListenAddr string
	// <client.id>:<client>
	conns map[string]*Client
}

func NewServer(listenAddr string) Server {
	return Server{
		ListenAddr: listenAddr,
		conns:      make(map[string]*Client),
	}
}

func (s Server) StartServer(ctx context.Context) {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		LogError(err)
		os.Exit(1)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			LogError(err)
		}
	}(listener)

	LogInfo("Listening on: ", s.ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogWarning(err)
			continue
		}

		id, err := gonanoid.New()
		if err != nil {
			LogError(err)
			err := conn.Close()
			if err != nil {
				LogError(err)
				continue
			}
			continue
		}

		cl := Client{
			conn: conn,
			id:   id,
		}

		s.conns[cl.id] = &cl

		LogInfo(fmt.Sprintf("Got a connection: %s", cl.id))

		go handleConnection(ctx, &cl, s)
	}
}

func handleConnection(ctx context.Context, client *Client, server Server) {
	rcs := NewRepetConnectionHandler(ctx, client)

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			LogError(err)
		}
	}(client.conn)

	for {
		rawMessage, err := rcs.ParseMessage()

		if err != nil {
			if errors.Is(err, io.EOF) {
				// client closed the connection
				// just get out of the handler
				LogInfo(fmt.Sprintf("Client %s closed connection", rcs.Client.id))
				break
			}
			logAndClose(err, client.conn)
			return
		}

		err = rcs.ExecuteCommand(rawMessage, server)

		if err != nil {
			logAndClose(err, client.conn)
			return
		}
	}
}

func logAndClose(err error, conn net.Conn) {
	LogError(err)

	connErr := conn.Close()

	if connErr != nil {
		LogError(err)
	}
}
