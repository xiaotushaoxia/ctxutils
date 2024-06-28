package ctxutils

import (
	"context"
	"fmt"
	"os"
	"runtime"
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
	ctx, stop := NotifyContext(parent, sigs...)
	return ctx, func() {
		stop(context.Canceled)
	}
}

// WithSignalsCause behaves like [WithSignals] but also sets the cause of the
// returned Context when the get sigs. The returned [CancelFunc] does
// not set the cause.
func WithSignalsCause(parent context.Context, sigs ...os.Signal) (ctx context.Context, cancel context.CancelCauseFunc) {
	return NotifyContext(parent, sigs...)
}

// SignalsCtx
// WithXXX usually return ctx,cancel. This func only return ctx, so I don't named it WithXXX
// not return cancel? how to avoid goroutine leaks?  use runtime.SetFinalizer.
// note: done use this like !!!
// ctx := SignalsCtx(context.Background(), syscall.SIGINT)
// done := ctx.Done()
// ctx = nil  // !!!! no refer to ctx. make ctx be gc
func SignalsCtx(parent context.Context, signals ...os.Signal) (ctx context.Context) {
	var cancel context.CancelCauseFunc
	ctx, cancel = NotifyContext(parent, signals...)

	ctx = &wrapContext{ctx}

	runtime.SetFinalizer(ctx, func(context.Context) {
		cancel(context.Canceled)
	})
	return ctx
}

func SignalsCtxDefault(signals ...os.Signal) (ctx context.Context) {
	ss := []os.Signal{syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT}
	for _, signal := range signals {
		existed := false
		for _, s := range ss {
			if s == signal {
				existed = true
				break
			}
		}
		if !existed {
			ss = append(ss, signal)
		}
	}
	return SignalsCtx(context.Background(), ss...)
}

func NotifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, causeFunc context.CancelCauseFunc) {
	return notifyContext(parent, signals...)
}

// wrapContext
type wrapContext struct {
	context.Context
}

func (c *wrapContext) String() string {
	return fmt.Sprintf("%v", c.Context)
}
