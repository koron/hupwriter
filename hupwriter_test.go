package hupwriter_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/koron/hupwriter"
)

func TestBasic(t *testing.T) {
	tmpdir := t.TempDir()
	name := filepath.Join(tmpdir, "basic.log")
	pid := filepath.Join(tmpdir, "basic.pid")
	w, err := hupwriter.New(name, pid)
	if err != nil {
		t.Fatalf("failed to create hupwriter: %s", err)
	}

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

func TestReopen(t *testing.T) {
	tmpdir := t.TempDir()
	name := filepath.Join(tmpdir, "basic.log")
	pid := filepath.Join(tmpdir, "basic.pid")
	w, err := hupwriter.New(name, pid)
	if err != nil {
		t.Fatalf("failed to create hupwriter: %s", err)
	}
	defer w.Close()

	logger := log.New(w, "hupwriter", 0)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		cnt := 0
		for cnt < 20 {
			time.Sleep(100 * time.Millisecond)
			cnt++
			logger.Printf("%d", cnt)
		}
	}()
	go func() {
		defer wg.Done()
		cnt := 0
		for cnt < 3 {
			time.Sleep(500 * time.Millisecond)
			cnt++
			rotate := filepath.Join(tmpdir, fmt.Sprintf("basic.%d.log", cnt))
			if err := os.Rename(name, rotate); err != nil {
				t.Errorf("failed to rename: %s", err)
			}
			if err := w.Reopen(); err != nil {
				t.Errorf("reopen failed: %s", err)
				break
			}
			break
		}
	}()
	wg.Wait()

	entries, err := os.ReadDir(tmpdir)
	if err != nil {
		t.Fatalf("failed to readdir: %s", err)
	}
	// TODO: check output
	for i, e := range entries {
		t.Logf("#%d name=%s isdir=%t mode=%04o", i, e.Name(), e.IsDir(), e.Type())
	}
}
