package outwardImport


type OutwardImportInterface interface {
	ImportProductsAndOffers(sourcePath string) (string, error)
	WaitForImportComplete(actionID string) (string, error)
	BuildProductAndOffersImportReport(actionID string, importReportFilePath string) error
}
