package reports

import (
	"go.uber.org/dig"
	"strings"
	"ts/adapters"
	"ts/productImport/mapping"
)

const (
	idIndex       = 0
	categoryIndex = 1
)

type ColumnMap struct {
	Category  string
	ProductID string
	Name      string
}

type ReportsHandler struct {
	Handler     adapters.HandlerInterface
	Header      *ReportLabels
	ColumnMap   *ColumnMap
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
		ColumnMap: &ColumnMap{
			Category:  m.Category,
			ProductID: m.ProductID,
			Name:      m.Name,
		},
		Header: initFirstRaw(m),
	}
}

func (r *ReportsHandler) WriteReport(
	feedPath string,
	isError bool,
	report []Report,
	sourceData []map[string]interface{},
) string {
	var data [][]string
	path := r.buildPath(feedPath, isError)
	if isError {
		data = r.buildFailuresReportData(report)
	} else {
		data = r.buildSuccessData(report, sourceData)
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

func (r *ReportsHandler) buildFailuresReportData(report []Report) [][]string {
	var res [][]string
	res = append(res, r.buildHeaderRaw())

	for _, item := range report {
		recordItem := r.buildRaw(item)
		res = append(res, recordItem)
	}
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
