package log_test

import (
	"context"
	"testing"

	"github.com/kosatnkn/log"
)

func newLogger(t *testing.T) log.AdapterInterface {

	cfg := log.Config{
		Level:   "ERROR",
		Colors:  true,
		Console: true,
	}

	l, err := log.NewAdapter(cfg)
	if err != nil {
		t.Fatalf("Error creating logger %v", err)
	}

	return l
}

func TestMessage(t *testing.T) {

	l := newLogger(t)

	l.Error(context.Background(), "Hello")
	l.Error(context.Background(), "Hello", "Additional 1", "Additional 2")
}
