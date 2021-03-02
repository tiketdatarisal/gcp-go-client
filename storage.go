package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type Storage struct {
	context *Context
}

func (s *Storage) newStorageManager(ctx context.Context, token string) (*storage.Client, error) {
	var result *storage.Client
	var err error
	if token != "" {
		t := &oauth2.Token{AccessToken: token, TokenType: "Bearer"}
		result, err = storage.NewClient(ctx, option.WithTokenSource(oauth2.StaticTokenSource(t)))
	} else {
		result, err = storage.NewClient(ctx)
	}

	if err != nil {
		status, err := parseGCPError(err)
		return nil, &HttpError{Code: status, Err: err}
	}

	return result, nil
}
