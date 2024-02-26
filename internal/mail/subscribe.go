package mail

import (
	"context"

	"github.com/mailerlite/mailerlite-go"
	"github.com/scrot/fourrabbitsstudio/internal/errors"
)

type Subscriber struct {
	client *mailerlite.Client
}

func NewSubscriber() (*Subscriber, error) {
	token, err := errors.Getenv("MAILERLITE_TOKEN")
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

	respSub, _, err := s.client.Subscriber.Create(ctx, sub)
	if err != nil {
		return err
	}

	if count.Total <= 100 {
		early := "113952487829931325"
		if _, _, err := s.client.Group.Assign(ctx, early, respSub.Data.ID); err != nil {
			return err
		}
	} else {
		normal := "113953297639933753"
		if _, _, err := s.client.Group.Assign(ctx, normal, respSub.Data.ID); err != nil {
			return err
		}
	}

	return nil
}
