package stage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bird-coder/manyo/pkg/logger"
)

type State uint32

const (
	StateInit State = iota
	StatePreparing
	StateStarting
	StateRunning
	StateStopping
	StateStopped
	StateFailed
)

func (s State) String() string {
	switch s {
	case StateInit:
		return "init"
	case StatePreparing:
		return "preparing"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateStopped:
		return "stopped"
	case StateFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", uint32(s))
	}
}

type Health struct {
	State    State
	Ready    bool
	Live     bool
	Stopping bool
	Err      error
}

type App struct {
	opts   *options
	ctx    context.Context
	cancel context.CancelFunc

	readyChan chan struct{}
	doneChan  chan struct{}

	readyOnce    sync.Once
	shutdownOnce sync.Once

	mu      sync.RWMutex
	state   State
	err     error
	managed []Server
}

func NewApp(opts ...optionFunc) *App {
	o := &options{
		ctx:         context.Background(),
		stopTimeout: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(o)
	}
	ctx, cancel := context.WithCancel(o.ctx)

	return &App{
		opts:      o,
		ctx:       ctx,
		cancel:    cancel,
		readyChan: make(chan struct{}),
		doneChan:  make(chan struct{}),
		state:     StateInit,
	}
}

func (a *App) Run() error {
	if err := a.setState(StateInit, StatePreparing); err != nil {
		return err
	}

	if err := a.beforeStart(); err != nil {
		return a.fail(err)
	}

	if err := a.prepare(); err != nil {
		return a.shutdownWithCause(err)
	}

	if err := a.setState(StatePreparing, StateStarting); err != nil {
		return a.shutdownWithCause(err)
	}

	runErrCh := a.start()

	if err := a.afterStart(); err != nil {
		return a.shutdownWithCause(err)
	}

	a.setStateDirect(StateRunning)

	return a.shutdownWithCause(a.wait(runErrCh))
}

func (a *App) Stop(ctx context.Context) error {
	switch a.State() {
	case StateInit, StateStopped, StateFailed:
		return a.Err()
	}

	a.cancel()

	select {
	case <-a.doneChan:
		return a.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *App) Ready() <-chan struct{} {
	return a.readyChan
}

func (a *App) Done() <-chan struct{} {
	return a.doneChan
}

func (a *App) State() State {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state
}

func (a *App) Err() error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.err
}

func (a *App) Health() Health {
	state := a.State()
	err := a.Err()
	ready := false
	select {
	case <-a.readyChan:
		ready = true
	default:
	}

	return Health{
		State:    state,
		Ready:    ready,
		Live:     state != StateInit && state != StateStopped && state != StateFailed,
		Stopping: state == StateStopping,
		Err:      err,
	}
}

func (a *App) prepare() error {
	prepared := make([]Server, 0, len(a.opts.servers))
	for _, srv := range a.opts.servers {
		preparable, ok := srv.(Preparable)
		if !ok {
			continue
		}
		if err := preparable.Prepare(a.ctx); err != nil {
			a.setManaged(prepared)
			return err
		}
		prepared = append(prepared, srv)
	}

	a.setManaged(a.opts.servers)
	return nil
}

func (a *App) start() <-chan error {
	errCh := make(chan error, len(a.opts.servers))
	for _, srv := range a.opts.servers {
		go func(srv Server) {
			errCh <- srv.Start(a.ctx)
		}(srv)
	}

	return errCh
}

func (a *App) wait(runErrCh <-chan error) error {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(term)

	select {
	case err := <-runErrCh:
		return err
	case sig := <-term:
		logger.Info("received signal: %s", sig.String())
		return nil
	case <-a.ctx.Done():
		if errors.Is(a.ctx.Err(), context.Canceled) {
			return nil
		}
		return a.ctx.Err()
	}
}

func (a *App) shutdownWithCause(cause error) error {
	a.shutdownOnce.Do(func() {
		a.setStateDirect(StateStopping)
		a.cancel()

		if cause != nil {
			a.setErr(cause)
		}

		if err := a.beforeStop(); err != nil {
			logger.Error("before stop hook failed: %v", err)
		}
		if err := a.stopServices(); err != nil {
			logger.Error("stop services failed: %v", err)
		}
		if err := a.afterStop(); err != nil {
			logger.Error("after stop hook failed: %v", err)
		}

		if cause != nil {
			a.setStateDirect(StateFailed)
		} else {
			a.setStateDirect(StateStopped)
		}
		close(a.doneChan)
	})

	return cause
}

func (a *App) fail(err error) error {
	a.shutdownOnce.Do(func() {
		a.setErr(err)
		a.setStateDirect(StateFailed)
		close(a.doneChan)
	})
	return err
}

func (a *App) stopServices() error {
	services := a.managedServices()
	if len(services) == 0 {
		return nil
	}

	// 优雅退出，由app层做超时控制
	stopCtx := context.Background()
	if a.opts.stopTimeout > 0 {
		var cancel func()
		stopCtx, cancel = context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer cancel()
	}

	errCh := make(chan error, len(services))
	var wg sync.WaitGroup

	// 倒序退出，避免服务间依赖问题
	for i := len(services) - 1; i >= 0; i-- {
		srv := services[i]
		wg.Add(1)
		go func(srv Server) {
			defer wg.Done()
			errCh <- srv.Stop(stopCtx)
		}(srv)
	}

	// 标识服务完全停止
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	if a.opts.stopTimeout > 0 {
		select {
		case <-done:
		case <-stopCtx.Done():
			return fmt.Errorf("stop timeout after %s", a.opts.stopTimeout)
		}
	} else {
		<-done
	}

	close(errCh)

	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (a *App) setManaged(services []Server) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.managed = append([]Server(nil), services...)
}

func (a *App) managedServices() []Server {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return append([]Server(nil), a.managed...)
}

func (a *App) setErr(err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.err = err
}

func (a *App) setState(expect, next State) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.state != expect {
		return fmt.Errorf("invalid state transition: %d -> %s", a.state, next)
	}
	a.state = next
	return nil
}

func (a *App) setStateDirect(next State) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state = next
}

/**
 * @description: 服务启动前执行
 * @return {error}
 */
func (a *App) beforeStart() error {
	for _, fn := range a.opts.beforeStart {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

/**
 * @description: 服务启动后执行
 * @return {error}
 */
func (a *App) afterStart() error {
	for _, fn := range a.opts.afterStart {
		if err := fn(); err != nil {
			return err
		}
	}
	a.readyOnce.Do(func() {
		close(a.readyChan)
	})
	return nil
}

/**
 * @description: 服务停止前执行
 * @return {error}
 */
func (a *App) beforeStop() error {
	var errs []error
	for _, fn := range a.opts.beforeStop {
		if err := fn(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

/**
 * @description: 服务停止后执行
 * @return {error}
 */
func (a *App) afterStop() error {
	var errs []error
	for _, fn := range a.opts.afterStop {
		if err := fn(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
