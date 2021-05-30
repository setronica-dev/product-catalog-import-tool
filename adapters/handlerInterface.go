package adapters

type HandlerInterface interface {
	Init(t FileType)
	GetHeader() []string
	Parse(file string) []map[string]interface{}
	Write(file string, data [][]string)
}
