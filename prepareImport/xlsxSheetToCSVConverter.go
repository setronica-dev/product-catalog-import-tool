package prepareImport

import (
	"fmt"
	"path/filepath"
	"ts/adapters"
	"ts/file"
)

type XLSXSheetToCSVConverter struct {
	sheet                 string
	headerRowsToSkip      int
	destinationPath       string
	destinationFileSuffix string
}

func NewXLSXSheetToCSVConverter(
	sheetName string,
	headerRowsToSkip int,
	destinationPath string,
	suffix string) *XLSXSheetToCSVConverter {
	return &XLSXSheetToCSVConverter{
		sheet:                 sheetName,
		headerRowsToSkip:      headerRowsToSkip,
		destinationPath:       destinationPath,
		destinationFileSuffix: suffix,
	}
}

func (c *XLSXSheetToCSVConverter) Convert(filePath string) error {
	destinationPath := c.buildPath(filePath)
	_, err := file.XLSXToCSV(filePath, c.sheet, c.headerRowsToSkip, destinationPath)
	return err
}

func (c *XLSXSheetToCSVConverter) buildPath(sourceFilePath string) string {
	fileName := adapters.GetFileName(sourceFilePath)

	res := filepath.Join(
		c.destinationPath,
		fmt.Sprintf(
			"%v%v.%v",
			fileName,
			c.destinationFileSuffix,
			adapters.CSV))
	return res
}
