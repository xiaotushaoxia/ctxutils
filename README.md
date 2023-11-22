# ctxutils

### SignalsCtx Example:

```go
func main() {
    ctx, cancel := ctxutil.WithSignals()
    ctx, cancelCause := ctxutil.WithSignalsCause()
    ctx := ctxutil.SignalsCtx()
    // use ctx ...
}
```

### AfterFunc Example:

doc in context.AfterFunc says "If the caller needs to know whether f is completed,  it must coordinate with f explicitly."
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