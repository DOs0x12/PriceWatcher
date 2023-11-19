package configer

import (
	"PriceWatcher/internal/entities/config"
	"testing"
)

var ymlData = `
services:
  - price_type: 'bank'
    sending_hours: [12, 17]
    email:
      from: 'test@mail.com'
      password: '123'
      #to: 'test@mail.com'
      to: 'test@mail.com'
      smtp_host: 'smtp.test.com'
      smtp_port: 465

  - price_type: 'marketplace'
    marketplace: 'ozon'
    items:
      'test1': 'test1'
      'test2': 'test2'
    email:
      from: 'test@mail.com'
      password: '123'
      #to: 'test@mail.com'
      to: 'test@mail.com'
      smtp_host: 'smtp.test.com'
      smtp_port: 465`

func TestUnmarshalConfig(t *testing.T) {
	got, err := unmarshalConf([]byte(ymlData))
	if err != nil {
		t.Errorf("The method of getting a config retuns the error: %v", err)

		return
	}

	wantedConf := createWantedConfig()

	gotServCnt := len(got.Services)
	wantServCnt := len(wantedConf.Services)

	if gotServCnt != wantServCnt {
		t.Errorf("Got service count %v, wanted service count %v", gotServCnt, wantServCnt)

		return
	}

	for _, got := range got.Services {
		want := getService(got.PriceType, wantedConf.Services)
		if got.PriceType != want.PriceType {
			t.Errorf("Got %v, wanted %v", got, want)

			return
		}
		if got.Marketplace != want.Marketplace {
			t.Errorf("Got %v, wanted %v", got, want)

			return
		}
	}

}

func createWantedConfig() config.Config {
	conf := config.Config{Services: make([]config.ServiceConf, 0, 2)}

	serv := config.ServiceConf{
		PriceType: "bank",
	}

	conf.Services = append(conf.Services, serv)

	serv = config.ServiceConf{
		PriceType:   "marketplace",
		Marketplace: "ozon",
	}

	conf.Services = append(conf.Services, serv)

	return conf
}

func getService(priceType string, services []config.ServiceConf) config.ServiceConf {
	for _, s := range services {
		if s.PriceType == priceType {
			return s
		}
	}

	return config.ServiceConf{}
}
