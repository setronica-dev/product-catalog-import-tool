package di

import (
	"go.uber.org/dig"
	"log"
	"ts/adapters"
	"ts/config"
	"ts/externalAPI/rest"
	"ts/externalAPI/tradeshiftAPI"
	"ts/offerImport"
	"ts/offerImport/importHandler"
	"ts/offerImport/offerReader"
	"ts/prepareImport"
	"ts/productImport"
	"ts/productImport/attribute"
	"ts/productImport/mapping"
	"ts/productImport/ontologyRead"
	"ts/productImport/ontologyValidator"
	"ts/productImport/product"
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
	{constructor: adapters.NewFileManager},
	{constructor: adapters.NewHandler},

	{constructor: mapping.NewMappingHandler},
	{constructor: ontologyRead.NewRulesHandler},
	{constructor: offerReader.NewOfferReader},
	{constructor: attribute.NewAttributeHandler},
	{constructor: product.NewProductHandler},

	{constructor: ontologyValidator.NewValidator},
	{constructor: reports.NewReportsHandler},

	{constructor: rest.NewRestClient},
	{constructor: tradeshiftAPI.NewTradeshiftAPI},
	{constructor: tradeshiftImportHandler.NewTradeshiftHandler},
	{constructor: importHandler.NewImportOfferHandler},

	{constructor: prepareImport.NewPrepareImportHandler},
	{constructor: productImport.NewProductImportHandler},
	{constructor: offerImport.NewOfferImportHandler},
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
