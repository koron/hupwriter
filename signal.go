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
			err := h.Reopen()
			if err != nil {
				panic("failed to reopen the log file: " + err.Error())
			}
		case os.Interrupt:
			h.removePid()
			os.Exit(0)
			return
		}
	}
}
