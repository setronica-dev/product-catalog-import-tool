package file

import (
	"fmt"
	"ts/file/csvFile"
	"ts/file/xlsxFile"
)

func XLSXToCSV(sourceFilePath string, destinationFilePath string) error {
	data, err := xlsxFile.Read(sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to comvert XLSX to csv: %v", err)
	}
	err = csvFile.Write(destinationFilePath, data)
	if err != nil {
		return fmt.Errorf("failed to comvert XLSX to csv: %v", err)
	}
	return nil
}
