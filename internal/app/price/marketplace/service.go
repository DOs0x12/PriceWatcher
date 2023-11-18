package marketplace

import (
	priceTime "PriceWatcher/internal/app/price/time"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	"PriceWatcher/internal/entities/config"
	"PriceWatcher/internal/interfaces/file"
	"PriceWatcher/internal/interfaces/requester"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type Service struct {
	wr       file.WriteReader
	req      requester.Requester
	ext      extractor.Extractor
	analyser analyser.Analyser
	conf     config.Config
}

func NewService(
	wr file.WriteReader,
	req requester.Requester,
	ext extractor.Extractor,
	analyser analyser.Analyser,
	conf config.Config) Service {

	return Service{
		wr:       wr,
		req:      req,
		ext:      ext,
		analyser: analyser,
		conf:     conf,
	}
}

func (s Service) ServePrice() (message, subject string, err error) {

	itemPrices := s.conf.Items

	curPrices, err := s.wr.Read()
	if err != nil {
		return "", "", err
	}

	crossedKeys := make([]string, len(curPrices))

	for k := range curPrices {
		if _, ok := itemPrices[k]; ok {
			crossedKeys = append(crossedKeys, k)

			continue
		}

		delete(curPrices, k)
	}

	for k := range itemPrices {
		if !slices.Contains(crossedKeys, k) {
			curPrices[k] = 0.0
		}
	}

	priceType := capitalize(s.conf.PriceType)
	sub := fmt.Sprintf("Цена на товар %v", priceType)
	messages := make([]string, 0)
	cnt := 0

	for k, v := range itemPrices {
		cnt++

		err := s.serveItemPrice(curPrices, &messages, k, v)
		if err != nil {
			return "", "", err
		}

		if len(itemPrices) > cnt {
			waitNextCallWithRand()
		}
	}

	err = s.wr.Write(curPrices)
	if err != nil {
		return "", "", err
	}

	return strings.Join(messages, "\n"), sub, nil
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (s Service) GetWaitTime(now time.Time) time.Duration {
	variation := 10
	randDur := priceTime.RandomMin(variation)
	callTime := getCallTime(now)

	return priceTime.GetWaitDurWithRandomComp(now, callTime, randDur)
}

func (Service) PerStartDur(now time.Time) time.Duration {
	return priceTime.PerStartDur(now, priceTime.HalfHour)
}

func waitNextCallWithRand() {
	permanentComp := 60
	randComp := 120

	dur := time.Duration(permanentComp+rand.Intn(randComp)) * time.Second
	time.Sleep(dur)
}

func (s Service) serveItemPrice(curPrices map[string]float64, messages *[]string, name, address string) error {
	response, err := s.req.RequestPage(address)
	if err != nil {
		return fmt.Errorf("cannot get a page with the current price: %w", err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return fmt.Errorf("cannot extract the price from the body: %w", err)
	}

	changed, up, amount := s.analyser.AnalysePrice(price, float32(curPrices[name]))

	if changed && !up {
		msg := fmt.Sprintf("Цена на %v %v уменьшилась на %.2fр. Текущая цена: %.2fр\n", name, address, amount, price)
		*messages = append(*messages, msg)

		curPrices[name] = float64(price)

		logrus.Info("The item price has been changed. A report is sended")

		return nil
	}

	if curPrices[name] == 0.0 {
		curPrices[name] = float64(price)
	}

	logrus.Info("The item price has been not changed")

	return nil
}
