package runner

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type Runner interface {
	Name() string
	Run()
	Shutdown(ctx context.Context)
}

type runner struct {
	name string

	wg        sync.WaitGroup
	ctx       context.Context
	cancelCtx context.CancelFunc
	log       logrus.FieldLogger

	task func(context.Context)
}

func NewRunner(name string, task func(context.Context)) Runner {
	ctx, cancelCtx := context.WithCancel(context.Background())
	return &runner{
		name:      name,
		ctx:       ctx,
		cancelCtx: cancelCtx,
		log:       logrus.WithField("component", name),
		task:      task,
	}
}

func (r *runner) Run() {
	r.log.Info("Start worker")
	r.wg.Add(1)
	go r.run()
}

func (r *runner) run() {
	r.log.Info("Run worker")
	defer r.wg.Done()

	r.task(r.ctx)
}

func (r *runner) Shutdown(ctx context.Context) {
	r.log.Info("Stop worker")
	r.cancelCtx()

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		done <- struct{}{}
	}()
	select {
	case <-done:
		r.log.Info("Stopped successfully")
		return
	case <-ctx.Done():
		r.log.Info("Stopped because of context deadline")
		return
	}
}

func (r *runner) Name() string {
	return r.name
}
