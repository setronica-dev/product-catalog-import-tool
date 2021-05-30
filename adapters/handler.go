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

func (c *Handler) Init(t FileType) {
	switch t {
	case XLSX:
		c.Adapter = &excelH.Adapter{}
	case CSV:
		c.Adapter = &csvH.Adapter{}
	case TXT:
		c.Adapter = &txtH.Adapter{}
	default:
		log.Fatal("unsupported source file type (only csv and xlsx are supported)")
	}
}

func (c *Handler) GetHeader() []string {
	return c.header
}

func (c *Handler) Parse(file string) []map[string]interface{} {
	res := c.Adapter.Parse(file)
	c.header = c.Adapter.GetHeader()
	return res
}

func (c *Handler) Write(filepath string, data [][]string) {
	c.Adapter.Write(filepath, data)
}
