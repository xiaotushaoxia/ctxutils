# ctxutils

### SignalsCtx Example:

```go
func main() {
    ctx1, cancel := ctxutils.WithSignals(context.Background(), syscall.SIGINT)
    ctx2, cancelCause := ctxutils.WithSignalsCause(context.Background(), syscall.SIGINT)
    
    // auto call cancel when ctx3 is gc, so don't lose reference to ctx3 !!!!!
	// don't lose reference to ctx3 !!!
	// don't lose reference to ctx3 !!!
	// don't lose reference to ctx3 !!!
    // example:
    //   done := ctxutils.SignalsCtx(context.Background(), syscall.SIGINT).Done()
    // no reference to `ctxutils.SignalsCtx(context.Background(), syscall.SIGINT)`, so `done` will be closed immediately
    ctx3 := ctxutils.SignalsCtx(context.Background(), syscall.SIGINT)
    
    // syntactic sugar. same as `ctx4 := ctxutils.SignalsCtx(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)`
    ctx4 := ctxutils.SignalsCtxDefault()
	
	....
}

```

### AfterFunc Example:

> If the caller needs to know whether f is completed,  it must coordinate with f explicitly. ——doc in context.AfterFunc

so I add it

```go
func TestAfterFuncSync1(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        time.Sleep(time.Second)
        cancel()
    }()
    
    stop, exit := AfterFunc(ctx, func() {
        for i := 0; i < 3; i++ {
            fmt.Println("倒计时", 3-i)
            time.Sleep(time.Second)
        }
    })
	
    time.Sleep(time.Second + time.Millisecond*20)
    fmt.Println(stop()) // cant stop anymore
    
    <-exit
    fmt.Println("exit")
}
```
