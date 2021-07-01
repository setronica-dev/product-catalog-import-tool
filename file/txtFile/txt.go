package txtFile

import (
	"fmt"
	"os"
	"strings"
)

func Read(filePath string) ([][]string, error) {
	err := fmt.Errorf("read method is not implemented yet")
	return nil, err
}

func Write(filepath string, data [][]string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed create result file: %v", err)
	}
	defer f.Close()
	for _, items := range data {
		_, err := f.WriteString(buildRaw(items))
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	return nil
}

func buildRaw(input []string) string {
	r := strings.Join(input, "")
	return r
}
