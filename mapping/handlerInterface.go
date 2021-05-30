package mapping

type HandlerInterface interface {
	Init(mappingConfigPath string)
	Get() map[string]string
}
