package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/netutil"

	"github.com/gromnsk/multiplexer/internal/infrastructure/config"
	"github.com/gromnsk/multiplexer/internal/usecase"
)

type Server struct {
	cfg               config.HttpConfig
	multiplexer       *usecase.Multiplexer
	multiplexerServer http.Server
}

type Response struct {
	Data  map[string][]byte `json:"data,omitempty"`
	Error string            `json:"error,omitempty"`
}

type Request struct {
	Urls []string `json:"urls"`
}

func NewServer(
	cfg config.HttpConfig,
	multiplexer *usecase.Multiplexer,
) *Server {
	s := Server{
		cfg:         cfg,
		multiplexer: multiplexer,
	}

	r := http.NewServeMux()
	r.HandleFunc("/", s.MainHandler())

	s.multiplexerServer = http.Server{
		Handler:           r,
		ReadHeaderTimeout: cfg.Timeout,
	}
	s.multiplexerServer.SetKeepAlivesEnabled(false)

	return &s
}

func (s *Server) Run() error {
	socket := fmt.Sprintf(":%d", s.cfg.Port)
	httpLis, err := net.Listen("tcp", socket)
	limitListener := netutil.LimitListener(httpLis, s.cfg.ConnectionsLimit)
	if err != nil {
		return err
	}
	return s.multiplexerServer.Serve(limitListener)
}

func failedResponse(err error, w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	response := Response{
		Error: err.Error(),
	}
	resp, err := json.Marshal(response)
	if err != nil {
		_, _ = w.Write([]byte("can't Marshal response"))
		return
	}

	_, _ = w.Write(resp)
}

func (s *Server) ShutDown() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.multiplexerServer.Shutdown(ctx); err != nil {
		log.Printf("server shuts down incorrectly: %#v", err)
	}
}
