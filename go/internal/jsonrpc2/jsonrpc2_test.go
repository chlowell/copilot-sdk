package jsonrpc2

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type discardWriteCloser struct{}

func (discardWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discardWriteCloser) Close() error {
	return nil
}

type nopReadCloser struct {
	io.Reader
}

func (nopReadCloser) Close() error {
	return nil
}

func TestRequest_ReturnsProcessErrorWhenDoneClosed(t *testing.T) {
	client := NewClient(discardWriteCloser{}, nopReadCloser{Reader: strings.NewReader("")})

	done := make(chan struct{})
	processErr := errors.New("cli process exited with stderr")
	client.SetProcessDone(done, &processErr)
	close(done)

	_, err := client.Request("test.method", map[string]string{"input": "x"})
	if err == nil {
		t.Fatal("Expected process error, got nil")
	}
	if !errors.Is(err, processErr) {
		t.Errorf("Expected error %q, got %q", processErr.Error(), err.Error())
	}
}
