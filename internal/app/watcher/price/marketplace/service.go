package marketplace

import (
	priceTime "PriceWatcher/internal/app/watcher/price/time"
	"PriceWatcher/internal/domain/price/analyser"
	"PriceWatcher/internal/domain/price/extractor"
	priceEnt "PriceWatcher/internal/entities/price"
	"PriceWatcher/internal/interfaces/configer"
	"PriceWatcher/internal/interfaces/file"
	"PriceWatcher/internal/interfaces/requester"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
)

type Service struct {
	wr       file.WriteReader
	req      requester.Requester
	ext      extractor.Extractor
	analyser analyser.Analyser
	configer configer.Configer
	name     string
}

func NewService(
	wr file.WriteReader,
	req requester.Requester,
	ext extractor.Extractor,
	analyser analyser.Analyser,
	configer configer.Configer,
	name string) Service {

	return Service{
		wr:       wr,
		req:      req,
		ext:      ext,
		analyser: analyser,
		configer: configer,
		name:     name,
	}
}

func (s Service) ServePrice() (message, subject string, err error) {
	conf, err := s.configer.GetMarketplaceConfig(s.name)
	if err != nil {
		return "", "", fmt.Errorf("cannot get the config for a service with the name %v: %w", s.name, err)
	}

	itemPrices := conf.Items

	s.wr.Lock()

	curPrices, err := s.wr.ReadPrices()
	if err != nil {
		return "", "", err
	}

	for k := range curPrices {
		if _, ok := itemPrices[k]; !ok {
			delete(curPrices, k)
		}
	}

	priceType := capitalize(conf.PriceType)
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

	err = s.wr.WritePrices(curPrices)
	if err != nil {
		return "", "", err
	}

	s.wr.Unlock()

	return strings.Join(messages, "\n"), sub, nil
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (s Service) GetWaitTime(now time.Time) time.Duration {
	variation := 600
	randDur := priceTime.RandomSec(variation)
	callTime := getCallTime(now)

	return getWaitDurWithRandomComp(now, callTime, randDur)
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

func (s Service) serveItemPrice(curPrices map[string]priceEnt.ItemPrice, messages *[]string, name, address string) error {
	serveName := s.GetName()

	response, err := s.req.RequestPage(address)
	if err != nil {
		return fmt.Errorf("%v: cannot get a page with the current price: %w", serveName, err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return fmt.Errorf("%v: cannot extract the price from the body: %w", serveName, err)
	}

	curPrice, ok := curPrices[name]
	if !ok {
		curPrice = priceEnt.ItemPrice{Address: address, Price: 0.0}
		curPrices[name] = curPrice
	}

	changed, up, amount := s.analyser.AnalysePrice(price, float32(curPrice.Price))

	if changed && !up {
		msg := fmt.Sprintf("Цена на %v %v уменьшилась на %.2fр. Текущая цена: %.2fр\n", name, address, amount, price)
		*messages = append(*messages, msg)

		curPrice.Price = float64(price)

		logrus.Infof("%v: the item price has been changed. A report is sended", serveName)

		return nil
	}

	if curPrice.Price == 0.0 {
		curPrice.Price = float64(price)
		curPrices[name] = curPrice
	}

	logrus.Infof("%v: the item price has been not changed", serveName)

	return nil
}

func (s Service) GetName() string {
	return s.name
}
