package internal

import (
	"PriceWatcher/internal/bank"
	"PriceWatcher/internal/entities/subscribing"
	"PriceWatcher/internal/telebot"
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func ServeMetalPrice(ctx context.Context,
	wg *sync.WaitGroup,
	bankService bank.Service,
	jobDone chan<- interface{},
	bot telebot.Telebot,
	subscribers *subscribing.Subscribers) {
	defer wg.Done()

	dur := getWaitTimeWithLogs(bankService, time.Now())

	t := time.NewTimer(dur)
	callChan := t.C
	defer t.Stop()

	defer func() {
		finishJobWithLogs(jobDone)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-callChan:
			go servePriceWithTiming(ctx, bankService, t, bot, subscribers)
		}
	}
}

func servePriceWithTiming(
	ctx context.Context,
	serv bank.Service,
	timer *time.Timer,
	bot telebot.Telebot,
	subscribers *subscribing.Subscribers) {
	msg, _ := serveWithLogs(serv)
	if msg != "" {
		for _, chatID := range subscribers.ChatIDs {
			bot.SendCurrentPrice(msg, chatID)
		}
	}

	now := time.Now()
	dur := perStartWithLogs(serv, now)

	select {
	case <-ctx.Done():
		logrus.Infoln("Interrupting waiting the next period")

		return
	case <-time.After(dur):
	}

	now = time.Now()
	dur = getWaitTimeWithLogs(serv, now)

	timer.Reset(dur)
}

func serveWithLogs(serv bank.Service) (string, string) {
	msg, sub, err := serv.ServePrice()
	if err != nil {
		logrus.Errorf("An error occurs while serving a price: %v", err)

		return "", ""
	}

	logrus.Info("The price is processed")

	return msg, sub
}

// func sendReportWithLogs(sender sender.Sender, msg, sub, servName string, email config.Email) {
// 	err := sender.Send(msg, sub, email)
// 	if err != nil {
// 		logrus.Errorf("%v: cannot send the report: %v", servName, err)
// 	}

// 	logrus.Info(servName + ": a report is sended")
// }

func perStartWithLogs(serv bank.Service, now time.Time) time.Duration {
	dur := serv.PerStartDur(now)
	logrus.Infof("Waiting the start of the next period %v", dur)

	return dur
}

func getWaitTimeWithLogs(serv bank.Service, now time.Time) time.Duration {
	dur := serv.GetWaitTime(now)
	logrus.Infof("Waiting %v", dur)

	return dur
}

func finishJobWithLogs(jobDone chan<- interface{}) {
	jobDone <- struct{}{}
	logrus.Infof("the job is done")
}

// func ServeWatchers(ctx context.Context,
// 	wg *sync.WaitGroup,
// 	configer configer.Configer,
// 	sen sender.Sender,
// 	wr file.WriteReader,
// 	restart <-chan interface{}) {
// 	defer wg.Done()

// 	jobDone := make(chan interface{})
// 	cancel, jobCnt, err := startJobs(configer, sen, wr, jobDone)
// 	if err != nil {
// 		logrus.Error(err)

// 		return
// 	}

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			cancel()
// 			WaithAllJobsDone(jobDone, jobCnt)

// 			return
// 		case <-restart:
// 			logrus.Infof("Restart the watcher jobs")
// 			cancel()
// 			WaithAllJobsDone(jobDone, jobCnt)
// 			cancel, jobCnt, err = startJobs(configer, sen, wr, jobDone)
// 			if err != nil {
// 				logrus.Error(err)

// 				return
// 			}
// 		}
// 	}
// }

// func startJobs(configer configer.Configer,
// 	sen sender.Sender,
// 	wr file.WriteReader,
// 	jobDone chan<- interface{}) (context.CancelFunc, int, error) {
// 	config, err := configer.GetConfig()
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("can not get the config data: %v", err)
// 	}

// 	services := config.Services
// 	jobCtx, cancel := context.WithCancel(context.Background())

// 	for _, s := range services {
// 		serv, err := price.NewPriceService(s, configer, wr)
// 		if err != nil {
// 			logrus.Errorf("%v: can not create a watcher service: %v", s.PriceType, err)

// 			continue
// 		}

// 		go watch(jobCtx, serv, sen, s.Email, jobDone)
// 	}

// 	return cancel, len(config.Services), nil
// }

// func WaithAllJobsDone(jobDone <-chan interface{}, jobCnt int) {
// 	for i := 0; i < jobCnt; i++ {
// 		<-jobDone
// 	}
// }
