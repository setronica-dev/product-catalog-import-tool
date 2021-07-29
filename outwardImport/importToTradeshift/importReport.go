package importToTradeshift


func (ti *TradeshiftImport) BuildProductAndOffersImportReport(actionID string, importReportFilePath string) error {
	report, err := ti.transport.GetImportResult(actionID)
	if err != nil {
		return err
	}

	ti.handler.Write(importReportFilePath, [][]string{{report}})
	return nil
}
