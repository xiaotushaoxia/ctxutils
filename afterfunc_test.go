package ctxutils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAfterFunc(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	exit := make(chan struct{})
	context.AfterFunc(ctx, func() {
		for i := 0; i < 3; i++ {
			fmt.Println("倒计时", 3-i)
		}
		close(exit)
	})

	<-exit
}

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

func TestAfterFuncSync2(t *testing.T) {

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
	fmt.Println(stop())
	<-exit
	fmt.Println("exit")

}

func TestAfterFuncSync3(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	defer SyncAfterFuncNoStop(ctx, func() {
		for i := 0; i < 3; i++ {
			fmt.Println("倒计时", 3-i)
			time.Sleep(time.Second)
		}
	})
	fmt.Println("exit")

}
