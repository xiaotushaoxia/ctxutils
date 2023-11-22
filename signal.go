package ctxutils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

// WithSignalsCause and WithSignals format is consistent with std context
// like context.WithTimeout和context.WithTimeoutCause

// WithSignals
// The returned [Context.Done] channel is closed when the deadline expires, when the
// returned cancel function is called, or when the parent context's Done channel
// is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithSignals(parent context.Context, sigs ...os.Signal) (ctx context.Context, cancel context.CancelFunc) {
	return context.WithCancel(SignalsCtx(parent, sigs...))
}

// WithSignalsCause behaves like [WithSignals] but also sets the cause of the
// returned Context when the get sigs. The returned [CancelFunc] does
// not set the cause.
func WithSignalsCause(parent context.Context, sigs ...os.Signal) (ctx context.Context, cancel context.CancelCauseFunc) {
	return context.WithCancelCause(SignalsCtx(parent, sigs...))
}

// SignalsCtx
// WithXXX usually return ctx,cancel. This func only return ctx, so I don't named it WithXXX
func SignalsCtx(parent context.Context, sigs ...os.Signal) (ctx context.Context) {
	return newSignalContext(parent, sigs...)
}

func newSignalContext(parent context.Context, sigs ...os.Signal) *signalContext {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	}

	ctx := &signalContext{
		Context:    parent,
		done:       make(chan struct{}),
		signalChan: make(chan os.Signal),
	}
	signal.Notify(ctx.signalChan, sigs...)

	go func() {
		select {
		case sig := <-ctx.signalChan:
			ctx.err.Store(SignalError{sig})
		case <-parent.Done():
			ctx.err.Store(parent.Err())
		}
		close(ctx.done)
	}()

	return ctx
}

type signalContext struct {
	// Context is the parent context
	context.Context
	done chan struct{}
	err  atomic.Value

	signalChan chan os.Signal // for easy test
}

func (s *signalContext) Done() <-chan struct{} {
	return s.done
}

func (s *signalContext) Err() error {
	val := s.err.Load()
	if val == nil {
		return nil
	}
	return val.(error)
}

type SignalError struct {
	Signal os.Signal
}

func (e SignalError) Error() string {
	return fmt.Sprintf("got signal: %s", e.Signal)
}
