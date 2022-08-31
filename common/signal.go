package common

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// 系统安全退出
var ServerStatus *serverStatus

func init() {
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

// 系统状态
type serverStatus struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	exitSignal chan os.Signal
	wg         *sync.WaitGroup
}

//等待操作系统信号，并发出退出信号
func (ss *serverStatus) waitSignal() {
	<-ss.exitSignal
	fmt.Println("系统收到退出信号...")
	ss.cancelFunc()
}

//获取系统全局context
func (ss *serverStatus) Context() context.Context {
	return ss.ctx
}

func (ss *serverStatus) Add(delta int) {
	ss.wg.Add(delta)
}

func (ss *serverStatus) Done() {
	ss.wg.Done()
}

//等待安全退出
func (ss *serverStatus) Wait() {
	ss.wg.Wait()
	fmt.Println("系统已退出.")
}
