package server

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
)

func StartServer(ctx context.Context) {
	ipAddr := "localhost:8080"
	listener, err := net.Listen("tcp", ipAddr)
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

	LogInfo("Listening on: ", ipAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogWarning(err)
			continue
		}

		LogInfo("Got a connection")

		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	rcs := NewRepetConnectionHandler(ctx, conn)

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			LogError(err)
		}
	}(conn)

	for {
		rawMessage, err := rcs.ParseMessage()

		if err != nil {
			if errors.Is(err, io.EOF) {
				// client closed the connection
				// just get out of the handler
				LogInfo("Client closed connection")
				break
			}
			logAndClose(err, conn)
			return
		}

		err = rcs.ExecuteCommand(rawMessage)

		if err != nil {
			logAndClose(err, conn)
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
