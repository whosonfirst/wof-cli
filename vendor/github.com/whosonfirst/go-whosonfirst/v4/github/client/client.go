package client

import (
	"context"

	"github.com/google/go-github/v88/github"
)

func NewClient(ctx context.Context, token string) (*github.Client, error) {

	client_opts := make([]github.ClientOptionsFunc, 0)

	if token != "" {
		client_opts = append(client_opts, github.WithAuthToken(token))
	}

	return github.NewClient(client_opts...)
}
