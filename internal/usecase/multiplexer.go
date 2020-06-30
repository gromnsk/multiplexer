package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gromnsk/multiplexer/internal/infrastructure/config"
	"github.com/gromnsk/multiplexer/internal/interface/client"
)

type Multiplexer struct {
	maxWorkers uint8
	maxUrls    int
	client     client.Client
}

type Response struct {
	sync.Mutex
	Results map[string][]byte
}

func NewMultiplexer(client client.Client, cfg config.ClientConfig) *Multiplexer {
	return &Multiplexer{
		client:     client,
		maxWorkers: cfg.MaxWorkers,
		maxUrls:    cfg.MaxUrls,
	}
}

func (m *Multiplexer) Validate(urls []string) error {
	if len(urls) > m.maxUrls {
		return fmt.Errorf("max amount of %d urls exceeded", m.maxUrls)
	}

	if len(urls) == 0 {
		return errors.New("empty urls list in request")
	}

	return nil
}

func (m *Multiplexer) ProcessRequests(ctx context.Context, urls []string) (response *Response, err error) {
	err = m.Validate(urls)
	if err != nil {
		return
	}
	var lastErr error
	response = &Response{
		Results: make(map[string][]byte, len(urls)),
	}
	limiter := make(chan struct{}, m.maxWorkers)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	for i := range urls {
		select {
		case <-ctx.Done():
			return response, lastErr
		default:
			wg.Add(1)
			limiter <- struct{}{}
			go func(i int) {
				resp, err := m.client.Process(ctx, urls[i])
				if err != nil {
					cancel()
					if lastErr == nil {
						lastErr = err
					}
				}
				response.Lock()
				response.Results[urls[i]] = resp
				response.Unlock()
				<-limiter
				wg.Done()
			}(i)
		}
	}
	wg.Wait()

	return
}
