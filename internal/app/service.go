package app

import (
	"GoldPriceGetter/internal/domain"
	"GoldPriceGetter/internal/interfaces"
)

type Service interface {
	HandlePrice()
}

type GoldPriceService struct {
	req    interfaces.Requester
	ext    domain.Extractor
	sender interfaces.Sender
}

func (s *GoldPriceService) HandlePrice() {
	response := s.req.RequestPage()
	price := s.ext.ExtractPrice(response.Body)
	s.sender.Send(price)
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
