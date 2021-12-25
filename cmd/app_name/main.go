package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"gitlab.com/dsp6/cloudfront/pkg/app_name"
	"gitlab.com/dsp6/cloudfront/pkg/apps"
	"gitlab.com/dsp6/cloudfront/pkg/config"
	"gitlab.com/dsp6/cloudfront/pkg/services/connections"
	"gitlab.com/dsp6/cloudfront/pkg/services/repositories"
)

func main() {
	var systemConfigPath = flag.String(
		"system-config",
		"/etc/cloudfront/conf.yaml",
		"Path to configuration file",
	)
	flag.Parse()

	sysConf, err := config.LoadFromFile(*systemConfigPath)
	if err != nil {
		logrus.Fatalf("loading config error: %v", err)
	}

	conn, err := connections.NewName(sysConf.Connections.Name)
	if err != nil {
		logrus.Fatal(err)
	}

	repo := repositories.NewName(conn)

	apps := application.New()

	appName := app_name.New(repo)
	apps.RegisterServer(appName)

	apps.RegisterOnShutdown(func(ctx context.Context) {
		err := conn.Close()
		if err != nil {
			logrus.Error(err)
		}
	})

	apps.Run()
}
