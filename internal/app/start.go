package app

import (
	"PriceWatcher/internal/app/service"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/sender"
	"context"

	"github.com/sirupsen/logrus"
)

func StartWatchers(ctx context.Context, configer configer.Configer, sender sender.Sender) {
	config, err := configer.GetConfig()
	if err != nil {
		logrus.Errorf("can not get the config data: %v", err)
	}

	for s, _ := range config.Services {
		service.NewWatcherService(sender, config)
	}
	//
	// watcher :=
}
