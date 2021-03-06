package product

import (
	"fmt"
	"os"
	"ts/adapters"
	"ts/utils"
)

type ProductHandler struct {
	handler   adapters.HandlerInterface
	ColumnMap *ProductColumnMap
}

func NewProductHandler(deps Deps) ProductHandlerInterface {
	m := deps.Mapping.GetColumnMapConfig()
	return &ProductHandler{
		handler: deps.Handler,
		ColumnMap: &ProductColumnMap{
			ProductID: m.ProductID,
			Category:  m.Category,
			Name:      m.Name,
		},
	}
}

func (ph *ProductHandler) InitSourceData(sourceFeedPath string) ([]map[string]interface{}, error) {
	sourceData, err := ph.read(sourceFeedPath)
	if err != nil {
		return nil, err
	}
	res := parse(sourceData)
	if len(res) == 0 {
		return nil, fmt.Errorf("empty source data")
	}
	return res, nil
}

func (ph *ProductHandler) read(filePath string) ([]map[string]interface{}, error) {
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		ph.handler.Init(adapters.GetFileType(filePath))
		parsedData := ph.handler.Parse(filePath)
		return parsedData, nil
	}
	return nil, err
}

func parse(sourceData []map[string]interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, row := range sourceData {
		if !utils.IsEmptyMap(row) {
			res = append(res, row)
		}
	}
	return res
}

func (ph *ProductHandler) InitParsedSourceData(sourceData []map[string]interface{}) *Products {
	if len(sourceData) == 0 {
		return NewProducts([]map[string]interface{}{}, ph.ColumnMap)
	}
	currentProductHeader := ph.GetCurrentHeader(sourceData[0])
	return NewProducts(sourceData, currentProductHeader)
}
