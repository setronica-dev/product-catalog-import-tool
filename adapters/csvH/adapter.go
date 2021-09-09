package csvH

import (
	"fmt"
	"ts/file/csvFile"
	"ts/utils"
)

const (
	alias = "csv"
)

type Adapter struct {
	header []string
}

func (h *Adapter) Alias() string {
	return alias
}

func (h *Adapter) GetHeader() []string {
	return h.header
}

func (h *Adapter) setHeader(header []string) {
	h.header = header
}

func (h *Adapter) Read(filePath string) ([][]string, error) {
	result, err := csvFile.Read(filePath)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return [][]string{}, fmt.Errorf("file %v does not contain rows", filePath)
	}
	validResult := getValidRows(result)
	if len(validResult) == 0 {
		return [][]string{}, fmt.Errorf("file %v does not contain valid rows", filePath)
	}
	h.setHeader(validResult[0])
	return validResult, nil
}

func (h *Adapter) Parse(filePath string) ([]map[string]interface{}, error) {
	data, err := h.Read(filePath)
	if err != nil {
		return nil, err
	}
	res, err := utils.RowsToMapRows(data, h.header)
	return res, err
}

func (h *Adapter) Write(filepath string, data [][]string) error {
	err := csvFile.Write(filepath, data)
	if err != nil {
		return err
	}
	return nil
}

func getValidRows(data [][]string) [][]string {
	var res [][]string

	for _, row := range data {
		if !utils.IsEmptyRow(row) {
			res = append(res, row)
		}
	}
	return res
}
