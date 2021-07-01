package excelH

import (
	"fmt"
	"strings"
	"ts/file/xlsxFile"
	"ts/utils"
)

const (
	alias         = "excel"
	pathDelimiter = "::"
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
	filePath, sheet, err := parsePath(path)
	if err != nil {
		return nil, err
	}
	res, err := xlsxFile.Read(filePath, sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to Read  %v: %v", path, err)
	}
	h.setHeader(res[0])
	return res, nil
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
