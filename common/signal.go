package common

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var ServerStatus *serverStatus

func init() { //主协程
	ctx, cancelFunc := context.WithCancel(context.Background())
	ServerStatus = &serverStatus{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		exitSignal: make(chan os.Signal, 1),
		wg:         &sync.WaitGroup{},
	}
	signal.Notify(ServerStatus.exitSignal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go ServerStatus.waitSignal()
}

type serverStatus struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	exitSignal chan os.Signal
	wg         *sync.WaitGroup
}

func (ss *serverStatus) waitSignal() {
	<-ss.exitSignal
	fmt.Println("系统收到退出信号...")
	ss.cancelFunc()
	ss.wg.Wait()
	fmt.Println("系统已退出.")
}

func (ss *serverStatus) Context() context.Context {
	return ss.ctx
}

func (ss *serverStatus) Add(delta int) {
	ss.wg.Add(delta)
}

func (ss *serverStatus) Done() {
	ss.wg.Done()
}
