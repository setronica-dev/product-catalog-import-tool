package importToTradeshift

import (
	"fmt"
	"log"
	"time"
)

const (
	inProgressImportState        = "in_progress"
	initImportState              = "init"
	completeImportState          = "complete"
	completeWithErrorImportState = "complete_with_error"
	identifier                   = "ID"
	retriesDelay                 = 1 * time.Second
)

func (ti *TradeshiftImport) ImportProducts(sourcePath string) (string, error) {
	return ti.runSupplierImport(sourcePath, false)
}

func (ti *TradeshiftImport) ImportOfferItems(sourcePath string) (string, error) {
	return ti.runSupplierImport(sourcePath, true)
}

func (ti *TradeshiftImport) runSupplierImport(sourcePath string, isOfferItemsImport bool) (string, error) {

	api := ti.transport
	// prepare supplier for import
	err := ti.defineSupplierProperties()
	if err != nil {
		return "", fmt.Errorf("failed to define supplier properties: %v", err)
	}
	//upload file
	r, err := api.UploadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed file upload %v", err)
	}
	fileID := fmt.Sprintf("%s", r["id"])

	log.Println("Uploaded file with file_id:", fileID)

	//import file
	actionID, err := api.RunImportAction(fmt.Sprintf("%s", fileID), ti.tsConfig.tsCurrency, ti.tsConfig.tsLocale, isOfferItemsImport)
	if err != nil {
		return actionID, fmt.Errorf("failed file import:%v", err)
	}
	return actionID, nil
}

func (ti *TradeshiftImport) defineSupplierProperties() error {
	r, err := ti.transport.GetIdentifier()

	if err != nil {
		return err
	}
	id := fmt.Sprintf("%s", r["name"])
	if id == "" {
		err := ti.transport.SetIdentifier(identifier)
		if err != nil {
			return err
		}
	}
	return nil
}

//---------------------//
func (ti *TradeshiftImport) WaitForImportComplete(actionID string) (string, error) {

	currentState, err := ti.getImportState(actionID)
	for currentState == inProgressImportState || currentState == initImportState {
		time.Sleep(retriesDelay)
		currentState, err = ti.getImportState(actionID)
		if err != nil {
			return "error", err
		}
	}
	if !isImportCompleted(currentState) {
		return currentState, fmt.Errorf("import failed with state %v", currentState)
	}

	return currentState, nil
}

func (ti *TradeshiftImport) getImportState(actionID string) (string, error) {
	res, err := ti.transport.GetActionResult(actionID)
	return fmt.Sprintf("%s", res["state"]), err
}

func isImportCompleted(state string) bool {
	switch state {
	case completeImportState, completeWithErrorImportState:
		return true
	default:
		return false
	}
}
