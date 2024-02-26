package storage

import (
	"context"
	"fmt"

	"github.com/scrot/fourrabbitsstudio/internal/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Bucket struct {
	name   string
	client *s3.Client
}

func NewBucket(ctx context.Context) (Bucket, error) {
	if _, err := errors.Getenv("AWS_ACCESS_KEY_ID"); err != nil {
		return nil, err
	}

	if _, err := errors.Getenv("AWS_SECRET_ACCESS_KEY"); err != nil {
		return nil, err
	}

	if _, err := errors.Getenv("AWS_REGION"); err != nil {
		return nil, err
	}

	name, err := errors.Getenv("BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Bucket{
		name:   name,
		client: client,
	}, nil
}

func (b *S3Bucket) All(ctx context.Context) ([]string, error) {
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

func (b *S3Bucket) ObjectURL(ctx context.Context, name string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
