package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/interfaces"
	"fmt"
)

type Service interface {
	HandlePrice() error
}

type GoldPriceService struct {
	req    interfaces.Requester
	ext    domain.Extractor
	sender interfaces.Sender
}

func (s *GoldPriceService) HandlePrice() error {
	response, err := s.req.RequestPage()
	if err != nil {
		serviceErr := "cannot get a page with the current price of gold: %w"

		return fmt.Errorf(serviceErr, err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	s.sender.Send(price)
	if err != nil {
		serviceErr := "cannot exttract the gold price from the body: %w"

		return fmt.Errorf(serviceErr, err)
	}

	return nil
}

func NewGoldPriceService(
	req interfaces.Requester,
	ext domain.Extractor,
	sender interfaces.Sender) *GoldPriceService {

	serv := GoldPriceService{
		req:    req,
		ext:    ext,
		sender: sender,
	}

	return &serv
}
