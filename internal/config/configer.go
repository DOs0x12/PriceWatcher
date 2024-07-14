package config

import (
	"PriceWatcher/internal/entities/config"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigDto struct {
	BotKey       string `yaml:"botKey"`
	SendingHours []int  `yaml:"sending_hours"`
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

func unmarshalConf(data []byte) (config.Config, error) {
	dto := ConfigDto{}
	err := yaml.Unmarshal(data, &dto)
	if err != nil {
		return config.Config{}, err
	}

	return cast(dto), nil
}

func cast(confDto ConfigDto) config.Config {
	return config.Config{BotKey: confDto.BotKey, SendingHours: confDto.SendingHours}
}
