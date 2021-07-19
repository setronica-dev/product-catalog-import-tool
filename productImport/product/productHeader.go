package product

import (
	"fmt"
	"ts/utils"
)

type ColumnMap struct {
	ProductID string
	Category  string
	Name      string
}

func (p *ProductHandler) GetCurrentHeader(row map[string]interface{}) *ColumnMap {
	var res ColumnMap
	columnMap := p.columnMap
	for k, _ := range row {
		key := fmt.Sprintf("%v", k)
		switch utils.TrimAll(key) {
		case utils.TrimAll(columnMap.Category):
			res.Category = key
		case utils.TrimAll(columnMap.ProductID):
			res.ProductID = key
		case utils.TrimAll(columnMap.Name):
			res.Name = key
		}
	}
	return &res
}
