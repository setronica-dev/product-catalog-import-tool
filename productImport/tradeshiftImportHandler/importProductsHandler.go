package tradeshiftImportHandler

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"path/filepath"
	"ts/adapters"
	"ts/config"
	"ts/externalAPI/tradeshiftAPI"
	"ts/outwardImport"
	"ts/outwardImport/importToTradeshift"
)

type TradeshiftHandler struct {
	transport            *tradeshiftAPI.TradeshiftAPI
	filemanager          *adapters.FileManager
	handler              adapters.HandlerInterface
	outwardImportHandler outwardImport.OutwardImportInterface
	reportPath           string
}


type DepsH struct {
	dig.In
	Config               *config.Config
	TradeshiftAPI        *tradeshiftAPI.TradeshiftAPI
	FileManager          *adapters.FileManager
	FilesHandler         adapters.HandlerInterface
	OutwardImportHandler outwardImport.OutwardImportInterface
}

func NewTradeshiftHandler(deps DepsH) *TradeshiftHandler {
	h := deps.FilesHandler
	h.Init(adapters.TXT)

	return &TradeshiftHandler{
		transport:            deps.TradeshiftAPI,
		filemanager:          deps.FileManager,
		handler:              h,
		outwardImportHandler: deps.OutwardImportHandler,
		reportPath:           deps.Config.ProductCatalog.ReportPath,
	}
}

func (th *TradeshiftHandler) ImportFeedToTradeshift(
	validationReportPath string) error {

	actionID, err := th.outwardImportHandler.ImportProducts(validationReportPath)
	if err != nil {
		return err
	}
	state, err := th.outwardImportHandler.WaitForImportComplete(actionID)
	if err != nil {
		return err
	}

	fileName := adapters.GetFileName(validationReportPath)

	err = th.outwardImportHandler.BuildProductAndOffersImportReport(
		actionID,
		filepath.Join(
			th.reportPath,
			fmt.Sprintf("report_products_%v", fileName)))
	if err != nil {
		return err
	}
	switch state {
	case importToTradeshift.CompleteImportState:
		log.Println("Product import has been finished successfully")
	case importToTradeshift.CompleteWithErrorImportState:
		log.Printf("Product import has been finished with errors. See report here '%v'", th.filemanager.ReportPath)
	default:
		log.Printf("Product import has been failed. See report here '%v'", th.filemanager.ReportPath)
	}
	return nil
}
