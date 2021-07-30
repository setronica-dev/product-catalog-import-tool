package offerItemImport

import (
	"path/filepath"
	"ts/adapters"
	"ts/file/csvFile"
	"ts/productImport/mapping"
	"ts/utils"
)

type OfferItemMappingHandler struct {
	columnMap         *mapping.ColumnMap
	sourcePath        string
	successReportPath string
}

func NewOfferItemMappingHandler(deps Deps) OfferItemMappingHandlerInterface {
	conf := deps.Config.OfferItemCatalog
	return &OfferItemMappingHandler{
		columnMap:         deps.Mapping.Parse(),
		sourcePath:        conf.SourcePath,
		successReportPath: conf.SuccessResultPath,
	}
}
func (oi *OfferItemMappingHandler) Run() error {
	files := adapters.GetFiles(oi.sourcePath)
	{
		for _, fileName := range files {
			err := oi.applyMapping(fileName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (oi *OfferItemMappingHandler) applyMapping(fileName string) error {
	data, _ := csvFile.Read(filepath.Join(oi.sourcePath, fileName))
	header := oi.buildHeader(data[0])
	data[0] = header
	err := csvFile.Write(filepath.Join(oi.successReportPath, fileName), data)
	if err != nil {
		return err
	}
	return nil
}

func (oi *OfferItemMappingHandler) buildHeader(row []string) []string {
	res := make([]string, len(row))
	for i, value := range row {
		if utils.TrimAll(value) == utils.TrimAll(oi.columnMap.ProductID) {
			res[i] = "ID"
			continue
		}
		res[i] = value
	}
	return res
}
