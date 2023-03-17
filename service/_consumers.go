package service

import "github.com/nats-io/nats.go"

func (s *Service) CreateConsumer(stream string, name string) (*nats.ConsumerInfo, error) {
	ci, err := s.JS.AddConsumer(stream, &nats.ConsumerConfig{
		Durable:        name,
		DeliverSubject: "socket." + name,
	})
	if err != nil {
		return nil, err
	}
	return ci, nil
}

func (s *Service) DeleteConsumer(stream string, name string) error {
	err := s.JS.DeleteConsumer(stream, name)
	return err
}
