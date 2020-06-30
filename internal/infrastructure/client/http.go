package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gromnsk/multiplexer/internal/infrastructure/config"
)

type HttpClient struct {
	timeout time.Duration
	client  *http.Client
}

func NewHttpClient(cfg config.ClientConfig) *HttpClient {
	return &HttpClient{
		client:  http.DefaultClient,
		timeout: cfg.PerRequestTimeout,
	}
}

func (c *HttpClient) Process(ctx context.Context, uri string) (response []byte, err error) {
	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(reqCtx, "GET", uri, nil)
	if err != nil {
		return
	}
	resp, err := c.client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
