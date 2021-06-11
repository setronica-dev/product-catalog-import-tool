package importHandler

import (
	"go.uber.org/dig"
	"ts/externalAPI/tradeshiftAPI"
)

type Status int

const (
	BuyerNotFound Status = 1
	OfferFound    Status = 2
	OfferCreated  Status = 4
	Failed        Status = 0
)

type Deps struct {
	dig.In
	Transport *tradeshiftAPI.TradeshiftAPI
}

type ImportOfferInterface interface {
	ImportOffer(offerName string, buyerID string) (Status, error)
}
