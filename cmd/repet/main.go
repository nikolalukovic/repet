package main

import (
	"context"

	"github.com/nikolalukovic/repet/cmd/repet/commands"
	"github.com/nikolalukovic/repet/internal/server"
)

func main() {
	ctx := context.Background()

	server.InitCache()

	srv := server.NewServer("localhost:8080")

	srv.StartServer(ctx)

	commands.Execute()
}
