package hupwriter_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/koron/hupwriter"
)

func TestBasic(t *testing.T) {
	tmpdir := t.TempDir()
	name := filepath.Join(tmpdir, "basic.log")
	pid := filepath.Join(tmpdir, "basic.pid")
	w := hupwriter.New(name, pid)

	// Check pid file
	if _, err := os.Stat(pid); err != nil {
		w.Close()
		t.Fatalf("failed to stat pid file %q: %s", pid, err)
	}
	pidBytes, err := os.ReadFile(pid)
	if err != nil {
		w.Close()
		t.Fatalf("failed to read pid file %q: %s", pid, err)
	}
	if want, got := fmt.Sprintf("%d", os.Getpid()), string(pidBytes); got != want {
		w.Close()
		t.Fatalf("the content of pid is missmatch:\nwant=%q\n got=%q", want, got)
	}

	// Write a line then Close.
	if _, err := io.WriteString(w, "Hello hupwriter!\n"); err != nil {
		w.Close()
		t.Fatalf("failed to write: %s", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close: %s", err)
	}

	// Check pid file is removed.
	if _, err := os.Stat(pid); !os.IsNotExist(err) {
		t.Fatalf("unexpected error of (removed) stat pid file: %s", err)
	}

	// Check contents of name file.
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("failed to read %q: %s", name, err)
	}
	if want, got := "Hello hupwriter!\n", string(b); got != want {
		t.Errorf("the content of output is missmatch:\nwant=%q\n got=%q", want, got)
	}
}
