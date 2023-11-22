package ctxutils

import (
	"context"
	"fmt"
	"syscall"
	"testing"
)

func TestWithSignalsCancelCause(t *testing.T) {
	cause := fmt.Errorf("test cause")

	ctx, cancelCause := WithSignalsCause(context.Background())
	cancelCause(cause)
	<-ctx.Done()
	//fmt.Println(ctx.Err())          //context canceled
	//fmt.Println(context.Cause(ctx)) // test cause

	if ctx.Err() != context.Canceled {
		t.Fatalf("ctx err is not context.Canceled")
	}
	if context.Cause(ctx) != cause {
		t.Fatalf("ctx Cause err is not input")
	}

	ctx2, cancel := WithSignals(context.Background())
	cancel()
	<-ctx2.Done()
	//fmt.Println(ctx2.Err())          // context canceled
	//fmt.Println(context.Cause(ctx2)) // context canceled

	if ctx2.Err() == context.Canceled {
		t.Fatalf("ctx2 err is not context.Canceled")
	}
	if context.Cause(ctx) != cause {
		t.Fatalf("ctx2 Cause err is not context.Canceled")
	}
}

func TestWithSignals(t *testing.T) {
	// make sure signalContext works like context.Context
	ctx, cancel := WithSignals(context.Background())
	var ch = make(chan int)
	for i := 0; i < 10; i++ {
		go func(_i int) {
			<-ctx.Done()
			fmt.Println(_i)
			ch <- _i
		}(i)
	}

	cancel()
	<-ctx.Done()
	for i := 0; i < 10; i++ {
		select {
		case v := <-ch:
			fmt.Println(v)
		}
	}
}

func TestSignalsCtx(t *testing.T) {
	// make sure signalContext works like context.Context

	ctx := SignalsCtx(context.Background())
	var c = make(chan int)
	context.AfterFunc(ctx, func() {
		fmt.Println(ctx.Err())
		close(c)
	})

	s := ctx.(*signalContext)
	s.signalChan <- syscall.SIGINT

	<-c
}
