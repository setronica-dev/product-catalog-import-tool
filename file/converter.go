package file

import (
	"fmt"
	"ts/file/csvFile"
	"ts/file/xlsxFile"
	"ts/utils"
)

func XLSXToCSV(sourceFilePath string, sheet string, headerLinesCount int, destinationFilePath string) (bool, error) {
	data, err := xlsxFile.Read(sourceFilePath, sheet)

	if err != nil {
		return false, fmt.Errorf("failed to convert XLSX to csv: %v", err)
	}
	if len(data) == 0 {
		return false, nil
	}
	clearedData := clearEmptyData(data, headerLinesCount)
	err = csvFile.Write(destinationFilePath, clearedData)
	if err != nil {
		return false, fmt.Errorf("failed to convert XLSX to csv: %v", err)
	}

	return true, nil
}

func clearEmptyData(data [][]string, headerLinesCount int) [][]string {
	var res [][]string

	columnIndexes := getValidColumnIndexes(data[headerLinesCount])

	for _, row := range data[headerLinesCount:] {
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
