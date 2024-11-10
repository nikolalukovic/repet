package main

import (
	"context"

	"github.com/nikolalukovic/repet/cmd/repet/commands"
	"github.com/nikolalukovic/repet/internal/server"
)

func main() {
	ctx := context.Background()

	server.InitCache()
	server.StartServer(ctx)

	commands.Execute()
}
