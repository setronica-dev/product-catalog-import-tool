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
)

type TradeshiftHandler struct {
	transport            *tradeshiftAPI.TradeshiftAPI
	filemanager          *adapters.FileManager
	handler              adapters.HandlerInterface
	outwardImportHandler outwardImport.OutwardImportInterface
	reportPath           string
}

const (
	completeImportState = "complete"
)

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

	actionID, err := th.outwardImportHandler.ImportProductsAndOffers(validationReportPath)
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
	if state == completeImportState {
		log.Printf("IMPORT TO TRADESHIFT WAS FINISHED: see report in %s", th.filemanager.ReportPath)
	} else {
		log.Printf("FAILED TO IMPORT VALID FEED TO TRADESHIFT: see report in %v", th.filemanager.ReportPath)
	}
	return nil
}
