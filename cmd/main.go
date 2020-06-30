package main

import (
	"log"
	"sync"

	"github.com/gromnsk/multiplexer/internal/infrastructure/client"
	"github.com/gromnsk/multiplexer/internal/infrastructure/config"
	"github.com/gromnsk/multiplexer/internal/infrastructure/handlers"
	"github.com/gromnsk/multiplexer/internal/infrastructure/signal"
	"github.com/gromnsk/multiplexer/internal/usecase"
)

func main() {
	cfg := config.MustConfigure()

	httpClient := client.NewHttpClient(cfg.Client)
	requestService := usecase.NewMultiplexer(httpClient, cfg.Client)

	server := handlers.NewServer(cfg.Http, requestService)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("starting multiplexer server on port: %d", cfg.Http.Port)
		err := server.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	sigHandler := signal.NewSignalHandler(func() error {
		log.Println("service shuts down")
		server.ShutDown()
		return nil
	})

	sigHandler.Poll()
	log.Println("Waiting while all goroutines will finish")
	wg.Wait()
	log.Println("Service stopped")
}
