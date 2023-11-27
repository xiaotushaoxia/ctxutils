package ctxutils

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
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

	if ctx2.Err() != context.Canceled {
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
			fmt.Println("put", _i)
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

	<-c
}

func TestCancel(t *testing.T) {
	ctx := SignalsCtx(context.Background())

	fmt.Println(context.Cause(ctx))
	fmt.Println(ctx.Err())
}

func TestGC(t *testing.T) {

	for i := 0; i < 20; i++ {
		ccc()
		fmt.Println(runtime.NumGoroutine())
		runtime.GC()
	}

}
func ccc() {
	var ctx context.Context
	//var causeFunc context.CancelCauseFunc
	//var cancel context.CancelFunc
	switch rand.Int() % 3 {
	case 0:
		ctx = SignalsCtx(context.TODO())
	case 1:
		ctx, _ = WithSignalsCause(context.Background())
		//causeFunc(fmt.Errorf("xsdsd"))
	case 2:
		ctx, _ = WithSignals(context.Background())
		//cancel()
	}
	//causeFunc(fmt.Errorf("xxxxxx"))
	//
	//<-ctx.Done()
	fmt.Println(ctx)

	fmt.Println(ctx.Err(), context.Cause(ctx))
}
