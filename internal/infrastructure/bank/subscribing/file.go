package subscribing

import (
	"PriceWatcher/internal/entities/subscribing"
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type subscribersDto struct {
	Subscribers []int64 `yaml:"subscribers"`
}

type SubscribingService struct{}

func (s SubscribingService) GetSubscribers(subscribersFilePath string) (*subscribing.Subscribers, error) {
	subFile, err := os.ReadFile(subscribersFilePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {

		return &subscribing.Subscribers{ChatIDs: make([]int64, 0)}, nil
	}

	if err != nil {
		return nil, err
	}

	return unmarshalConf(subFile)
}

func (s SubscribingService) SaveSubscribers(subs *subscribing.Subscribers,
	subscribersFilePath string) error {
	subsDto := castSubscribers(subs)
	subsData, err := yaml.Marshal(&subsDto)
	if err != nil {
		return err
	}

	err = os.WriteFile(subscribersFilePath, subsData, 0755)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalConf(data []byte) (*subscribing.Subscribers, error) {
	dto := subscribersDto{}
	err := yaml.Unmarshal(data, &dto)
	if err != nil {
		return &subscribing.Subscribers{ChatIDs: make([]int64, 0)}, err
	}

	return castSubDto(dto), nil
}

func castSubDto(subDto subscribersDto) *subscribing.Subscribers {
	return &subscribing.Subscribers{ChatIDs: subDto.Subscribers}
}

func castSubscribers(subs *subscribing.Subscribers) subscribersDto {
	return subscribersDto{Subscribers: subs.ChatIDs}
}
