package configer

import (
	"PriceWatcher/internal/entities/config"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	debugPath = "../config.yml"
	path      = "config.yml"
)

type ConfigDto struct {
	SendingHours []int    `yaml:"sending_hours"`
	Email        EmailDto `yaml:"email"`
}

type EmailDto struct {
	From     string `yaml:"from"`
	Pass     string `yaml:"password"`
	To       string `yaml:"to"`
	SmtpHost string `yaml:"smtp_host"`
	SmtpPort int    `yaml:"smtp_port"`
}

type Configer struct{}

func (Configer) GetConfig() (config.Config, error) {
	confFile, err := os.ReadFile(debugPath)
	//confFile, err := os.ReadFile(path)
	if err != nil {
		return config.Config{}, err
	}

	dto := ConfigDto{Email: EmailDto{}}
	err = yaml.Unmarshal(confFile, dto)
	if err != nil {
		return config.Config{}, err
	}

	return cast(dto), nil
}

func cast(confDto ConfigDto) config.Config {
	return config.Config{
		SendingHours: confDto.SendingHours,
		Email: config.Email{
			From:     confDto.Email.From,
			Pass:     confDto.Email.Pass,
			To:       confDto.Email.To,
			SmtpHost: confDto.Email.SmtpHost,
			SmtpPort: confDto.Email.SmtpPort,
		},
	}
}
