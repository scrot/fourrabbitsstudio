package main

import (
	"context"

	"github.com/mailerlite/mailerlite-go"
)

type Subscriber struct {
	client *mailerlite.Client
}

func NewSubscriber() (*Subscriber, error) {
	token, err := Getenv("MAILERLITE_TOKEN")
	if err != nil {
		return nil, err
	}
	client := mailerlite.NewClient(token)
	return &Subscriber{client}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, email string) error {
	sub := &mailerlite.Subscriber{
		Email: email,
	}

	_, _, err := s.client.Subscriber.Create(ctx, sub)
	if err != nil {
		return err
	}

	return nil
}
