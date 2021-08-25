package reports

import (
	"fmt"
	"path/filepath"
	"ts/adapters"
	"ts/productImport/mapping"
	"ts/utils"
)

const (
	idIndex        = 0
	categoryIndex  = 1
	tsCategoryKey  = "Category"
	tsProductIdKey = "ID"
)

func (r *ReportsHandler) writeSuccessReport(report []Report, sourceData []map[string]interface{}, feedFilePath string) string {
	filePath := filepath.Join(r.SuccessResultPath, buildSuccessFileName(feedFilePath))
	var data [][]string
	data = r.buildSuccessData(report, sourceData)
	r.Handler.Write(filePath, data)
	return filePath
}

func buildSuccessFileName(feedPath string) string {
	ext := filepath.Ext(feedPath)
	sourceFileName := adapters.GetFileName(feedPath)
	return fmt.Sprintf("%v%v", sourceFileName, ext)
}

func (r *ReportsHandler) buildSuccessData(report []Report, source []map[string]interface{}) [][]string {
	var res [][]string
	res = r.buildSuccessMapRaw(source, report)
	return res
}

func buildSuccessReportHeader(sourceRow map[string]interface{}, reportItems []Report, columnMapConfig *mapping.ColumnMapConfig) ([]string, map[string]int64) {
	sourceHeaderKeys := getSortedHeader(sourceRow, columnMapConfig)
	headerTs, headerIndex := buildProductsHeaderPart(sourceHeaderKeys, columnMapConfig)

	//check all the ontology attribute columns are defined and if not - extend headerIndex and headerTs
	for _, reportItem := range reportItems {
		if _, ok := headerIndex[reportItem.AttrName]; !ok {
			headerTs = append(headerTs, reportItem.AttrName)
			headerIndex[reportItem.AttrName] = int64(len(headerTs) - 1)
		}
	}
	return headerTs, headerIndex
}

func getSortedHeader(sourceRow map[string]interface{}, columnMapConfig *mapping.ColumnMapConfig) []string {
	requiredKeys := make([]string, 2)
	otherKeys := make([]string, 0)

	for k, _ := range sourceRow {
		switch utils.TrimAll(k) {
		case utils.TrimAll(columnMapConfig.ProductID):
			requiredKeys[idIndex] = k
		case utils.TrimAll(columnMapConfig.Category):
			requiredKeys[categoryIndex] = k
		default:
			otherKeys = append(otherKeys, k)
		}
	}
	res := make([]string, len(requiredKeys)+len(otherKeys))
	copy(res, requiredKeys)
	copy(res[len(requiredKeys):], otherKeys)
	return res
}

func buildProductsHeaderPart(sourceRow []string, columnMapConfig *mapping.ColumnMapConfig) ([]string, map[string]int64) {
	headerTs := make([]string, 2)
	headerIndex := make(map[string]int64, 0)
	headerTs[categoryIndex] = tsCategoryKey
	headerTs[idIndex] = tsProductIdKey

	for i, sourceColumnName := range sourceRow {
		switch utils.TrimAll(sourceColumnName) {
		case utils.TrimAll(columnMapConfig.Category):
			headerIndex[sourceColumnName] = categoryIndex
		case utils.TrimAll(columnMapConfig.ProductID):
			headerIndex[sourceColumnName] = idIndex
		default:
			f := columnMapConfig.GetDefaultValueByMapped(sourceColumnName)
			if f != nil {
				headerTs = append(headerTs, f.DefaultKey)
			} else {
				headerTs = append(headerTs, sourceColumnName)
			}
			headerIndex[sourceColumnName] = int64(i) //todo check why
		}
	}
	return headerTs, headerIndex
}

/**
* transformation flow
 */
func (r *ReportsHandler) buildSuccessMapRaw(source []map[string]interface{}, reportAttrItems []Report) [][]string {
	tsFormattedHeader, headerIndex := buildSuccessReportHeader(source[0], reportAttrItems, r.ColumnMapConfig)
	// reset length
	headerLength := len(tsFormattedHeader)

	// get parsable products list
	parsedSourceProducts := r.productHandler.InitParsedSourceData(source)
	// build final report
	report := make([][]string, 1)
	report[0] = tsFormattedHeader

	for _, sourceProduct := range parsedSourceProducts.GetProducts() {
		reportItem := make([]string, headerLength)
		reportItem[idIndex] = sourceProduct.ID
		reportItem[categoryIndex] = sourceProduct.Category

		productAttrs := findAttributesByProduct(sourceProduct.ID, reportAttrItems)

		for attrName, attrValue := range sourceProduct.Attributes {
			if i, ok := headerIndex[attrName]; ok {
				attr := findAttributeByName(attrName, productAttrs)
				if attr != nil {
					reportItem[i] = fmt.Sprintf("%v", attr.AttrValue)
				} else {
					reportItem[i] = fmt.Sprintf("%v", attrValue)
				}
			}
		}
		for _, productAttrItem := range productAttrs {
			if i, ok := headerIndex[productAttrItem.AttrName]; ok {
				reportItem[i] = fmt.Sprintf("%v", productAttrItem.AttrValue)
				if productAttrItem.Category != "" {
					reportItem[categoryIndex] = productAttrItem.Category
				}
			}
		}
		if len(reportItem) > 0 {
			report = append(report, reportItem)
		}
	}

	return report
}

func findAttributesByProduct(productID string, attributes []Report) []Report {
	res := make([]Report, 0)
	for _, attr := range attributes {
		if attr.ProductId == productID {
			res = append(res, attr)
		}
	}
	return res
}

func findAttributeByName(attrName string, attributes []Report) *Report {
	for _, attr := range attributes {
		if attr.AttrName == attrName {
			return &attr
		}
	}
	return nil
}
