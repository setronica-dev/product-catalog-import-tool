package xlsxFile

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
	"ts/utils"
)

const (
	pathDelimiter = "::"
)

func Read(path string) ([][]string, error) {
	res, err := processFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read  %v: %v", path, err)
	}
	return res, err
}

func processFile(path string) ([][]string, error) {
	filePath, sheetName, err := parsePath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	xlFile, err := openFile(filePath)
	if err != nil || xlFile == nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	rows := getSheetData(xlFile.Sheets, sheetName)

	return rows, nil
}

func parsePath(path string) (string, string, error) {
	res := strings.SplitN(path, pathDelimiter, 2)
	if len(res[0]) == 0 {
		return "", "", fmt.Errorf("file path is not defined")
	}
	if len(res) == 1 || res[1] == "" {
		return "", "", fmt.Errorf("sheet name is not defined")
	}
	return res[0], res[1], nil
}

func openFile(filePath string) (*xlsx.File, error) {
	if len(filePath) == 0 {
		return nil, fmt.Errorf("file path is not detected: %v", filePath)
	}
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading the file %v: %v", filePath, err)
	}
	if xlFile == nil {
		return nil, fmt.Errorf("failed to open file %v: file is empty", filePath)
	}
	return xlFile, nil
}

func getSheetData(sheets []*xlsx.Sheet, name string) [][]string {
	rows := make([][]string, 0, 0)

	for _, sheet := range sheets {
		if utils.TrimAll(sheet.Name) == utils.TrimAll(name) {
			rows = processSheetData(sheet)
			break
		}
	}
	return rows
}

func processSheetData(sheet *xlsx.Sheet) [][]string {
	rows := sheet.Rows
	res := make([][]string, 0, 0)
	for _, row := range rows {
		singleMap := make([]string, 0)
		for _, cell := range row.Cells {
			text := cell.String()
			singleMap = append(singleMap, text)
		}
		res = append(res, singleMap)
	}
	return res
}

func Write(filepath string, data [][]string) error {
	return fmt.Errorf("write method is not implemented yet")
}
