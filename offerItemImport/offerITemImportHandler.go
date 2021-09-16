package offerItemImport

import (
	"fmt"
	"log"
	"path/filepath"
	"ts/adapters"
	offerItemMapping2 "ts/offerItemImport/offerItemMapping"
	"ts/outwardImport"
	"ts/outwardImport/importToTradeshift"
)

type OfferItemImportHandler struct {
	outwardImportHandler outwardImport.OutwardImportInterface
	offerItemMapping     offerItemMapping2.OfferItemMappingHandlerInterface
	sourcePath           string
	successPath          string
	sentPath             string
	reportPath           string
}

func NewOfferItemImportHandler(deps Deps) OfferItemImportHandlerInterface {
	conf := deps.Config.OfferItemCatalog
	return &OfferItemImportHandler{
		outwardImportHandler: deps.OutwardImportHandler,
		offerItemMapping:     deps.OfferItemMapping,
		sourcePath:           conf.SourcePath,
		sentPath:             conf.SentPath,
		successPath:          conf.SuccessResultPath,
		reportPath:           conf.ReportPath,
	}
}

func (oi *OfferItemImportHandler) Run() {
	log.Println("_________________________________")
	files := adapters.GetFiles(oi.sourcePath)
	if len(files) == 0 {
		log.Printf("Offer Items import failed: please, put file with offer items into %v", oi.sourcePath)
		return
	}

	log.Println("Import Offer Items to Tradeshift has been started")
	for _, fileName := range files {
		err := oi.runOfferItemImportFlow(fileName)
		if err != nil {
			log.Println(err)
		}
	}
}

func (oi *OfferItemImportHandler) runOfferItemImportFlow(fileName string) error {
	err := oi.offerItemMapping.Run()
	if err != nil {
		return fmt.Errorf("failed apply offerItems mapping: %v", err)
	}
	_, err = adapters.MoveToPath(filepath.Join(oi.sourcePath, fileName), oi.sentPath)
	if err != nil {
		return fmt.Errorf("failed offerItems source file replacement: %v", err)
	}

	err = oi.importToTradeshift(fileName)
	if err != nil {
		return fmt.Errorf("failed offerItems import: %v", err)
	}
	return nil
}

func (oi *OfferItemImportHandler) importToTradeshift(fileName string) error {
	actionID, err := oi.outwardImportHandler.ImportOfferItems(filepath.Join(oi.successPath, fileName))
	if err != nil {
		return err
	}
	state, err := oi.outwardImportHandler.WaitForImportComplete(actionID)
	if err != nil {
		return err
	}
	err = oi.outwardImportHandler.BuildProductAndOffersImportReport(
		actionID,
		filepath.Join(
			oi.reportPath,
			fmt.Sprintf("report_offer_items_%v", fileName)))
	if err != nil {
		return err
	}
	switch state {
	case importToTradeshift.CompleteImportState:
		log.Println("Offer Items import has been finished successfully")
	default:
		log.Printf("Offer Items import has been finished with errors. See report here '%v'", oi.reportPath)
	}
	return nil
}
