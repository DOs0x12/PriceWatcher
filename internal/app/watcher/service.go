package watcher

import (
	"PriceWatcher/internal/app/watcher/price"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/sender"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

func ServeWatchers(ctx context.Context,
	wg *sync.WaitGroup,
	configer configer.Configer,
	sen sender.Sender) {
	defer wg.Done()

	config, err := configer.GetConfig()
	if err != nil {
		logrus.Errorf("can not get the config data: %v", err)

		return
	}

	services := config.Services
	servWG := sync.WaitGroup{}

	for _, s := range services {
		serv, err := price.NewPriceService(s)
		if err != nil {
			logrus.Errorf("%v: can not create a watcher service: %v", s.PriceType, err)

			continue
		}

		servWG.Add(1)
		go watch(ctx, &servWG, serv, sen, s.Email)
	}

	servWG.Wait()

	logrus.Infof("All the jobs is done")
}
