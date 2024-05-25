package price

import (
	"PriceWatcher/internal/bank"
	"PriceWatcher/internal/config"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CurrentPriceComm struct{}

func NewPriceCommand() CurrentPriceComm {
	return CurrentPriceComm{}
}

func (c CurrentPriceComm) GetCurrentPrices() string {
	configer := GetConfiger()
	conf, err := configer.GetConfig()
	if err != nil {
		logrus.Error("%w", err)
	}

	bankService := bank.NewService(bank.BankRequester{}, bank.NewPriceExtractor(`([0-9])*(&nbsp;)([0-9])*,([0-9])*`, "div"), conf)
	msg, _, err := bankService.ServePrice()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	return msg
}

func GetConfiger() config.Configer {
	configPath := "config.yml"

	return config.NewConfiger(configPath)
}
