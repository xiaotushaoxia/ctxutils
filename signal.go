package ctxutils

import (
	"context"
	"fmt"
	"os"
	"runtime"
)

type SignalError struct {
	Signal os.Signal
}

func (e SignalError) Error() string {
	return fmt.Sprintf("got signal: %s", e.Signal)
}

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
func SignalsCtx(parent context.Context, signals ...os.Signal) (ctx context.Context) {
	ctx, _ = NotifyContext(parent, signals...)
	return ctx
}

func NotifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, causeFunc context.CancelCauseFunc) {
	ctx1, _causeFunc := notifyContext(parent, signals...)

	// ctx1 can’t be gc automatically
	// because I make a daemon goroutine for watching signal in notifyContext  and
	// context package may make a daemon goroutine for watch parent
	// wrap ctx1 into wrapContext, wrapContext can be gc simply
	ctx = &wrapContext{ctx1}
	//ctx = ctx1
	runtime.SetFinalizer(ctx, func(context.Context) {
		_causeFunc(context.Canceled)
	})

	return ctx, func(err error) {
		runtime.SetFinalizer(ctx, nil)
		_causeFunc(err)
	}
}

// wrapContext
type wrapContext struct {
	context.Context
}

func (c *wrapContext) String() string {
	return fmt.Sprintf("%v", c.Context)
}
