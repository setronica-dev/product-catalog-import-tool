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
	columnMap *ProductColumnMap
	data      []*Product
}

func NewProducts(rawProducts []map[string]interface{}, columnMap *ProductColumnMap) *Products {
	return &Products{
		columnMap: columnMap,
		data:      getParsedProducts(rawProducts, columnMap),
	}
}

func (ps *Products) GetProducts() []*Product {
	return ps.data
}

func getParsedProducts(rawProducts []map[string]interface{}, columnMap *ProductColumnMap) []*Product {
	var res []*Product

	for _, rawProduct := range rawProducts {
		res = append(res, parseProduct(rawProduct, columnMap))
	}
	return res
}

func parseProduct(rawProduct map[string]interface{}, columnMap *ProductColumnMap) *Product {
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

func (ps *Products) FindProductByID(productID string) *Product {
	if productID == "" {
		return nil
	}
	for _, item := range ps.data {
		if item.ID == productID {
			return item
		}
	}
	return nil
}
