package stage

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"
)

type testService struct {
	name string

	mu      sync.Mutex
	events  *[]string
	startCh chan struct{}
	stopCh  chan struct{}

	prepareErr error
	startErr   error
	stopErr    error

	prepareFn func() error
	startFn   func() error
	stopFn    func() error
}

func newTestService(name string, events *[]string) *testService {
	return &testService{
		name:    name,
		events:  events,
		startCh: make(chan struct{}),
		stopCh:  make(chan struct{}),
	}
}

func (s *testService) record(event string) {
	if s.events == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	*s.events = append(*s.events, event)
}

func (s *testService) Name() string {
	return s.name
}

func (s *testService) Prepare(ctx context.Context) error {
	s.record(s.name + ":prepare")
	if s.prepareFn != nil {
		return s.prepareFn()
	}
	return s.prepareErr
}

func (s *testService) Start(ctx context.Context) error {
	s.record(s.name + ":start")
	close(s.startCh)
	if s.startFn != nil {
		return s.startFn()
	}
	if s.startErr != nil {
		return s.startErr
	}
	<-s.stopCh
	return nil
}

func (s *testService) Stop(ctx context.Context) error {
	s.record(s.name + ":stop")
	select {
	case <-s.stopCh:
	default:
		close(s.stopCh)
	}
	if s.stopFn != nil {
		return s.stopFn()
	}
	return s.stopErr
}

func TestAppRunLifecycle(t *testing.T) {
	var events []string
	svc := newTestService("svc", &events)

	app := NewApp(
		WithServer(svc),
		BeforeStart(func() error {
			events = append(events, "hook:before-start")
			return nil
		}),
		AfterStart(func() error {
			events = append(events, "hook:after-start")
			return nil
		}),
		BeforeStop(func() error {
			events = append(events, "hook:before-stop")
			return nil
		}),
		AfterStop(func() error {
			events = append(events, "hook:after-stop")
			return nil
		}),
	)

	runDone := make(chan error, 1)
	go func() {
		runDone <- app.Run()
	}()

	select {
	case <-app.Ready():
	case <-time.After(2 * time.Second):
		t.Fatal("app did not become ready")
	}

	if got := app.State(); got != StateRunning {
		t.Fatalf("expected running state, got %s", got)
	}

	if health := app.Health(); !health.Ready || !health.Live || health.Stopping {
		t.Fatalf("unexpected health: %+v", health)
	}

	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("stop returned error: %v", err)
	}

	select {
	case err := <-runDone:
		if err != nil {
			t.Fatalf("run returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("run did not exit")
	}

	select {
	case <-app.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("app done was not closed")
	}

	if got := app.State(); got != StateStopped {
		t.Fatalf("expected stopped state, got %s", got)
	}

	wantOrder := []string{
		"hook:before-start",
		"svc:prepare",
		"hook:after-start",
	}
	if !slices.Equal(events[:3], wantOrder) {
		t.Fatalf("unexpected startup order: %v", events)
	}

	if !slices.Contains(events, "svc:start") || !slices.Contains(events, "svc:stop") {
		t.Fatalf("service start/stop events missing: %v", events)
	}
}

func TestAppPrepareFailureStopsPreparedServices(t *testing.T) {
	var events []string
	first := newTestService("first", &events)
	second := newTestService("second", &events)
	second.prepareErr = errors.New("prepare failed")

	app := NewApp(WithServers(first, second))

	err := app.Run()
	if err == nil || !errors.Is(err, second.prepareErr) {
		t.Fatalf("expected prepare error, got %v", err)
	}

	if !slices.Contains(events, "first:stop") {
		t.Fatalf("prepared service was not stopped on prepare failure: %v", events)
	}

	if got := app.State(); got != StateFailed {
		t.Fatalf("expected failed state, got %s", got)
	}
}

func TestAppRunReturnsRuntimeError(t *testing.T) {
	var events []string
	runErr := errors.New("boom")
	svc := newTestService("svc", &events)
	svc.startErr = runErr

	app := NewApp(WithServer(svc))

	err := app.Run()
	if err == nil || !errors.Is(err, runErr) {
		t.Fatalf("expected runtime error, got %v", err)
	}

	if got := app.State(); got != StateFailed {
		t.Fatalf("expected failed state, got %s", got)
	}

	if !slices.Contains(events, "svc:stop") {
		t.Fatalf("service stop was not called after runtime error: %v", events)
	}
}
