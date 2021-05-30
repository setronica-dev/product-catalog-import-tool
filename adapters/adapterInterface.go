package adapters

type AdapterInterface interface {
	Alias() string
	GetHeader() []string
	Read(filePath string) []map[string]interface{}
	Parse(filePath string) []map[string]interface{}
	Write(filepath string, data [][]string)
}
