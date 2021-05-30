package txtH

import (
	"log"
	"os"
	"strings"
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

func (h *Adapter) Read(filePath string) []map[string]interface{} {
	log.Fatalf("read method is not implemented yet")
	return nil
}

func (h *Adapter) Parse(blobPath string) []map[string]interface{} {
	log.Fatalf("parse method is not implemented yet")
	return nil
}

func (h *Adapter) Write(filepath string, data [][]string) {
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("failed create result file: %v", err)
	}
	defer f.Close()
	for _, items := range data {
		_, err := f.WriteString(buildRaw(items))
		if err != nil {
			log.Printf("failed to write to file: %v", err)
		}
	}
}

func buildRaw(input []string) string {
	r := strings.Join(input, "")
	return r
}
