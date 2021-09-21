package adapters

import (
	"log"
	"ts/adapters/csvH"
	"ts/adapters/excelH"
	"ts/adapters/txtH"
)

type FileType string

const (
	CSV  FileType = "csv"
	XLSX FileType = "xlsx"
	TXT  FileType = "txt"
)

type Handler struct {
	Adapter   AdapterInterface
	Delimiter rune
	header    []string
	LineChan  chan interface{}
}

func NewHandler() HandlerInterface {
	return &Handler{}
}

func (h *Handler) Init(t FileType) {
	switch t {
	case XLSX:
		h.Adapter = &excelH.Adapter{}
	case CSV:
		h.Adapter = &csvH.Adapter{}
	case TXT:
		h.Adapter = &txtH.Adapter{}
	default:
		log.Fatal("unsupported source file type (only csv and xlsx are supported)")
	}
}

func (h *Handler) GetHeader() []string {
	return h.header
}

func (h *Handler) Write(filepath string, data [][]string) {
	err := h.Adapter.Write(filepath, data)
	if err != nil {
		log.Fatalf("failed to write %v file: %v", filepath, err)
	}
}

func (h *Handler) Parse(filePath string) []map[string]interface{} {
	res, err := h.Adapter.Parse(filePath)
	h.header = h.Adapter.GetHeader()
	if err != nil {
		log.Fatalf("failed to Read file %v: %v", filePath, err)
	}
	return res
}
