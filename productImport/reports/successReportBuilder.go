package reports

import (
	"fmt"
	"ts/utils"
)

const (
	tsCategoryKey  = "Category"
	tsProductIdKey = "ID"
	tsNameKey      = "Name"
)

func (r *ReportsHandler) buildSuccessData(report []Report, source []map[string]interface{}) [][]string {
	var res [][]string
	res = r.buildSuccessMapRaw(source, report)
	return res
}

/**
* transformation iteration
 */
func (r *ReportsHandler) buildSuccessMapRaw(source []map[string]interface{}, reportAttrItems []Report) [][]string {

	tsFormattedHeader, headerIndex := buildSuccessReportHeader(source, reportAttrItems, r.ColumnMap)
	// reset length
	headerLength := len(tsFormattedHeader)

	// build final report
	report := make([][]string, 1)
	report[0] = tsFormattedHeader

	// get parsable products list
	parsedSourceProducts := r.productHandler.InitParsedSourceData(source)

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

func buildSuccessReportHeader(source []map[string]interface{}, reportItems []Report, columnMap *ColumnMap) ([]string, map[string]int64) {
	// headers from the source but ID and Category go first
	sourceRow := source[0]
	sourceOrderedHeader := make([]string, 2)
	sourceOrderedHeader[idIndex] = columnMap.ProductID
	sourceOrderedHeader[categoryIndex] = columnMap.Category
	for k := range sourceRow {
		if utils.TrimAll(k) != utils.TrimAll(columnMap.ProductID) && utils.TrimAll(k) != utils.TrimAll(columnMap.Category) {
			sourceOrderedHeader = append(sourceOrderedHeader, k)
		}
	}
	// build index of header cols
	headerLength := len(sourceOrderedHeader)
	headerIndex := make(map[string]int64, headerLength)
	for i, v := range sourceOrderedHeader {
		headerIndex[v] = int64(i)
	}
	//form correct header with Mappings
	headerTs := make([]string, len(sourceOrderedHeader))
	for i, v := range sourceOrderedHeader {
		switch utils.TrimAll(v) {
		case utils.TrimAll(columnMap.Category):
			headerTs[i] = tsCategoryKey
		case utils.TrimAll(columnMap.ProductID):
			headerTs[i] = tsProductIdKey
		case utils.TrimAll(columnMap.Name):
			headerTs[i] = tsNameKey
		default:
			headerTs[i] = v
		}
	}
	//check all the ontology attribute columns are defined and if not - extend headerIndex and headerTs
	for _, reportItem := range reportItems {
		if _, ok := headerIndex[reportItem.AttrName]; !ok {
			headerTs = append(headerTs, reportItem.AttrName)
			headerIndex[reportItem.AttrName] = int64(len(headerTs) - 1)
		}
	}
	return headerTs, headerIndex
}
