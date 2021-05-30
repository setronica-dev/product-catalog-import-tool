package excelH

import (
	"container/list"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
)

const (
	alias = "excel"
)

type Adapter struct {
	header list.List
}

func (h *Adapter) Alias() string {
	return alias
}

func (h *Adapter) GetHeader() []string {
	log.Printf("getHeader method is not implemented yet")
	return nil
}

func (h *Adapter) Read(filePath string) []map[string]interface{} {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Println("Error reading the file")
	}

	parsedData := make([]map[string]interface{}, 0, 0)
	headerName := list.New()
	// sheet
	for _, sheet := range xlFile.Sheets {
		// rows
		for rowCounter, row := range sheet.Rows {
			// column
			headerIterator := headerName.Front()
			var singleMap = make(map[string]interface{})

			for _, cell := range row.Cells {
				if rowCounter == 0 {
					text := cell.String()
					headerName.PushBack(text)
				} else {
					text := cell.String()
					singleMap[headerIterator.Value.(string)] = text
					headerIterator = headerIterator.Next()
				}
			}
			if rowCounter != 0 && len(singleMap) > 0 {
				parsedData = append(parsedData, singleMap)
			}
		}
	}
	return parsedData
}

func (h *Adapter) Parse(blobPath string) []map[string]interface{} {
	parsedData := h.Read(blobPath)
	return parsedData
}

func (h *Adapter) Write(filepath string, data [][]string) {
	log.Fatalf("write method is not implemented yet")
}
