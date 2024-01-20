package watcher

import (
	"PriceWatcher/internal/app/watcher/price"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/file"
	"PriceWatcher/internal/interfaces/sender"
	"context"
	"fmt"
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

	jobDone := make(chan interface{})
	cancel, jobCnt, err := startJobs(configer, sen, wr, jobDone)
	if err != nil {
		logrus.Error(err)

		return
	}

	for {
		select {
		case <-ctx.Done():
			cancel()
			WaithAllJobsDone(jobDone, jobCnt)

			return
		case <-restart:
			logrus.Infof("Restart the watcher jobs")
			cancel()
			WaithAllJobsDone(jobDone, jobCnt)
			cancel, jobCnt, err = startJobs(configer, sen, wr, jobDone)
			if err != nil {
				logrus.Error(err)

				return
			}
		}
	}
}

func startJobs(configer configer.Configer,
	sen sender.Sender,
	wr file.WriteReader,
	jobDone chan<- interface{}) (context.CancelFunc, int, error) {
	config, err := configer.GetConfig()
	if err != nil {
		return nil, 0, fmt.Errorf("can not get the config data: %v", err)
	}

	services := config.Services
	jobCtx, cancel := context.WithCancel(context.Background())

	for _, s := range services {
		serv, err := price.NewPriceService(s, configer, wr)
		if err != nil {
			logrus.Errorf("%v: can not create a watcher service: %v", s.PriceType, err)

			continue
		}

		go watch(jobCtx, serv, sen, s.Email, jobDone)
	}

	return cancel, len(config.Services), nil
}

func WaithAllJobsDone(jobDone <-chan interface{}, jobCnt int) {
	for i := 0; i < jobCnt; i++ {
		<-jobDone
	}
}
