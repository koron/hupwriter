/*
Package hupwriter provides wrapper of os.File.

The wrapper will close and reopen the file when a HUP signal is received.
It allows easier log file rotation.

By logging to hupwriter you can create log files that can be used with log
rotation management programs such as logrotate, newsyslog, or so.
*/
package hupwriter

import (
	"fmt"
	"os"
	"os/signal"
)

// HupWriter wraps os.File, and will close and reopen the file when a HUP
// signal is received.
type HupWriter struct {
	log  string
	pid  string
	sig  chan os.Signal
	file *os.File
}

// New creates a HupWriter with a file of name.
// It create PID file which records current process ID when pid is not empty.
func New(name, pid string) *HupWriter {
	if len(pid) != 0 {
		writePid(pid)
	}
	sig := make(chan os.Signal, 1)
	h := &HupWriter{log: name, pid: pid, sig: sig}
	h.signalStart()
	return h
}

// Write writes data to an underlying file.
func (h *HupWriter) Write(p []byte) (int, error) {
	if h.file == nil {
		_, err := h.newLogFile()
		if err != nil {
			return 0, err
		}
	}
	return h.file.Write(p)
}

// Close closes an underlying file, and stop to listening signals.
func (h *HupWriter) Close() error {
	h.closeLogFile()
	if h.sig != nil {
		signal.Stop(h.sig)
		close(h.sig)
		h.sig = nil
	}
	h.removePid()
	return nil
}

func (h *HupWriter) closeLogFile() {
	if h.file != nil {
		h.file.Sync()
		h.file.Close()
		h.file = nil
	}
}

func (h *HupWriter) newLogFile() (*os.File, error) {
	// Close old file.
	h.closeLogFile()
	// Open new load file.
	f, err := os.OpenFile(h.log, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open log file:", err)
		return nil, err
	}
	h.file = f
	return h.file, nil
}

// writePid
func writePid(pid string) {
	f, err := os.Create(pid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create pid file:", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "%d", os.Getpid())
}

// removePid removes a PID file.
func (h *HupWriter) removePid() {
	if len(h.pid) == 0 {
		return
	}
	os.Remove(h.pid)
}
