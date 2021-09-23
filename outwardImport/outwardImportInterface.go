package outwardImport

type OutwardImportInterface interface {
	ImportProducts(sourcePath string) (string, error)
	ImportOfferItems(sourcePath string) (string, error)
	WaitForImportComplete(actionID string) (string, error)
	BuildProductAndOffersImportReport(actionID string, importReportFilePath string) error
}
