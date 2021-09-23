package csvFile

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
)

func Read(filePath string) ([][]string, error) {
	// Load a csv file.
	f, _ := os.Open(filePath)
	defer f.Close()
	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	result, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf( /**/ "the file %v is not valid", filePath)
	}
	return result, nil
}

func Write(filepath string, data [][]string) error {
	f, err := os.Create(filepath)
	if f == nil {
		return fmt.Errorf("failed create file %v", filepath)
	}
	defer f.Close()

	if err != nil {
		return fmt.Errorf("failed create result file: %v", err)
	}
	w := csv.NewWriter(f)
	w.Comma = ','
	defer w.Flush()

	for _, item := range data {
		err := w.Write(item)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	return nil
}
