package product

import (
	"fmt"
)

type Product struct {
	ID         string
	Category   string
	Attributes map[string]string
}

type Products struct {
	columnMap *ColumnMap
	data      []*Product
}

func NewProducts(rawProducts []map[string]interface{}, columnMap *ColumnMap, ) *Products {
	return &Products{
		columnMap: columnMap,
		data:      getParsedProducts(rawProducts, columnMap),
	}
}

func (ps *Products) GetProducts() []*Product {
	return ps.data
}

func getParsedProducts(rawProducts []map[string]interface{}, columnMap *ColumnMap) []*Product {
	var res []*Product

	for _, rawProduct := range rawProducts {
		res = append(res, parseProduct(rawProduct, columnMap))
	}
	return res
}

func parseProduct(rawProduct map[string]interface{}, columnMap *ColumnMap) *Product {
	rawAttributes := make(map[string]string, 0)
	product := Product{
		ID:       fmt.Sprintf("%v", rawProduct[columnMap.ProductID]),
		Category: fmt.Sprintf("%v", rawProduct[columnMap.Category]),
		//	Name:     fmt.Sprintf("%v", rawProduct[columnMap.Name]),
	}
	for key, value := range rawProduct {
		if key != columnMap.ProductID && key != columnMap.Category {
			rawAttributes[key] = fmt.Sprintf("%v", value)
		}
	}
	product.Attributes = rawAttributes
	return &product
}
