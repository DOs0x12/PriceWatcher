package watcher

import (
	"context"

	"github.com/sirupsen/logrus"
)

func waitJobs(ctx context.Context, jobNames <-chan string, jobCount int) {
	<-ctx.Done()

	for i := 0; i < jobCount; i++ {
		jobName := <-jobNames
		logrus.Infof("The job %v is done", jobName)
	}

	logrus.Infof("All jobs is done")
}
