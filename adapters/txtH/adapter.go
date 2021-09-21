package txtH

import (
	"fmt"
	"ts/file/txtFile"
)

const (
	alias = "txt"
)

type Adapter struct {
	header string
}

func (h *Adapter) Alias() string {
	return alias
}

func (h *Adapter) GetHeader() []string {
	return []string{
		h.header,
	}
}

func (h *Adapter) setHeader(header string) {
	h.header = header
}

func (h *Adapter) Read(filePath string) ([][]string, error) {
	_, err := txtFile.Read(filePath)
	return nil, err
}

func (h *Adapter) Parse(filePath string) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("parse method is not implemented yet for txt adapter")
}

func (h *Adapter) Write(filepath string, data [][]string) error {
	err := txtFile.Write(filepath, data)
	return err
}
