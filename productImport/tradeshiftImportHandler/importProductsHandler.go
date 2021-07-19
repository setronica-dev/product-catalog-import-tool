package tradeshiftImportHandler

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"time"
	"ts/adapters"
	"ts/config"
	"ts/externalAPI/tradeshiftAPI"
)

type TradeshiftHandler struct {
	transport   *tradeshiftAPI.TradeshiftAPI
	filemanager *adapters.FileManager
	handler     adapters.HandlerInterface
}

const (
	inProgressImportState        = "in_progress"
	initImportState              = "init"
	completeImportState          = "complete"
	completeWithErrorImportState = "complete_with_error"
	identifier                   = "ID"
	retriesDelay                 = 1 * time.Second
)

type DepsH struct {
	dig.In
	Config        *config.Config
	TradeshiftAPI *tradeshiftAPI.TradeshiftAPI
	FileManager   *adapters.FileManager
	FilesHandler  adapters.HandlerInterface
}

func NewTradeshiftHandler(deps DepsH) *TradeshiftHandler {
	h := deps.FilesHandler
	h.Init(adapters.TXT)

	return &TradeshiftHandler{
		transport:   deps.TradeshiftAPI,
		filemanager: deps.FileManager,
		handler:     h,
	}
}

func (th *TradeshiftHandler) ImportFeedToTradeshift(
	sourceFeedPath string,
	validationReportPath string) error {
	api := th.transport
	// prepare supplier for import
	err := th.defineSupplierProperties()
	if err != nil {
		return err
	}
	//upload file
	r, err := api.UploadFile(validationReportPath)
	if err != nil {
		return err
	}
	fileID := fmt.Sprintf("%s", r["id"])

	log.Println("Uploaded file with file_id:", fileID)

	//import file
	actionID, err := api.RunImportAction(fmt.Sprintf("%s", fileID))
	if err != nil {
		return err
	}

	state, err := th.waitForImportComplete(actionID)
	if err != nil {
		return err
	}

	//get import report
	err = th.processImportResults(state, sourceFeedPath, actionID)
	return err
}

func (th *TradeshiftHandler) defineSupplierProperties() error {
	r, err := th.transport.GetIdentifier()
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%s", r["name"])
	if id == "" {
		err := th.transport.SetIdentifier(identifier)
		if err != nil {
			return err
		}
	}
	return nil
}

func (th *TradeshiftHandler) getImportState(actionID string) (string, error) {
	res, err := th.transport.GetActionResult(actionID)
	return fmt.Sprintf("%s", res["state"]), err
}

func (th *TradeshiftHandler) waitForImportComplete(actionID string) (string, error) {
	currentState, err := th.getImportState(actionID)
	for currentState == inProgressImportState || currentState == initImportState {
		time.Sleep(retriesDelay)
		currentState, err = th.getImportState(actionID)
		if err != nil {
			return "error", err
		}
	}
	return currentState, nil
}

func isImportCompleted(state string) bool {
	switch state {
	case completeImportState, completeWithErrorImportState:
		return true
	default:
		return false
	}
}

func (th *TradeshiftHandler) processImportResults(importState string, sourceFeedPath string, actionID string) error {
	if !isImportCompleted(importState) {
		return fmt.Errorf("import failed with state %v", importState)
	}
	err := th.WriteReport(actionID, th.filemanager.BuildTradeshiftImportResultsPath(sourceFeedPath))
	if err != nil {
		return err
	}
	if importState == completeImportState {
		log.Printf("IMPORT TO TRADESHIFT WAS FINISHED: see report in %s", th.filemanager.ReportPath)
	} else {
		log.Printf("FAILED TO IMPORT VALID FEED TO TRADESHIFT: see report in %v", th.filemanager.ReportPath)
	}
	return nil
}

func (th *TradeshiftHandler) WriteReport(actionID string, importReportPath string) error {
	report, err := th.transport.GetImportResult(actionID)
	if err != nil {
		return err
	}

	th.handler.Write(importReportPath, [][]string{{report}})
	return nil
}
