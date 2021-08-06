package importToTradeshift

import (
	"fmt"
	"os"
	"time"
)

const (
	inProgressImportState        = "in_progress"
	initImportState              = "init"
	CompleteImportState          = "complete"
	CompleteWithErrorImportState = "complete_with_error"
	identifier                   = "ID"
	retriesDelay                 = 1 * time.Second
)

func (ti *TradeshiftImport) ImportProducts(sourceFilePath string) (string, error) {
	return ti.runSupplierImport(sourceFilePath, false)
}

func (ti *TradeshiftImport) ImportOfferItems(sourceFilePath string) (string, error) {
	return ti.runSupplierImport(sourceFilePath, true)
}

func (ti *TradeshiftImport) runSupplierImport(sourceFilePath string, isOfferItemsImport bool) (string, error) {
	_, err := os.Stat(sourceFilePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("file for import %v not found", sourceFilePath)
	}
	api := ti.transport

	// prepare supplier for import
	err = ti.defineSupplierProperties()
	if err != nil {
		return "", fmt.Errorf("failed to define supplier properties: %v", err)
	}
	//upload file
	r, err := api.UploadFile(sourceFilePath)
	if err != nil {
		return "", fmt.Errorf("failed file upload %v", err)
	}
	fileID := fmt.Sprintf("%s", r["id"])

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
	return currentState, nil
}

func (ti *TradeshiftImport) getImportState(actionID string) (string, error) {
	res, err := ti.transport.GetActionResult(actionID)
	return fmt.Sprintf("%s", res["state"]), err
}
