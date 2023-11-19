package configer

import (
	"PriceWatcher/internal/entities/config"

	"gopkg.in/yaml.v3"
)

func unmarshalConf(data []byte) (config.Config, error) {
	dto := ConfigDto{}
	err := yaml.Unmarshal(data, &dto)
	if err != nil {
		return config.Config{}, err
	}

	return cast(dto), nil
}

func cast(confDto ConfigDto) config.Config {
	serviceCount := len(confDto.Services)
	conf := config.Config{Services: make([]config.ServiceConf, 0, serviceCount)}

	for _, s := range confDto.Services {
		serv := config.ServiceConf{
			SendingHours: s.SendingHours,
			PriceType:    s.PriceType,
			Items:        s.Items,
			Marketplace:  s.Marketplace,
			Email: config.Email{
				From:     s.Email.From,
				Pass:     s.Email.Pass,
				To:       s.Email.To,
				SmtpHost: s.Email.SmtpHost,
				SmtpPort: s.Email.SmtpPort,
			},
		}

		conf.Services = append(conf.Services, serv)
	}

	return conf
}
