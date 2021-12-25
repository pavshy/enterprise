package app_name

import (
	"context"
	"gitlab.com/dsp6/cloudfront/pkg/apps/runner"
	"gitlab.com/dsp6/cloudfront/pkg/services/repositories"
	"time"

	"github.com/sirupsen/logrus"
)

type AppName struct {
	name string
	runner.Runner

	logger *logrus.Entry

	repoName repositories.Name
}

func New(repoName repositories.Name) *AppName {
	var appName = "app_name"

	cf := &AppName{
		name: appName,

		logger: logrus.WithField("app", appName),

		repoName: repoName,
	}

	cf.Runner = runner.NewRunner(appName, cf.run)
	return cf
}

func (cf *AppName) Shutdown(ctx context.Context) {
	cf.Runner.Shutdown(ctx)
	cf.logger.Infof("%s service exited", cf.name)
}

func (cf *AppName) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cf.process()
		}
	}
}

func (cf *AppName) process() {
	_ = cf.repoName.SaveData([]struct{}{})
	_, _ = cf.repoName.GetData()
}
