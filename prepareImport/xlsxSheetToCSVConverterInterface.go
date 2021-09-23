package prepareImport

type XLSXSheetToCSVConvertInterface interface {
	Convert(filePath string) error
}
