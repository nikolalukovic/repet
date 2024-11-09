package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
)

type Server struct {
	listenAddr net.TCPAddr
}

func StartServer() {
	ipAddr := "localhost:8080"
	listener, err := net.Listen("tcp", ipAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Printf("Listening on: %s\n", ipAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		go handleConnection(&conn)
	}
}

func handleConnection(conn *net.Conn) {
	reader := bufio.NewReader(*conn)

	rawMessage, err := ExtractMessage(reader)

	if err != nil {
		logAndClose(err, conn)
		return
	}

	err = ExecuteCommand(context.Background(), rawMessage)

	if err != nil {
		logAndClose(err, conn)
		return
	}

	// version, err := reader.ReadString(';')
	// if err != nil {
	// 	LogError(err)
	// 	err := (*conn).Close()
	// 	if err != nil {
	// 		LogError(err)
	// 	}
	// }
	//
	// length, err := reader.ReadBytes(';')

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			err := (*conn).Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			return
		}

		(*conn).Write([]byte(fmt.Sprintf("Received: %s", message)))
	}

}

func logAndClose(err error, conn *net.Conn) {
	LogError(err)

	connErr := (*conn).Close()

	if connErr != nil {
		LogError(err)
	}
}
