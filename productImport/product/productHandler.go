package product

import (
	"os"
	"ts/adapters"
	"ts/utils"
)

type ProductHandler struct {
	fileManager *adapters.FileManager
	handler     adapters.HandlerInterface
	ColumnMap   *ColumnMap
}

func NewProductHandler(deps Deps) ProductHandlerInterface {
	m := deps.Mapping.Parse()
	return &ProductHandler{
		fileManager: deps.FileManager,
		handler:     deps.Handler,
		ColumnMap: &ColumnMap{
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
	return res, nil
}

func (ph *ProductHandler) read(filePath string) ([]map[string]interface{}, error) {
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		ph.handler.Init(ph.fileManager.GetFileType(filePath))
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
	currentProductHeader := ph.GetCurrentHeader(sourceData[0])
	return NewProducts(sourceData, currentProductHeader)
}
