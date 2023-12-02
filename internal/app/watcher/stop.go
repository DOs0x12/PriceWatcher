package watcher

import (
	"context"

	"github.com/sirupsen/logrus"
)

func waitJobs(ctx context.Context, finishedJobs <-chan string, jobCount int) {
	<-ctx.Done()

	for i := 0; i < jobCount; i++ {
		jobName := <-finishedJobs
		logrus.Infof("%v: the job is done", jobName)
	}

	logrus.Infof("All jobs is done")
}
