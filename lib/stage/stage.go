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

// State 表示应用生命周期所处的阶段。
// 这里的状态是 app 级别的，不是单个服务的状态。
type State uint32

const (
	// StateInit 表示 App 已创建，但还没有开始执行启动流程。
	StateInit State = iota
	// StatePreparing 表示正在执行 beforeStart/Prepare 这类启动前准备工作。
	StatePreparing
	// StateStarting 表示准备工作已完成，开始进入服务运行阶段。
	StateStarting
	// StateRunning 表示所有服务的 Start 已经被发起，且 afterStart 钩子也执行完成。
	StateRunning
	// StateStopping 表示已经进入统一停机流程。
	StateStopping
	// StateStopped 表示正常停机完成。
	StateStopped
	// StateFailed 表示启动或运行过程中出现了首因错误，最终以失败状态退出。
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

// Health 是对外暴露的应用健康快照。
// 这里有意保持为 app 级别的简单视图，方便外部健康检查或管理接口直接使用。
type Health struct {
	// State 是当前应用生命周期状态。
	State State
	// Ready 表示应用是否已经对外宣告“可用”。
	// 当前语义是：所有 afterStart 钩子执行成功后 Ready=true。
	Ready bool
	// Live 表示应用进程当前仍然处于生命周期中，而不是 init/stopped/failed 这些终态。
	Live bool
	// Stopping 表示应用已经进入停机流程。
	Stopping bool
	// Err 只保存触发停机的首因。
	// 停机过程中的附带错误只打日志，不覆盖这里的值。
	Err error
}

// App 负责统一编排多个服务的启动、运行和停止流程。
// 它只管理 app 级生命周期，不介入单个服务自己的内部细节。
type App struct {
	opts *options
	// ctx/cancel 是本次运行的共享上下文。
	// 它主要用于：
	// 1. 传给 Prepare/Start，让服务知道应用是否还在运行
	// 2. 在 shutdown 时统一广播“应用该停了”
	ctx    context.Context
	cancel context.CancelFunc

	// readyChan 在应用对外宣布 ready 时关闭。
	// doneChan 在整个停机流程彻底结束时关闭。
	readyChan chan struct{}
	doneChan  chan struct{}

	// readyOnce/shutdownOnce 用于保证 ready/shutdown 逻辑只执行一次。
	readyOnce    sync.Once
	shutdownOnce sync.Once

	mu sync.RWMutex
	// state 是 app 当前生命周期状态。
	state State
	// err 只记录触发本次退出的首因。
	err error
	// managed 记录当前这次运行中，已经纳入生命周期管理、需要参与 Stop 的服务集合。
	// 这个字段主要用于“部分 Prepare 成功、后续失败”的场景，避免漏清理已占用的资源。
	managed []Server
}

// NewApp 构造一个新的生命周期协调器。
// 默认 stopTimeout 为 5 秒，默认父 context 为 Background。
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

// Run 按固定顺序执行整个应用生命周期：
// 1. beforeStart
// 2. Prepare
// 3. Start（非阻塞发起）
// 4. afterStart
// 5. wait 直到收到退出信号
// 6. shutdown
//
// 这里的设计重点是：
// - Prepare 用来完成真正的“可启动”准备，比如 listen/连接初始化
// - Start 只负责把长时间运行的 serve loop 发出去
// - afterStart 是 app 级别的“启动后”扩展点
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

	// Start 只负责把各个服务的长生命周期循环发出去，它本身不阻塞主协程。
	// runErrCh 用来接收“谁先退出”的信号；真正的退出编排统一在 wait/shutdown 中处理。
	runErrCh := a.start()

	if err := a.afterStart(); err != nil {
		return a.shutdownWithCause(err)
	}

	// afterStart 成功后，应用才真正进入 running。
	a.setStateDirect(StateRunning)

	return a.shutdownWithCause(a.wait(runErrCh))
}

// Stop 提供给外部主动停机使用。
// 调用 Stop 不会立刻返回，而是会等待整个 shutdown 流程走完，或者等待调用方自己的 ctx 超时。
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

// Ready 返回应用 ready 信号。
func (a *App) Ready() <-chan struct{} {
	return a.readyChan
}

// Done 返回应用生命周期结束信号。
// 注意：Done 关闭表示整个停机流程已经走完，不只是收到停机请求。
func (a *App) Done() <-chan struct{} {
	return a.doneChan
}

// State 返回当前状态快照。
func (a *App) State() State {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state
}

// Err 返回本次运行的首因错误。
func (a *App) Err() error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.err
}

// Health 构造一个适合对外暴露的健康快照。
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

// prepare 顺序执行所有实现了 Preparable 的服务。
// 这里之所以维护 prepared 列表，是因为一旦中途 Prepare 失败，
// 之前成功的服务可能已经占用了端口、建立了连接或持有其他资源，
// 后续 shutdown 时必须只清理这批“已经成功进入管理范围”的服务。
func (a *App) prepare() error {
	prepared := make([]Server, 0, len(a.opts.servers))
	for _, srv := range a.opts.servers {
		preparable, ok := srv.(Preparable)
		if !ok {
			continue
		}
		if err := preparable.Prepare(a.ctx); err != nil {
			// 发生部分 Prepare 失败时，只让已经成功 Prepare 的服务参与后续 Stop。
			a.setManaged(prepared)
			return err
		}
		prepared = append(prepared, srv)
	}

	// 所有 Prepare 都成功后，说明本次运行里的全部服务都已经纳入生命周期管理。
	a.setManaged(a.opts.servers)
	return nil
}

