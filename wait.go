package ctxutils

import (
	"context"
	"time"
)

func WaitWithCtx(ctx context.Context, duration time.Duration) {
	if ctx == nil || ctx.Done() == nil {
		time.Sleep(duration)
		return
	}
	t := time.NewTimer(duration)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return
	case <-t.C:
		return
	}
}
