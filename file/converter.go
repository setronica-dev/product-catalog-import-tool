package file

import (
	"fmt"
	"ts/file/csvFile"
	"ts/file/xlsxFile"
	"ts/utils"
)

func XLSXToCSV(sourceFilePath string, sheet string, headerRowsToSkip int, destinationFilePath string) (bool, error) {
	data, err := xlsxFile.Read(sourceFilePath, sheet)

	if err != nil {
		return false, fmt.Errorf("failed to read XLSX file: %v", err)
	}
	if len(data) == 0 {
		return false, nil
	}
	clearedData := clearEmptyData(data, headerRowsToSkip)
	if len(clearedData) == 0 {
		return false, nil
	}

	err = csvFile.Write(destinationFilePath, clearedData)
	if err != nil {
		return false, fmt.Errorf("failed to convert XLSX to csv: %v", err)
	}

	return true, nil
}

func clearEmptyData(data [][]string, headerRowsToSkip int) [][]string {
	// will be removed when all sheet's headers will be configurable
	if headerRowsToSkip == 0 {
		return data
	}

	var res [][]string
	columnIndexes := getValidColumnIndexes(data[headerRowsToSkip])

	for _, row := range data[headerRowsToSkip:] {
		if !utils.IsEmptyRow(row) {
			cells := getValidRowCells(row, columnIndexes)
			res = append(res, cells)
		}
	}
	return res
}

func getValidColumnIndexes(data []string) []int {
	columnIndex := make([]int, 0)
	for i, v := range data {
		if v != "" {
			columnIndex = append(columnIndex, i)
		}
	}
	return columnIndex
}

func getValidRowCells(rawData []string, columnIndex []int) []string {
	res := make([]string, 0)
	l := len(rawData)
	for _, i := range columnIndex {
		if i < l {
			res = append(res, rawData[i])
		}
	}
	return res
}
