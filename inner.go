package ctxutils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// notifyContext  add CancelCause for signal.NotifyContext
func notifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, stop context.CancelCauseFunc) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	}
	ctx, cancelCause := context.WithCancelCause(parent)
	c := &signalCtx{
		Context:     ctx,
		cancelCause: cancelCause,
		signals:     signals,
	}
	c.ch = make(chan os.Signal, 1)

	signal.Notify(c.ch, signals...)
	if ctx.Err() == nil {
		go func() {
			select {
			case sig := <-c.ch:
				c.stop(SignalError{sig})
			case <-c.Done():
				c.stop(c.Err())
			}
		}()
	}

	return c, c.stop
}

type signalCtx struct {
	context.Context
	cancelCause context.CancelCauseFunc

	signals []os.Signal // for String()
	ch      chan os.Signal
}

func (c *signalCtx) stop(err error) {
	c.cancelCause(err)
	signal.Stop(c.ch)
}

type stringer interface {
	String() string
}

func (c *signalCtx) String() string {
	var buf []byte
	// We know that the type of c.Context is context.cancelCtx, and we know that the
	// String method of cancelCtx returns a string that ends with ".WithCancel".
	name := c.Context.(stringer).String()
	name = name[:len(name)-len(".WithCancel")]
	buf = append(buf, "signal.NotifyContext("+name...)
	if len(c.signals) != 0 {
		buf = append(buf, ", ["...)
		for i, s := range c.signals {
			buf = append(buf, s.String()...)
			if i != len(c.signals)-1 {
				buf = append(buf, ',')
			}
		}
		buf = append(buf, ']')
	}
	buf = append(buf, ')')
	return string(buf)
}
