package signal

import (
	"os"
	"os/signal"
	"syscall"
)

type signalHandler struct {
	sig      chan os.Signal
	callback func() error
}

func NewSignalHandler(callback func() error) *signalHandler {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return &signalHandler{
		sig:      sig,
		callback: callback,
	}
}

func (sh *signalHandler) Poll() {
	<-sh.sig
	callbackError := sh.callback()
	if callbackError != nil {
		panic(callbackError)
	}
}
