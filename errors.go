package main

import (
	"fmt"
	"os"
	"strings"
)

type EnvironmentError struct {
	Missing []string
}

func (e *EnvironmentError) Error() string {
	vars := strings.Join(e.Missing, ", ")
	if len(e.Missing) > 0 {
		return fmt.Sprintf("environment variable(s) %s not set", vars)
	}
	return "environment variables missing"
}

func Getenv(key string) (string, error) {
	var (
		v   = os.Getenv(key)
		err = new(EnvironmentError)
	)

	if v == "" {
		err.Missing = append(err.Missing, key)
		return "", err
	}

	return v, nil
}
