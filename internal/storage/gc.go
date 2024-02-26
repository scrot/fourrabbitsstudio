package storage

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	gstorage "cloud.google.com/go/storage"
	xerrors "github.com/scrot/fourrabbitsstudio/internal/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GCBucket struct {
	name   string
	client *gstorage.Client
}

func NewGCBucket(ctx context.Context, name string) (Bucket, error) {
	enc, err := xerrors.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		return nil, errors.New("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	cred, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(cred)
	client, err := gstorage.NewClient(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &GCBucket{
		name:   name,
		client: client,
	}, nil
}

func (b *GCBucket) All(ctx context.Context) ([]string, error) {
	bkt := b.client.Bucket(b.name)

	var names []string
	it := bkt.Objects(ctx, &gstorage.Query{Prefix: ""})
	for {
		obj, err := it.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return []string{}, err
		}
		names = append(names, obj.Name)
	}

	return names, nil
}

func (b *GCBucket) ObjectURL(ctx context.Context, name string) (string, error) {
	return b.client.Bucket(b.name).SignedURL(name, &gstorage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(10 * time.Minute),
	})
}
