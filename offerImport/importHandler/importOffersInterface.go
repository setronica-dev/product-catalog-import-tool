package importHandler

import (
	"go.uber.org/dig"
	"ts/config"
	"ts/externalAPI/tradeshiftAPI"
	"ts/offerImport/offerReader"
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
	Config    *config.Config
}

type ImportOfferInterface interface {
	ImportOffers(offers []offerReader.RawOffer)
}
