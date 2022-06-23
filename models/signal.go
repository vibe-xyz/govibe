package models

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	SigHup = syscall.SIGHUP
)

type SignalHandler struct {
	sig syscall.Signal
	ch  chan os.Signal
}

func NewSignalHandler(sig syscall.Signal) *SignalHandler {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig)

	h := &SignalHandler{sig: sig, ch: ch}
	return h
}
func (h *SignalHandler) GetChan() <-chan os.Signal {
	return h.ch
}
