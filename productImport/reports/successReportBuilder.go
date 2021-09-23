package reports

import (
	"fmt"
	"path/filepath"
	"ts/adapters"
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

		// setup product attributes from source
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
		//actualize values from attributes file
		for _, productAttrItem := range productAttrs {
			if i, ok := headerIndex[productAttrItem.AttrName]; ok {
				reportItem[i] = fmt.Sprintf("%v", productAttrItem.AttrValue)
				if productAttrItem.UoM != "" {
					if j, ok := headerIndex[buildUOMColumnName(productAttrItem.Name)]; ok {
						reportItem[j] = productAttrItem.UoM
					}
				}
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
