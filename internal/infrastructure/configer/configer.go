package configer

import (
	"PriceWatcher/internal/entities/config"
	"os"
)

// var (
// 	path      = "config.yml"
// )

type ConfigDto struct {
	Services []ServiceConfDto `yaml:"services"`
}

type ServiceConfDto struct {
	PriceType    string            `yaml:"price_type"`
	SendingHours []int             `yaml:"sending_hours"`
	Items        map[string]string `yaml:"items"`
	Marketplace  string            `yaml:"marketplace"`
	Email        EmailDto          `yaml:"email"`
}

type EmailDto struct {
	From     string `yaml:"from"`
	Pass     string `yaml:"password"`
	To       string `yaml:"to"`
	SmtpHost string `yaml:"smtp_host"`
	SmtpPort int    `yaml:"smtp_port"`
}

type Configer struct {
	path string
}

func (c Configer) GetConfig() (config.Config, error) {
	confFile, err := os.ReadFile(c.path)
	if err != nil {
		return config.Config{}, err
	}

	return unmarshalConf(confFile)
}

func NewConfiger(path string) Configer {
	return Configer{path: path}
}
