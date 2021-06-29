package mapping

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type mapping struct {
	Map map[string]string
}

func NewMappingHandler() MappingHandlerInterface {
	return &mapping{}
}

func (m *mapping) Init(path string) map[string]string {
	var columnMap map[string]string
	if path != "" {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			m.upload(path)
			columnMap = m.Get()
		}
	}
	return columnMap
}

func (m *mapping) upload(mappingConfigPath string) {
	data, err := ioutil.ReadFile(mappingConfigPath)
	if err != nil {
		log.Fatalf("unable to load mapping file from %s\n%s", mappingConfigPath, err)
	}
	rawMapping := &RawMapping{}
	if err := yaml.Unmarshal(data, rawMapping); err != nil {
		log.Fatalf("unable to unmarshal mapping file %s\n%s", mappingConfigPath, err)
	}

	m.Map = rawMapping.ToConfig()
}

func (m *mapping) Get() map[string]string {
	return m.Map
}
