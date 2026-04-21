//go:build !cgo

package main

import (
	"bytes"
	"os"
	"sync"
	"testing"
)

// stdioMutex serializes tests that temporarily redirect stdout/stderr in
// non-CGO builds. The cgo test path defines the same helper in its own
// test-only file, so we keep this behind !cgo to avoid duplicate symbols.
var stdioMutex sync.Mutex

func captureStdout(t *testing.T, fn func() error) string {
	t.Helper()

	stdioMutex.Lock()
	defer stdioMutex.Unlock()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	return buf.String()
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	stdioMutex.Lock()
	defer stdioMutex.Unlock()

	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stderr = oldStderr

	return buf.String()
}
