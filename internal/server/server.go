package server

import (
	"bufio"
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
