package storage

import "context"

type Bucket interface {
	All(context.Context) ([]string, error)
	ObjectURL(context.Context, string) (string, error)
}
