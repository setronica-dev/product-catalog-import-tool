package di

import (
	"go.uber.org/dig"
	"log"
	"ts/adapters"
	"ts/config"
	"ts/externalAPI/rest"
	"ts/externalAPI/tradeshiftAPI"
	"ts/offerImport/importHandler"
	"ts/offerImport/offerReader"
	"ts/productImport"
	"ts/productImport/mapping"
	"ts/productImport/ontologyRead"
	"ts/productImport/ontologyValidator"
	"ts/productImport/reports"
	"ts/productImport/tradeshiftImportHandler"
)

type options = []dig.ProvideOption

type entry struct {
	constructor interface{}
	opts        options
}

var diConfig = []entry{
	{constructor: config.Get},
	{constructor: mapping.NewMappingHandler},
	{constructor: adapters.NewFileManager},
	{constructor: adapters.NewHandler},
	{constructor: ontologyRead.NewRulesHandler},
	{constructor: offerReader.NewOfferReader},
	{constructor: ontologyValidator.NewValidator},
	{constructor: reports.NewReportsHandler},
	{constructor: rest.NewRestClient},
	{constructor: tradeshiftAPI.NewTradeshiftAPI},
	{constructor: tradeshiftImportHandler.NewTradeshiftHandler},
	{constructor: productImport.NewProductImportHandler},
	{constructor: importHandler.NewImportOfferHandler},
}

func BuildContainer() *dig.Container {
	container := dig.New()
	for _, entry := range diConfig {
		if err := container.Provide(entry.constructor, entry.opts...); err != nil {
			log.Fatalf("DI provider error\n%s", err)
		}
	}
	return container
}
