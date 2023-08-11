package main

import (
	"log"
	"runtime"

	"github.com/Awareness-Labs/rainforest/pkg/server"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	s := server.NewServer(nc)
	s.Start()

	// Keep the connection alive until terminated
	runtime.Goexit()
	nc.Close()
}
