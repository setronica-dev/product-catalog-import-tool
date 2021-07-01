package prepareImport

import (
	"fmt"
	"path/filepath"
	"ts/adapters"
	"ts/file"
)

type XLSXSheetToCSVConverter struct {
	sheet                 string
	destinationPath       string
	destinationFileSuffix string
}

func NewXLSXSheetToCSVConverter(
	sheet string,
	destinationPath string,
	suffix string) *XLSXSheetToCSVConverter {
	return &XLSXSheetToCSVConverter{
		sheet:                 sheet,
		destinationPath:       destinationPath,
		destinationFileSuffix: suffix,
	}
}

func (c *XLSXSheetToCSVConverter) Convert(filePath string) error {
	destinationPath := c.buildPath(filePath)
	err := file.XLSXToCSV(filePath, c.sheet, destinationPath)
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
