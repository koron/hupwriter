/*
Package hupwriter provides wrapper of os.File.

The wrapper will close and reopen the file when a HUP signal is received.
It allows easier log file rotation.

By logging to hupwriter you can create log files that can be used with log
rotation management programs such as logrotate, newsyslog, or so.
*/
package hupwriter

import (
	"io"
	"io/fs"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

// HupWriter wraps os.File, and will close and reopen the file when a HUP
// signal is received.
type HupWriter struct {
	log  string
	pid  string
	sig  chan os.Signal
	file *os.File

	lock   sync.Mutex
	closed bool
}

func openFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
}

// New creates a HupWriter with a file of name.
// It create PID file which records current process ID when pid is not empty.
func New(name, pid string) (*HupWriter, error) {
	// Write pid file.
	if len(pid) != 0 {
		if err := writePid(pid); err != nil {
			return nil, err
		}
	}
	// Open "name" file to append log.
	file, err := openFile(name)
	if err != nil {
		os.Remove(pid)
		return nil, err
	}
	// Compose HupWriter.
	sig := make(chan os.Signal, 1)
	h := &HupWriter{log: name, pid: pid, sig: sig, file: file}
	h.signalStart()
	return h, nil
}

// Write writes data to an underlying file.
func (h *HupWriter) Write(p []byte) (int, error) {
	// status check
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return 0, fs.ErrClosed
	}
	return h.file.Write(p)
}

// Close closes an underlying file, and stop to listening signals.
func (h *HupWriter) Close() error {
	// status check
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return fs.ErrClosed
	}
	// close log file.
	h.file.Close()
	h.closed = true
	// terminate signal monitor.
	if h.sig != nil {
		signal.Stop(h.sig)
		close(h.sig)
		h.sig = nil
	}
	// remove pid file.
	h.removePid()
	return nil
}

// Reopen closes the output file and reopens it.
func (h *HupWriter) Reopen() error {
	// status check
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return fs.ErrClosed
	}
	// close and reopen the log file
	if err := h.file.Close(); err != nil {
		return err
	}
	f, err := openFile(h.log)
	if err != nil {
		return err
	}
	h.file = f
	return nil
}

// writePid
func writePid(pid string) error {
	f, err := os.Create(pid)
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, strconv.Itoa(os.Getpid()))
	f.Close()
	return err
}

// removePid removes a PID file.
func (h *HupWriter) removePid() {
	if len(h.pid) == 0 {
		return
	}
	os.Remove(h.pid)
}