// start 并发发起所有服务的 Start。
// 这里返回的 errCh 只表示“服务退出/报错”的信号，不表示“服务已经启动完成”。
func (a *App) start() <-chan error {
	errCh := make(chan error, len(a.opts.servers))
	for _, srv := range a.opts.servers {
		go func(srv Server) {
			errCh <- srv.Start(a.ctx)
		}(srv)
	}

	return errCh
}

// wait 统一等待应用退出的首个触发条件。
// 谁先到来，就以谁作为本次退出的原因：
// 1. 某个服务提前退出并返回 error
// 2. 收到操作系统信号
// 3. app 运行 ctx 被取消
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

// shutdownWithCause 负责统一执行停机流程。
// 这里有两个设计选择：
// 1. 只记录触发停机的首因 cause 到 a.err
// 2. 停机过程中的 beforeStop/Stop/afterStop 错误只打日志，不覆盖首因
//
// 这样做是为了让 Err() 的语义足够稳定：它只表示“为什么退出”，
// 而不是把所有停机过程中的附带问题全部混到一起。
func (a *App) shutdownWithCause(cause error) error {
	a.shutdownOnce.Do(func() {
		a.setStateDirect(StateStopping)
		// 先广播运行 ctx 取消，通知各服务和内部流程“应用准备停机”。
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
		// doneChan 关闭意味着整个 shutdown 流程已经彻底结束。
		close(a.doneChan)
	})

	return cause
}

// fail 只用于“还没有进入统一 shutdown 流程之前”的失败场景。
// 例如 beforeStart 阶段失败，此时还没有需要清理的资源，直接进入 failed 即可。
func (a *App) fail(err error) error {
	a.shutdownOnce.Do(func() {
		a.setErr(err)
		a.setStateDirect(StateFailed)
		close(a.doneChan)
	})
	return err
}

// stopServices 统一停止已经纳入 managed 集合的服务。
// 注意这里不用 a.ctx 直接作为 Stop 的 ctx，
// 因为进入 shutdown 时 a.ctx 已经会被 cancel；
// 如果继续把这个“已取消的 ctx”传给 Stop，很多服务会立刻失败，拿不到优雅退出窗口。
func (a *App) stopServices() error {
	services := a.managedServices()
	if len(services) == 0 {
		return nil
	}

	// Stop 使用单独的 stopCtx，由 app 层统一控制停机预算。
	// stopTimeout == 0 表示不限制优雅停机时间。
	stopCtx := context.Background()
	if a.opts.stopTimeout > 0 {
		var cancel func()
		stopCtx, cancel = context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer cancel()
	}

	errCh := make(chan error, len(services))
	var wg sync.WaitGroup

	// 倒序停止，避免服务之间存在依赖时发生“上游先停，下游还在使用”的问题。
	for i := len(services) - 1; i >= 0; i-- {
		srv := services[i]
		wg.Add(1)
		go func(srv Server) {
			defer wg.Done()
			errCh <- srv.Stop(stopCtx)
		}(srv)
	}

	// done 关闭表示所有 Stop 调用都已经返回。
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

	// 这里仍然聚合 Stop 阶段的错误，但调用方不会把它作为首因暴露出去；
	// 它只用于在 shutdownWithCause 中统一记录日志。
	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// setManaged 记录当前需要参与 Stop 的服务集合。
func (a *App) setManaged(services []Server) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.managed = append([]Server(nil), services...)
}

// managedServices 返回一份拷贝，避免外部直接修改内部服务切片。
func (a *App) managedServices() []Server {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return append([]Server(nil), a.managed...)
}

// setErr 更新首因错误。
func (a *App) setErr(err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.err = err
}

// setState 用于受控状态迁移，要求当前状态必须匹配 expect。
func (a *App) setState(expect, next State) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.state != expect {
		return fmt.Errorf("invalid state transition: %d -> %s", a.state, next)
	}
	a.state = next
	return nil
}

// setStateDirect 用于 shutdown 等明确流程中的直接状态切换。
func (a *App) setStateDirect(next State) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state = next
}

// beforeStart 顺序执行所有启动前 hook。
func (a *App) beforeStart() error {
	for _, fn := range a.opts.beforeStart {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// afterStart 顺序执行所有启动后 hook。
// 当前语义是：
// - 所有服务的 Start 已经被发起
// - app 级启动后逻辑（如注册、通知、ready 标记）在这里执行
// - 只有这些 hook 全部成功后，readyChan 才会关闭
func (a *App) afterStart() error {
	for _, fn := range a.opts.afterStart {
		if err := fn(); err != nil {
			return err
		}
	}
	// readyChan 只在 afterStart 全部成功后关闭，保证对外 ready 语义稳定。
	a.readyOnce.Do(func() {
		close(a.readyChan)
	})
	return nil
}

// beforeStop 顺序执行所有停机前 hook。
// 停机阶段采用“尽量执行完”的策略，因此这里会收集全部错误而不是遇错即停。
func (a *App) beforeStop() error {
	var errs []error
	for _, fn := range a.opts.beforeStop {
		if err := fn(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// afterStop 顺序执行所有停机后 hook。
// 和 beforeStop 一样，停机阶段优先保证流程完整性，错误统一收集返回。
func (a *App) afterStop() error {
	var errs []error
	for _, fn := range a.opts.afterStop {
		if err := fn(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
