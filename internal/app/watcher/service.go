package watcher

import (
	"PriceWatcher/internal/app/watcher/price"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/file"
	"PriceWatcher/internal/interfaces/sender"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

func ServeWatchers(ctx context.Context,
	wg *sync.WaitGroup,
	configer configer.Configer,
	sen sender.Sender,
	wr file.WriteReader,
	restart <-chan interface{}) {
	defer wg.Done()

	cancel := startJobs(configer, sen, wr)

	for {
		select {
		case <-ctx.Done():
			return
		case <-restart:
			logrus.Infof("Restart the watcher jobs")
			cancel()
			cancel = startJobs(configer, sen, wr)
		}
	}
}

func startJobs(configer configer.Configer,
	sen sender.Sender,
	wr file.WriteReader) context.CancelFunc {
	config, err := configer.GetConfig()
	if err != nil {
		logrus.Errorf("can not get the config data: %v", err)

		return nil
	}

	services := config.Services
	jobCtx, cancel := context.WithCancel(context.Background())

	for _, s := range services {
		serv, err := price.NewPriceService(s, configer, wr)
		if err != nil {
			logrus.Errorf("%v: can not create a watcher service: %v", s.PriceType, err)

			continue
		}

		go watch(jobCtx, serv, sen, s.Email)
	}

	return cancel
}
