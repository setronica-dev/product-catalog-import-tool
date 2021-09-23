package reports

import (
	"fmt"
	"go.uber.org/dig"
	"strings"
	"ts/adapters"
	"ts/productImport/mapping"
	"ts/utils"
)

const (
	categoryKey  = "Category"
	productIdKey = "ID"
)

type ReportsHandler struct {
	Handler     adapters.HandlerInterface
	Header      *ReportLabels
	FileManager *adapters.FileManager
}

type Deps struct {
	dig.In
	Handler     adapters.HandlerInterface
	FileManager *adapters.FileManager
	Mapping     mapping.MappingHandlerInterface
}

func NewReportsHandler(deps Deps) *ReportsHandler {
	h := deps.Handler
	h.Init(adapters.CSV)
	m := deps.Mapping.Parse()

	return &ReportsHandler{
		Handler:     h,
		FileManager: deps.FileManager,
		Header:      initFirstRaw(m),
	}
}

func (r *ReportsHandler) WriteReport(
	feedPath string,
	isError bool,
	report []Report,
	sourceData []map[string]interface{},
	columnMap map[string]string,
) string {
	var data [][]string
	path := r.buildPath(feedPath, isError)
	if isError {
		data = r.buildReportData(report)
	} else {
		data = r.buildSuccessData(report, sourceData, columnMap)
	}
	r.Handler.Write(path, data)
	return path
}

func (r *ReportsHandler) buildPath(feedPath string, isError bool) string {
	if isError {
		return r.FileManager.BuildFailReportPath(feedPath)
	} else {
		return r.FileManager.BuildSuccessReportPath(feedPath)
	}
}

func (r *ReportsHandler) buildReportData(report []Report) [][]string {
	var res [][]string
	res = append(res, r.buildHeaderRaw())

	for _, item := range report {
		recordItem := r.buildRaw(item)
		res = append(res, recordItem)
	}
	return res
}

func (r *ReportsHandler) buildSuccessData(report []Report, source []map[string]interface{}, columnMap map[string]string) [][]string {
	var res [][]string
	res = r.buildSuccessMapRaw(source, report, columnMap)
	return res
}

func initFirstRaw(m *mapping.ColumnMap) *ReportLabels {
	labels := ReportLabels{
		CategoryName: "Category Name",
		AttrName:     "Attribute Name*",
		AttrValue:    "Attribute Value*",
		UoM:          "UOM",
		Errors:       "Error Message",
		Description:  "Description",
		DataType:     "Data Type",
		IsMandatory:  "Is Mandatory",
		CodedVal:     "Coded Value",
	}

	if m.ProductID != "" {
		labels.ProductId = m.ProductID
	} else {
		labels.ProductId = "ProductID*"
	}
	if m.Category != "" {
		labels.Category = m.Category
	} else {
		labels.Category = "Category"
	}
	if m.Name != "" {
		labels.Name = m.Name
	} else {
		labels.Name = "Name"
	}

	return &labels
}

func (r *ReportsHandler) buildRaw(item Report) []string {
	return []string{
		item.ProductId,
		item.Name,
		item.Category,
		item.CategoryName,
		item.AttrName,
		item.AttrValue,
		item.UoM,
		strings.Join(item.Errors, " "),
		item.Description,
		item.DataType,
		item.IsMandatory,
		item.CodedVal,
	}
}

func (r *ReportsHandler) buildHeaderRaw() []string {
	header := r.Header
	return []string{
		header.ProductId,
		header.Name,
		header.Category,
		header.CategoryName,
		header.AttrName,
		header.AttrValue,
		header.UoM,
		header.Errors,
		header.Description,
		header.DataType,
		header.IsMandatory,
		header.CodedVal,
	}
}

func (r *ReportsHandler) buildSuccessMapRaw(source []map[string]interface{}, item []Report, columnMap map[string]string) [][]string {
	catKey := categoryKey
	prodIdKey := productIdKey
	var columnMapIndex map[string]string
	if columnMap != nil && len(columnMap) > 0 {
		if val, ok := columnMap[catKey]; ok {
			catKey = val
		}
		if val, ok := columnMap[prodIdKey]; ok {
			prodIdKey = val
		}
		columnMapIndex = utils.RevertMapKeyValue(columnMap)
	}
	// headers from the source but ID and Category go first
	idIndex := 0
	categoryIndex := 1
	sourceRow := source[0]
	slice := make([]string, 2)
	slice[idIndex] = prodIdKey
	slice[categoryIndex] = catKey
	for k := range sourceRow {
		if k != prodIdKey && k != catKey {
			slice = append(slice, k)
		}
	}
	// build index of header cols
	length := len(slice)
	index := make(map[string]int64, length)
	for i, v := range slice {
		index[v] = int64(i)
	}
	//form correct header with Mappings
	headerTs := make([]string, len(slice))
	for i, v := range slice {
		headerTs[i] = utils.GetMapOrDefault(v, columnMapIndex)
	}
	//check all the ontology attribute columns are defined and if not - extend index and headers
	for _, v := range item {
		if _, ok := index[v.AttrName]; !ok {
			headerTs = append(headerTs, v.AttrName)
			index[v.AttrName] = int64(len(headerTs) - 1)
		}
	}
	// reset length
	length = len(headerTs)
	// build final report
	report := make([][]string, 1)
	report[0] = headerTs
	var product []string
	productId := ""
	for _, attribute := range item {
		if productId != attribute.ProductId {
			productId = attribute.ProductId
			//if product does not exist in the report yet
			if len(product) > 0 {
				report = append(report, product)
			}
			product = make([]string, length)
			// first feel data from source file
			for _, sourceItem := range source {
				if sourceItem[prodIdKey] == attribute.ProductId {
					for itemAttr, attrValue := range sourceItem {
						if i, ok := index[itemAttr]; ok {
							product[i] = fmt.Sprintf("%v", attrValue)
						}
					}
					break
				}
			}
			product[categoryIndex] = attribute.Category
			product[idIndex] = attribute.ProductId

			if i, ok := index[attribute.AttrName]; ok {
				product[i] = attribute.AttrValue
			}
		} else {
			if i, ok := index[attribute.AttrName]; ok {
				product[i] = attribute.AttrValue
			}
		}
	}

	if len(product) > 0 {
		report = append(report, product)
	}
	return report
}
