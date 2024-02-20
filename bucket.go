package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	name   string
	client *s3.Client
}

func NewBucket(ctx context.Context) (*Bucket, error) {
	name, err := Getenv("BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	if err := checkAWSVariables(); err != nil {
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

func checkAWSVariables() error {
	_, confErr := os.Stat(config.DefaultSharedConfigFilename())
	_, credErr := os.Stat(config.DefaultSharedCredentialsFilename())
	if confErr == nil && credErr == nil {
		return nil
	}

	if _, err := Getenv("AWS_ACCESS_KEY_ID"); err != nil {
		return err
	}

	if _, err := Getenv("AWS_SECRET_ACCESS_KEY"); err != nil {
		return err
	}

	if _, err := Getenv("AWS_REGION"); err != nil {
		return err
	}

	return nil
}
