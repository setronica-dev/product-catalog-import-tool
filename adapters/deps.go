package adapters

import (
	"go.uber.org/dig"
	"ts/config"
)

type Deps struct{
	dig.In
	Config *config.Config
}