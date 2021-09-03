package reports

import (
	"fmt"
	"ts/productImport/mapping"
	"ts/utils"
)

type Header struct {
	headerTs    []string
	headerIndex map[string]int64
}

type HeaderBuilder struct {
	sourceRow  map[string]interface{}
	attributes []Report
	columnMap  *mapping.ColumnMapConfig
}

func NewSuccessReportHeaderBuilder(sourceRow map[string]interface{}, reportItems []Report, columnMapConfig *mapping.ColumnMapConfig) *HeaderBuilder {
	return &HeaderBuilder{
		sourceRow:  sourceRow,
		attributes: reportItems,
		columnMap:  columnMapConfig,
	}
}

func buildSuccessReportHeader(sourceRow map[string]interface{}, reportItems []Report, columnMapConfig *mapping.ColumnMapConfig) ([]string, map[string]int64) {
	hb := NewSuccessReportHeaderBuilder(sourceRow, reportItems, columnMapConfig)
	header := hb.buildSuccessReportHeader()
	return header.headerTs, header.headerIndex
}

func (h *HeaderBuilder) buildSuccessReportHeader() *Header {
	sourceHeaderSortedKeys := h.buildSortedHeader()
	header := h.buildHeaderWithMapping(sourceHeaderSortedKeys)

	// Check - are all the ontology attribute columns defined.
	// If not - extend headerIndex and headerTs
	for _, attribute := range h.attributes {
		if _, ok := header.headerIndex[attribute.AttrName]; !ok {
			header.headerTs = append(header.headerTs, attribute.AttrName)
			header.headerIndex[attribute.AttrName] = int64(len(header.headerTs) - 1)
			if attribute.UoM != "" {
				header.headerTs = append(header.headerTs, buildUOMColumnName(attribute.AttrName))
				header.headerIndex[buildUOMColumnName(attribute.AttrName)] = int64(len(header.headerTs) - 1)
			}
		}
	}
	return header
}

func (h *HeaderBuilder) buildSortedHeader() []string {
	requiredKeys := make([]string, 2)
	otherKeys := make([]string, 0)

	for k, _ := range h.sourceRow {
		switch utils.TrimAll(k) {
		case utils.TrimAll(h.columnMap.ProductID):
			requiredKeys[idIndex] = k
		case utils.TrimAll(h.columnMap.Category):
			requiredKeys[categoryIndex] = k
		default:
			otherKeys = append(otherKeys, k)
			// add column for UoM
			attr := findAttributeByName(k, h.attributes)
			if attr != nil && attr.UoM != "" {
				otherKeys = append(otherKeys, buildUOMColumnName(k))
			}
		}
	}
	res := make([]string, len(requiredKeys)+len(otherKeys))
	copy(res, requiredKeys)
	copy(res[len(requiredKeys):], otherKeys)
	return res
}

func buildUOMColumnName(attrName string) string {
	return fmt.Sprintf("%v_UOM", attrName)
}

func NewSuccessHeader() *Header {
	headerTs := make([]string, 2)
	headerIndex := make(map[string]int64, 0)

	return &Header{
		headerTs:    headerTs,
		headerIndex: headerIndex,
	}
}

func (h *HeaderBuilder) buildHeaderWithMapping(sourceRow []string) *Header {

	res := NewSuccessHeader()
	res.headerTs[categoryIndex] = tsCategoryKey
	res.headerTs[idIndex] = tsProductIdKey
	for _, sourceColumnName := range sourceRow {
		switch utils.TrimAll(sourceColumnName) {
		case utils.TrimAll(h.columnMap.Category):
			res.headerIndex[sourceColumnName] = categoryIndex
		case utils.TrimAll(h.columnMap.ProductID):
			res.headerIndex[sourceColumnName] = idIndex
		default:
			f := h.columnMap.GetDefaultValueByMapped(sourceColumnName)
			if f != nil {
				res.headerTs = append(res.headerTs, f.DefaultKey)
			} else {
				res.headerTs = append(res.headerTs, sourceColumnName)
			}
			res.headerIndex[sourceColumnName] = int64(len(res.headerTs) - 1)
		}
	}
	return res
}
