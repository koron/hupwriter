package hupwriter_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
	if runtime.GOOS == "windows" {
		t.Skip("windows can't rename opening log file")
	}
	tmpdir := t.TempDir()
	name := filepath.Join(tmpdir, "reopen.log")
	pid := filepath.Join(tmpdir, "reopen.pid")
	w, err := hupwriter.New(name, pid)
	if err != nil {
		t.Fatalf("failed to create hupwriter: %s", err)
	}
	defer w.Close()

	logger := log.New(w, "[hupwriter]", 0)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		cnt := 0
		for cnt < 20 {
			time.Sleep(10 * time.Millisecond)
			cnt++
			logger.Printf("%d", cnt)
			time.Sleep(90 * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		cnt := 0
		for cnt < 3 {
			time.Sleep(500 * time.Millisecond)
			cnt++
			rotate := filepath.Join(tmpdir, fmt.Sprintf("reopen.%d.log", cnt))
			if err := os.Rename(name, rotate); err != nil {
				t.Errorf("failed to rename: %s", err)
			}
			if err := w.Reopen(); err != nil {
				t.Errorf("reopen failed: %s", err)
				break
			}
		}
	}()
	wg.Wait()

	// Check contents of log files.
	testFileContents(t, filepath.Join(tmpdir, "reopen.1.log"), `[hupwriter]1
[hupwriter]2
[hupwriter]3
[hupwriter]4
[hupwriter]5
`)
	testFileContents(t, filepath.Join(tmpdir, "reopen.2.log"), `[hupwriter]6
[hupwriter]7
[hupwriter]8
[hupwriter]9
[hupwriter]10
`)
	testFileContents(t, filepath.Join(tmpdir, "reopen.3.log"), `[hupwriter]11
[hupwriter]12
[hupwriter]13
[hupwriter]14
[hupwriter]15
`)

	testFileContents(t, filepath.Join(tmpdir, "reopen.log"), `[hupwriter]16
[hupwriter]17
[hupwriter]18
[hupwriter]19
[hupwriter]20
`)
}

func testFileContents(t *testing.T, name string, want string) {
	t.Helper()
	f, err := os.Open(name)
	if err != nil {
		t.Errorf("failed to open file: %s", err)
		return
	}
	b, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		t.Errorf("failed to read file: %s", err)
		return
	}
	got := string(b)
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("content mismatch: -want +got\n%s", d)
	}
}
