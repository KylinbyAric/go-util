package goroutine

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

func NewGoroutineWaitGroup() *GoroutineWG {
	return &GoroutineWG{}
}

// GoroutineWG 用于多个子协程组以wait方式执行，主协程阻塞方式等待子协程组，无panic
type GoroutineWG struct {
	m   sync.Mutex
	fns []func()
}

func (g *GoroutineWG) AddFn(fn func()) {
	g.m.Lock()
	defer g.m.Unlock()

	if len(g.fns) == 0 {
		g.fns = make([]func(), 0)
	}
	g.fns = append(g.fns, fn)
}

func (g *GoroutineWG) RunWithWait(ctx context.Context) {
	g.m.Lock()
	defer g.m.Unlock()
	if len(g.fns) == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(g.fns))
	for _, fn := range g.fns {
		runFn := fn
		// 无panic协程
		NewGoroutine(ctx, GoroutineWgName, func(ctx context.Context) {
			defer wg.Done()
			runFn()
		})
	}
	wg.Wait()
	g.fns = nil
	return
}

func (g *GoroutineWG) RunWithWaitTimeOut(ctx context.Context, duration time.Duration) {
	g.m.Lock()
	defer g.m.Unlock()
	if len(g.fns) == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(g.fns))
	for _, fn := range g.fns {
		runFn := fn
		// 无panic协程
		NewGoroutine(ctx, GoroutineWgName, func(ctx context.Context) {
			defer wg.Done()
			runFn()
		})
	}
	waitTimeout(&wg, duration)
	g.fns = nil
	return
}

// time.after + channel 方式 实现超时控制
// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// NewGoroutine 无panic的方式启动goroutine
func NewGoroutine(ctx context.Context, name string, fn func(ctx context.Context)) {
	if name == "" {
		name = DefaultGoroutineName
	}
	spawned := context.WithValue(ctx, GoroutineNameContextKey, name)
	go func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				err := fmt.Errorf(string(debug.Stack()))
				funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				fmt.Println(ctx, "PANIC:%s\nfuncName:%s\n%s", err, funcName, debug.Stack())
			}
		}()
		fmt.Println(ctx, "msg=start a goroutine||goroutine_name=%s")
		fn(spawned)
	}()
}
