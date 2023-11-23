package ctxutils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
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
	//return context.WithCancel(SignalsCtx(parent, sigs...))
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
	return
}

func NotifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, stop context.CancelCauseFunc) {
	parent2, cancelCause := context.WithCancelCause(parent)
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	}
	ctx = &signalCtx{
		Context:     parent2,
		cancelCause: cancelCause,
		signals:     signals,
	}
	ch := make(chan os.Signal, 1)

	var once sync.Once
	stop = func(cause error) {
		once.Do(func() {
			cancelCause(cause)
			signal.Stop(ch)
		})
	}

	signal.Notify(ch, signals...)
	if parent2.Err() == nil {
		go watch(ch, parent2, stop)
	}

	runtime.SetFinalizer(ctx, func(*signalCtx) {
		stop(context.Canceled)
	})
	return
}

type signalCtx struct {
	context.Context
	cancelCause context.CancelCauseFunc

	signals []os.Signal // for String()
}

func (c *signalCtx) String() string {
	var buf []byte
	// We know that the type of c.Context is context.cancelCtx, and we know that the
	// String method of cancelCtx returns a string that ends with ".WithCancel".
	name := c.Context.(interface {
		String() string
	}).String()
	name = name[:len(name)-len(".WithCancel")]
	buf = append(buf, "signalCtx("+name...)
	if len(c.signals) != 0 {
		buf = append(buf, ", ["...)
		for i, s := range c.signals {
			buf = append(buf, s.String()...)
			if i != len(c.signals)-1 {
				buf = append(buf, ' ')
			}
		}
		buf = append(buf, ']')
	}
	buf = append(buf, ')')
	return string(buf)
}

func watch(ch chan os.Signal, parent2 context.Context, do func(err error)) {
	var err error
	select {
	case sig := <-ch:
		err = SignalError{sig}
	case <-parent2.Done():
		err = parent2.Err()
	}
	do(err)
}
