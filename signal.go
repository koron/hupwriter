package hupwriter

import (
	"os"
	"os/signal"
	"syscall"
)

func (h *HupWriter) signalStart() {
	signal.Notify(h.sig, syscall.SIGHUP, os.Interrupt)
	go h.signalMonitor()
}

func (h *HupWriter) signalMonitor() {
	for s := range h.sig {
		switch s {
		case syscall.SIGHUP:
			h.newLogFile()
		case os.Interrupt:
			h.removePid()
			return
		}
	}
}
