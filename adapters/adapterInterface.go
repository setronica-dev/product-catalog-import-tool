package adapters

type AdapterInterface interface {
	Alias() string
	GetHeader() []string
	Read(filePath string) ([][]string, error)
	Parse(filePath string) ([]map[string]interface{}, error)
	Write(filepath string, data [][]string) error
}
