package main

import (
	"log"
	"ts/config"
	"ts/di"
	"ts/offerImport"
	"ts/offerItemImport"
	"ts/prepareImport"
	"ts/productImport"
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
	prepareImportHandler *prepareImport.Handler,
	productImportHandler *productImport.ProductImportHandler,
	offerImportHandler *offerImport.OfferImportHandler,
	offerItemImportHandler offerItemImport.OfferItemImportHandlerInterface,
) {
	prepareImportHandler.Run()
	offerImportHandler.RunCSV()
	productImportHandler.RunCSV()
	offerItemImportHandler.Run()
	return
}
