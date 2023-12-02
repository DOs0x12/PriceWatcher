package watcher

import (
	"PriceWatcher/internal/app/service"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/sender"
	"context"

	"github.com/sirupsen/logrus"
)

func ServeWatchers(ctx context.Context, configer configer.Configer, sender sender.Sender) {
	config, err := configer.GetConfig()
	if err != nil {
		logrus.Errorf("can not get the config data: %v", err)

		return
	}

	for _, s := range config.Services {
		serv, err := service.NewWatcherService(sender, s)
		if err != nil {
			logrus.Errorf("%v: can not create a watcher service: %v", s.PriceType, err)

			continue
		}

		servCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		go watch(servCtx, serv, s.PriceType)
	}

	<-ctx.Done()
}
