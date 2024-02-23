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

	count, _, err := s.client.Subscriber.Count(ctx)
	if err != nil {
		return err
	}

	if count.Total <= 100 {
		early := mailerlite.Group{ID: "113952487829931325"}
		sub.Groups = append(sub.Groups, early)
	} else {
		normal := mailerlite.Group{ID: "113953297639933753"}
		sub.Groups = append(sub.Groups, normal)
	}

	if _, _, err := s.client.Subscriber.Create(ctx, sub); err != nil {
		return err
	}

	return nil
}
