package mapping

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Mapping struct {
	Map map[string]string
}

func GetHandler() HandlerInterface {
	return &Mapping{}
}

func (m *Mapping) Init(mappingConfigPath string) {
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

func (m *Mapping) Get() map[string]string {
	return m.Map
}
