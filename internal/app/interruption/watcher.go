package interruption

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WatchForInterruption(cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
}
