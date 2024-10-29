package main

import (
	"fmt"
	"os"

	"github.com/nikolalukovic/repet/internal/server"
)

func main() {
	cfg, err := server.ConsumeConfiguration()

	if err != nil {
		server.LogError(err)
		os.Exit(1)
	}

	fmt.Printf("%#v\n", cfg)
}
