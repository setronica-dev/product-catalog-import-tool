package tradeshiftAPI

import (
	"go.uber.org/dig"
	"ts/config"
	"ts/externalAPI/rest"
)

type TradeshiftAPI struct {
	TennantId string
	Client    rest.RestClientInterface
}

type Deps struct {
	dig.In
	Connection rest.RestClientInterface
	Config     *config.Config
}

func NewTradeshiftAPI(deps Deps) *TradeshiftAPI {
	return &TradeshiftAPI{
		Client: deps.Connection,
	}
}

