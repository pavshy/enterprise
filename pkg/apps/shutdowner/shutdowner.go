package shutdowner

import (
	"context"
)

type OnShutdown interface {
	Register(func(ctx context.Context))
	Shutdown(context.Context)
	ReverseShutdown(context.Context)
}

type onShutdown []func(ctx context.Context)

func New() OnShutdown {
	return &onShutdown{}
}

func (ons *onShutdown) Register(shutdown func(ctx context.Context)) {
	*ons = append(*ons, shutdown)
}

func (ons *onShutdown) Shutdown(ctx context.Context) {
	done := make(chan struct{}, 1)
	go func() {
		for _, shutdown := range *ons {
			shutdown(ctx)
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	case <-done:
	}
}

func (ons *onShutdown) ReverseShutdown(ctx context.Context) {
	done := make(chan struct{}, 1)
	go func() {
		for i := len(*ons) - 1; i >= 0; i-- {
			(*ons)[i](ctx)
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	case <-done:
	}
}
