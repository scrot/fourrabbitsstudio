package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	name   string
	client *s3.Client
}

func NewBucket(ctx context.Context, name string) *Bucket {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "aws error: %s", err)
	}
	client := s3.NewFromConfig(cfg)

	return &Bucket{
		name:   name,
		client: client,
	}
}

func (b *Bucket) allObjects() ([]string, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.name),
	}

	res, err := b.client.ListObjectsV2(context.TODO(), params)
	if err != nil {
		return []string{}, err
	}

	var objects []string
	for _, object := range res.Contents {
		objects = append(objects, aws.ToString(object.Key))
	}

	return objects, nil
}
