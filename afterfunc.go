package ctxutils

import (
	"context"
)

func AfterFunc(ctx context.Context, do func()) (func() bool, <-chan struct{}) {
	exit := make(chan struct{})

	closeExitChan := func() {
		close(exit)
	}

	// don't need once because closeExitChan after stop and closeExitChan before do cant call both
	//var once sync.Once
	//closeExitChan := func() {
	//	once.Do(func() {
	//		close(exit)
	//	})
	//}

	stop := context.AfterFunc(ctx, func() {
		defer closeExitChan()
		do()
	})
	return func() bool {
		if stop() {
			closeExitChan()
			return true
		}
		return false
	}, exit
}
