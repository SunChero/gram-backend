package service

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func (s *Service) CreateStream(n *nats.StreamConfig) (*nats.StreamInfo, error) {
	//var si *nats.StreamInfo
	si, err := s.JS.AddStream(n)
	if err != nil {
		return nil, err
	}
	return si, nil
}

func (s *Service) UpdateStream(n *nats.StreamConfig) (*nats.StreamInfo, error) {
	si, err := s.JS.UpdateStream(n)
	if err != nil {
		return nil, err
	}
	return si, nil
}

func (s *Service) DeleteStream(str string) (string, error) {
	err := s.JS.DeleteStream(str)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (s *Service) ListStreams() []*nats.StreamInfo {
	var str []*nats.StreamInfo
	for info := range s.JS.StreamsInfo() {
		fmt.Println("stream name:", info.Config.Name)
		str = append(str, info)
	}
	return str
}

//create a durable consumer on the stream
func (s *Service) JoinStream(stream string, name string) (*nats.ConsumerInfo, error) {
	ci, err := s.JS.AddConsumer(stream, &nats.ConsumerConfig{
		Durable:        name,
		DeliverSubject: "socket." + name,
		AckPolicy:      nats.AckAllPolicy,
	})
	if err != nil {
		return nil, err
	}
	return ci, nil
}

//delete durable consumer from stream
func (s *Service) LeaveStream(stream string, name string) error {
	err := s.JS.DeleteConsumer(stream, name)
	return err
}
