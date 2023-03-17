package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/nats-io/nats.go"
	"github.com/sunchero/http/handler"
	"github.com/sunchero/http/service"
)

var (
	nc *nats.Conn
	js nats.JetStreamContext
)

func main() {
	err := startServer()
	if err != nil {
		panic(fmt.Sprintf("cant start server: %v", err))
	}
	defer stopServer()
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println(err.Error())
	}
	defer nc.Close()
	// Create JetStream context
	js, err = nc.JetStream()
	if err != nil {
		log.Fatalln(err.Error())
	}
	svc := &service.Service{
		JS: js,
	}
	h := &handler.Handler{
		Svc: svc,
		NC:  nc,
	}
	server := &http.Server{
		Addr:              ":8080",
		Handler:           h.New(),
		TLSConfig:         &tls.Config{},
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){},
		ConnState: func(net.Conn, http.ConnState) {
		},
		ErrorLog: &log.Logger{},
		// BaseContext: func(net.Listener) context.Context {
		// },
		// ConnContext: func(ctx context.Context, c net.Conn) context.Context {
		// },
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
