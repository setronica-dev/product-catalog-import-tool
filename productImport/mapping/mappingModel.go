package mapping

import "ts/utils"

const (
	categoryKey  = "Category" // TS min required column
	productIdKey = "ID"       // TS min required
	nameKey      = "Name"
)

type ColumnMapConfig struct {
	Category     string
	ProductID    string
	Name         string
	OtherColumns []*ColumnItem //todo name
}

type ColumnItem struct {
	DefaultKey string
	MappedKey  string
}

func (m *mapping) NewColumnMap(rawMap map[string]string) *ColumnMapConfig {
	parsedMap := ColumnMapConfig{}

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

	parsedMap.OtherColumns = parseNotRequiredColumns(rawMap)

	return &parsedMap
}

func parseNotRequiredColumns(rawMap map[string]string) []*ColumnItem {
	otherColumns := make([]*ColumnItem, 0)
	for key, value := range rawMap {
		if key != nameKey && key != productIdKey && key != categoryKey {
			otherColumns = append(
				otherColumns,
				&ColumnItem{
					DefaultKey: key,
					MappedKey:  value,
				})
		}
	}
	return otherColumns
}

func (c *ColumnMapConfig) GetDefaultValueByMapped(mappedValue string) *ColumnItem {
	if utils.TrimAll(c.ProductID) == utils.TrimAll(mappedValue) {
		return &ColumnItem{
			DefaultKey: productIdKey,
			MappedKey:  c.ProductID,
		}
	}
	if utils.TrimAll(c.Category) == utils.TrimAll(mappedValue) {
		return &ColumnItem{
			DefaultKey: categoryKey,
			MappedKey:  c.Category,
		}
	}

	if utils.TrimAll(c.Name) == utils.TrimAll(mappedValue) {
		return &ColumnItem{
			DefaultKey: nameKey,
			MappedKey:  c.Name,
		}
	}

	for _, item := range c.OtherColumns {
		if utils.TrimAll(item.MappedKey) == utils.TrimAll(mappedValue) {
			return item
		}
	}
	return nil
}
