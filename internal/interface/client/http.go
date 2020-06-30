package client

import (
	"context"
)

type Client interface {
	Process(ctx context.Context, uri string) ([]byte, error)
}
