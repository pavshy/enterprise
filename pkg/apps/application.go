package application

import (
	"context"
	"gitlab.com/dsp6/cloudfront/pkg/apps/runner"
	"gitlab.com/dsp6/cloudfront/pkg/apps/shutdowner"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const ShutdownTimeout = 15 * time.Second

type Application interface {
	Run()
	RegisterServer(server runner.Runner)
	RegisterOnShutdown(f func(ctx context.Context))
	Shutdown()
}

type application struct {
	servers   []runner.Runner
	shutdowns shutdowner.OnShutdown
}

func New() Application {
	return newApplication(shutdowner.New())
}

func newApplication(shutdown shutdowner.OnShutdown) *application {
	return &application{shutdowns: shutdown}
}

func (app *application) Run() {
	for _, server := range app.servers {
		server := server
		go func() {
			server.Run()
		}()
	}
	app.waitForInterruptSignal()
	app.Shutdown()
}

func (app *application) RegisterServer(server runner.Runner) {
	app.servers = append(app.servers, server)
	app.shutdowns.Register(server.Shutdown)
}

func (app *application) waitForInterruptSignal() {
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)
	<-signalChannel
	logrus.Info("shutting down application")
}

func (app *application) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	app.shutdowns.ReverseShutdown(ctx)
	cancel()
}

func (app *application) RegisterOnShutdown(f func(ctx context.Context)) {
	app.shutdowns.Register(f)
}
