package excelH

import (
	"fmt"
	"ts/file/xlsxFile"
	"ts/utils"
)

const (
	alias = "excel"
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

func (h *Adapter) Read(path string) ([][]string, error) {
	res, err := xlsxFile.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to Read  %v: %v", path, err)
	}
	h.setHeader(res[0])
	return res, nil
}

func (h *Adapter) Parse(path string) ([]map[string]interface{}, error) {
	data, err := h.Read(path)
	if err != nil {
		return nil, err
	}
	res, err := utils.RowsToMapRows(data, h.header)
	return res, err
}

func (h *Adapter) Write(filepath string, data [][]string) error {
	return fmt.Errorf("write method is not implemented yet")
}
