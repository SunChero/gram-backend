package service

import "github.com/nats-io/nats.go"

type Service struct {
	JS nats.JetStreamContext
}
