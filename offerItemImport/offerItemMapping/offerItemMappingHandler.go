package offerItemMapping

import (
	"fmt"
	"path/filepath"
	"ts/adapters"
	"ts/file/csvFile"
	"ts/productImport/mapping"
	"ts/utils"
)

type OfferItemMappingHandler struct {
	columnMap         *mapping.ColumnMapConfig
	sourcePath        string
	successReportPath string
}

func NewOfferItemMappingHandler(deps Deps) OfferItemMappingHandlerInterface {
	conf := deps.Config.OfferItemCatalog
	return &OfferItemMappingHandler{
		columnMap:         deps.Mapping.GetColumnMapConfig(),
		sourcePath:        conf.SourcePath,
		successReportPath: conf.SuccessResultPath,
	}
}
func (oi *OfferItemMappingHandler) Run() error {
	fileNames := adapters.GetFiles(oi.sourcePath)
	if len(fileNames) == 0 {
		return fmt.Errorf("no offer items files found in %v", oi.sourcePath)
	}
	{
		for _, fileName := range fileNames {
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
