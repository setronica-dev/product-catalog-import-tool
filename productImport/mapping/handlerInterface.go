package mapping

type MappingHandlerInterface interface {
	Init(mappingConfigPath string) map[string]string
	Get() map[string]string
}
