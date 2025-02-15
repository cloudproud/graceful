package graceful

import (
	"context"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestContextSignal(t *testing.T) {
	ctx := NewContext(context.Background(), zaptest.NewLogger(t))

	mu := sync.Mutex{}
	expected := 3
	closed := 0

	counter := func() {
		time.Sleep(time.Second)
		mu.Lock()
		closed++
		mu.Unlock()
	}

	go ctx.Closer(counter)
	go ctx.Closer(counter)
	go ctx.Closer(counter)

	go func() {
		time.Sleep(time.Second)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		require.NoError(t, err)
	}()

	ctx.AwaitKillSignal()

	mu.Lock()
	require.Equal(t, expected, closed)
	mu.Unlock()
}

func TestForceContextSignal(t *testing.T) {
	ctx := NewContext(context.Background(), zaptest.NewLogger(t))

	mu := sync.Mutex{}
	expected := 0
	closed := 0

	counter := func() {
		time.Sleep(10 * time.Second)
		mu.Lock()
		closed++
		mu.Unlock()
	}

	go ctx.Closer(counter)
	go ctx.Closer(counter)
	go ctx.Closer(counter)

	go func() {
		time.Sleep(time.Second)
		// NOTE: the second signal represents the force kill
		syscall.Kill(syscall.Getpid(), syscall.SIGINT) //nolint:errcheck
		time.Sleep(time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT) //nolint:errcheck
	}()

	ctx.AwaitKillSignal()

	mu.Lock()
	require.Equal(t, expected, closed)
	mu.Unlock()
}
