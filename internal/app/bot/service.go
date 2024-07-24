package telebot

import (
	"PriceWatcher/internal/entities/subscribing"
	"PriceWatcher/internal/interfaces"
	"context"
	"fmt"
	"sync"
)

func Start(ctx context.Context,
	wg *sync.WaitGroup,
	bot interfaces.Bot,
	subscribers *subscribing.Subscribers) error {
	defer wg.Done()

	if err := bot.Start(ctx); err != nil {
		return fmt.Errorf("can not start the bot: %v", err)
	}

	go func() {
		<-ctx.Done()
		bot.Stop()
	}()

	return nil
}
