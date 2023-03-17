package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

var (
	ns *server.Server
)

func startServer() error {
	serverOptions := &server.Options{
		JetStream: true,
	}
	ns, err := server.NewServer(serverOptions)
	if err != nil {
		panic(fmt.Sprintf("No NATS Server object returned: %v", err))
	}
	go ns.Start()
	if !ns.ReadyForConnections(10 * time.Second) {
		panic("Unable to start NATS Server in Go Routine")
	}
	return nil
}

func stopServer() {
	ns.Shutdown()
}
