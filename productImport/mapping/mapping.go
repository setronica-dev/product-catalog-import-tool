package mapping

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

const (
	categoryKey  = "Category" // TS min required column
	productIdKey = "ID"       // TS min required
	nameKey      = "Name"
)

type mapping struct {
	rawMap    map[string]string
	parsedMap *ColumnMap
}

func NewMappingHandler(deps Deps) MappingHandlerInterface {
	rawMap := mapping{}
	rawMap.Init(deps.Config.ProductCatalog.MappingPath)
	return &rawMap
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
	m.rawMap = rawMapping.ToConfig()
}

func (m *mapping) Get() map[string]string {
	return m.rawMap
}

func (m *mapping) Parse() *ColumnMap {
	rawMap := m.Get()
	parsedMap := ColumnMap{}
	if rawMap[categoryKey] != "" {
		parsedMap.Category = rawMap[categoryKey]
	} else {
		parsedMap.Category = categoryKey
	}

	if rawMap[productIdKey] != "" {
		parsedMap.ProductID = rawMap[productIdKey]
	} else {
		parsedMap.ProductID = productIdKey
	}

	if rawMap[nameKey] != "" {
		parsedMap.Name = rawMap[nameKey]
	} else {
		parsedMap.Name = nameKey
	}
	return &parsedMap
}
