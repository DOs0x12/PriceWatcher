package watcher

import (
	"PriceWatcher/internal/app/service"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/sender"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

func ServeWatchers(ctx context.Context,
	wg *sync.WaitGroup,
	configer configer.Configer,
	sender sender.Sender) {
	defer wg.Done()

	config, err := configer.GetConfig()
	if err != nil {
		logrus.Errorf("can not get the config data: %v", err)

		return
	}

	services := config.Services
	servCount := len(services)
	servWG := sync.WaitGroup{}
	servWG.Add(servCount)

	for _, s := range services {
		serv, err := service.NewWatcherService(sender, s)
		if err != nil {
			logrus.Errorf("%v: can not create a watcher service: %v", s.PriceType, err)

			continue
		}

		go watch(ctx, &servWG, serv)
	}

	servWG.Wait()

	logrus.Infof("All the jobs is done")
}
