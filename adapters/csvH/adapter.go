package csvH

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

const (
	alias = "csv"
)

type Adapter struct {
	Header []string
}

func (h *Adapter) Alias() string {
	return alias
}

func (h *Adapter) Parse(blobPath string) []map[string]interface{} {
	parsedData := h.Read(blobPath)
	return parsedData
}

func (h *Adapter) GetHeader() []string {
	return h.Header
}

func (h *Adapter) Read(filePath string) []map[string]interface{} {
	// Load a csv file.
	f, _ := os.Open(filePath)
	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	result, err := r.ReadAll()
	if err != nil {
		log.Fatalf(fmt.Sprintf("The file %v is not valid", filePath))
	}
	parsedData := make([]map[string]interface{}, 0, 0)
	h.Header = result[0]

	for rowCounter, row := range result {

		if rowCounter != 0 {
			var singleMap = make(map[string]interface{})
			for colCounter, col := range row {
				singleMap[h.Header[colCounter]] = col
			}
			if len(singleMap) > 0 {

				parsedData = append(parsedData, singleMap)
			}
		}
	}
	return parsedData
}

func (h *Adapter) Write(filepath string, data [][]string) {
	f, err := os.Create(filepath)

	if err != nil {
		log.Fatalf("failed create result file: %v", err)
	}
	w := csv.NewWriter(f)
	w.Comma = ','
	defer w.Flush()

	for _, item := range data {
		err := w.Write(item)
		if err != nil {
			log.Printf("failed to write to file: %v", err)
		}
	}
}
