package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	name   string
	client *s3.Client
}

func NewBucket(ctx context.Context) (*Bucket, error) {
	if _, err := Getenv("AWS_ACCESS_KEY_ID"); err != nil {
		return nil, err
	}

	if _, err := Getenv("AWS_SECRET_ACCESS_KEY"); err != nil {
		return nil, err
	}

	if _, err := Getenv("AWS_REGION"); err != nil {
		return nil, err
	}

	name, err := Getenv("BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &Bucket{
		name:   name,
		client: client,
	}, nil
}

func (b *Bucket) allObjects(ctx context.Context) ([]string, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.name),
	}

	res, err := b.client.ListObjectsV2(ctx, params)
	if err != nil {
		return []string{}, err
	}

	var objects []string
	for _, object := range res.Contents {
		objects = append(objects, aws.ToString(object.Key))
	}

	return objects, nil
}

func (b *Bucket) downloadObject(ctx context.Context, key string) ([]byte, error) {
	object := &s3.GetObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
	}

	dl := manager.NewDownloader(b.client)
	buf := manager.NewWriteAtBuffer([]byte{})

	if _, err := dl.Download(ctx, buf, object); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
