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
func (r *ReportsHandler) buildSuccessMapRaw(source []map[string]interface{}, reportItems []Report) [][]string {

	tsFormattedHeader, headerIndex := buildSuccessReportHeader(source, reportItems, r.ColumnMap)
	// reset length
	headerLength := len(tsFormattedHeader)

	// build final report
	report := make([][]string, 1)
	report[0] = tsFormattedHeader
	product := make([]string, 0)
	productId := ""

	sourceColumnMap := getSourceKeys(source[0], r.ColumnMap)
	for _, attribute := range reportItems {
		if productId != attribute.ProductId {
			productId = attribute.ProductId
			//if product does not exist in the report yet
			if len(product) > 0 {
				report = append(report, product)
			}

			// first feel data from source file:
			// - find product for attribute
			var foundProduct map[string]interface{}
			for _, sourceItem := range source {
				if sourceItem[sourceColumnMap.ProductID] == attribute.ProductId {
					foundProduct = sourceItem
					break
				}
			}

			// - fill attributes
			product = make([]string, headerLength)
			for itemAttr, attrValue := range foundProduct {
				if i, ok := headerIndex[itemAttr]; ok {
					product[i] = fmt.Sprintf("%v", attrValue)
				}
			}

			// fill main info for product:
			if attribute.Category == "" {
				product[categoryIndex] = fmt.Sprintf("%v", foundProduct[sourceColumnMap.Category])
			} else {
				product[categoryIndex] = attribute.Category
			}

			if attribute.Name == "" {
				product[idIndex] = fmt.Sprintf("%v", foundProduct[sourceColumnMap.ProductID])
			} else {
				product[idIndex] = attribute.ProductId
			}

			// fill fixed attribute value:
			if i, ok := headerIndex[attribute.AttrName]; ok {
				product[i] = attribute.AttrValue
			}
		} else {
			if i, ok := headerIndex[attribute.AttrName]; ok {
				product[i] = attribute.AttrValue
			}
		}
	}

	if len(product) > 0 {
		report = append(report, product)
	}
	return report
}

func getSourceKeys(sourceRow map[string]interface{}, columnMap *ColumnMap) *ColumnMap {
	var sourceColumnMap ColumnMap
	for k, _ := range sourceRow {
		key := fmt.Sprintf("%v", k)
		switch utils.TrimAll(key) {
		case utils.TrimAll(columnMap.Category):
			sourceColumnMap.Category = key
		case utils.TrimAll(columnMap.ProductID):
			sourceColumnMap.ProductID = key
		case utils.TrimAll(columnMap.Name):
			sourceColumnMap.Name = key
		}
	}
	return &sourceColumnMap
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
