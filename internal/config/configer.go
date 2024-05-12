package config

import (
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

func (c Configer) GetConfig() (Config, error) {
	confFile, err := os.ReadFile(c.path)
	if err != nil {
		return Config{}, err
	}

	return unmarshalConf(confFile)
}

func NewConfiger(path string) Configer {
	return Configer{path: path}
}

func unmarshalConf(data []byte) (Config, error) {
	dto := ConfigDto{}
	err := yaml.Unmarshal(data, &dto)
	if err != nil {
		return Config{}, err
	}

	return cast(dto), nil
}

func cast(confDto ConfigDto) Config {
	return Config{BotKey: confDto.BotKey, SendingHours: confDto.SendingHours}
}
