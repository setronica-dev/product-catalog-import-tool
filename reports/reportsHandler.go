package reports

import (
	"fmt"
	"go.uber.org/dig"
	"strings"
	"ts/adapters"
	"ts/models"
	"ts/utils"
)

const (
	categoryKey  = "Category"
	productIdKey = "ID"
)

type ReportsHandler struct {
	Handler     adapters.HandlerInterface
	Header      *models.ReportLabels
	FileManager *adapters.FileManager
}

type Deps struct {
	dig.In
	Handler     adapters.HandlerInterface
	FileManager *adapters.FileManager
}

func NewReportsHandler(deps Deps) *ReportsHandler {
	h := deps.Handler
	h.Init(adapters.CSV)

	return &ReportsHandler{
		Handler:     h,
		FileManager: deps.FileManager,
		Header:      initFirstRaw(),
	}
}

func (r *ReportsHandler) WriteReport(
	feedPath string,
	isError bool,
	report []models.Report,
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

func (r *ReportsHandler) buildReportData(report []models.Report) [][]string {
	var res [][]string
	res = append(res, r.buildHeaderRaw())

	for _, item := range report {
		recordItem := r.buildRaw(item)
		res = append(res, recordItem)
	}
	return res
}

func (r *ReportsHandler) buildSuccessData(report []models.Report, source []map[string]interface{}, columnMap map[string]string) [][]string {
	var res [][]string
	res = r.buildSuccessMapRaw(source, report, columnMap)
	return res
}

func initFirstRaw() *models.ReportLabels {
	return &models.ReportLabels{
		ProductId:    "ProductID*",
		Name:         "Name",
		Category:     "Category",
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
}

func (r *ReportsHandler) buildRaw(item models.Report) []string {
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
	item := r.Header
	return []string{
		item.ProductId,
		item.Name,
		item.Category,
		item.CategoryName,
		item.AttrName,
		item.AttrValue,
		item.UoM,
		item.Errors,
		item.Description,
		item.DataType,
		item.IsMandatory,
		item.CodedVal,
	}
}

func (r *ReportsHandler) buildSuccessMapRaw(source []map[string]interface{}, item []models.Report, columnMap map[string]string) [][]string {
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
