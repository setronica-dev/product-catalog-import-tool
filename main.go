package main

import (
	"log"
	"ts/adapters"
	"ts/config"
	"ts/di"
	"ts/externalAPI/rest"
	"ts/externalAPI/tradeshiftAPI"
	"ts/offerImport"
	"ts/offerImport/importHandler"
	"ts/productImport"
	"ts/productImport/mapping"
	"ts/productImport/ontologyRead"
	"ts/productImport/ontologyValidator"
	"ts/productImport/reports"
	"ts/productImport/tradeshiftImportHandler"
)

func main() {
	config.Init()

	// init di build container
	bc := di.BuildContainer()

	// inject stuff and start service
	if err := bc.Invoke(start); err != nil {
		log.Fatalf("instantiation error\n%s", err)
	}
}

func start(
	config *config.Config,
	mapHandler mapping.MappingHandlerInterface,
	rulesHandler *ontologyRead.RulesHandler,
	handler adapters.HandlerInterface,
	validator ontologyValidator.ValidatorInterface,
	reports *reports.ReportsHandler,
	fileManager *adapters.FileManager,
	rest rest.RestClientInterface,
	tradeshiftAPI *tradeshiftAPI.TradeshiftAPI,
	tradeshiftProductHandler *tradeshiftImportHandler.TradeshiftHandler,
	tradeshiftOfferHandler importHandler.ImportOfferInterface,
	productImportHandler *productImport.ProductImportHandler,
	offerImportHandler *offerImport.OfferImportHandler,
) {
	offerImportHandler.Run()
	productImportHandler.Run()
	return
}
