package server

import (
	"context"
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

	defer listener.Close()

	LogInfo("Listening on: ", ipAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogWarning(err)
			continue
		}

		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	rcs := NewRepetConnectionHandler(ctx, conn)

	defer conn.Close()

	for {
		rawMessage, err := rcs.ParseMessage()

		if err != nil {
			logAndClose(err, conn)
			return
		}

		LogInfo(rawMessage)

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
