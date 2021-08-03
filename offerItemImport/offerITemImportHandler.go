package offerItemImport

import (
	"fmt"
	"log"
	"path/filepath"
	"ts/adapters"
	"ts/outwardImport"
)

type OfferItemImportHandler struct {
	outwardImportHandler outwardImport.OutwardImportInterface
	sourcePath           string
	reportPath           string
}

func NewOfferItemImportHandler(deps Deps) OfferItemImportHandlerInterface {
	return &OfferItemImportHandler{
		outwardImportHandler: deps.OutwardImportHandler,
		sourcePath:           deps.Config.OfferItemCatalog.SuccessResultPath,
		reportPath:           deps.Config.OfferItemCatalog.ReportPath,
	}
}

func (oi *OfferItemImportHandler) Run() {
	files := adapters.GetFiles(oi.sourcePath)
	if len(files) == 0 {
		log.Printf("Offer Import failed: please, put file with offers into %v", oi.sourcePath)
		return
	}

	for _, fileName := range files {
		err := oi.runImport(fileName)
		if err != nil {
			log.Printf("failed offerItems import: %v", err)
		}
	}
}

func (oi *OfferItemImportHandler) runImport(fileName string) error {
	actionID, err := oi.outwardImportHandler.ImportOfferItems(filepath.Join(oi.sourcePath, fileName))
	if err != nil {
		return err
	}
	state, err := oi.outwardImportHandler.WaitForImportComplete(actionID)
	if err != nil {
		return err
	}
	log.Printf("offer items import was complete with state %v", state)
	err = oi.outwardImportHandler.BuildProductAndOffersImportReport(
		actionID,
		filepath.Join(
			oi.reportPath,
			fmt.Sprintf("report_offer_items_%v", fileName)))
	if err != nil {
		return err
	}
	return nil
}
