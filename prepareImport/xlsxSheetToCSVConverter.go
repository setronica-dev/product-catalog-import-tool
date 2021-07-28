package prepareImport

import (
	"fmt"
	"path/filepath"
	"ts/adapters"
	"ts/file"
)

type XLSXSheetToCSVConverter struct {
	sheet                 string
	headerRowsCount       int
	destinationPath       string
	destinationFileSuffix string
}

func NewXLSXSheetToCSVConverter(
	sheetName string,
	headerRowsCount int,
	destinationPath string,
	suffix string) *XLSXSheetToCSVConverter {
	return &XLSXSheetToCSVConverter{
		sheet:                 sheetName,
		headerRowsCount:       headerRowsCount,
		destinationPath:       destinationPath,
		destinationFileSuffix: suffix,
	}
}

func (c *XLSXSheetToCSVConverter) Convert(filePath string) error {
	destinationPath := c.buildPath(filePath)
	_, err := file.XLSXToCSV(filePath, c.sheet, c.headerRowsCount, destinationPath)
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
