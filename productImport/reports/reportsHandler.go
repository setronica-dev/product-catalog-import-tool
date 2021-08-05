package reports

import (
	"go.uber.org/dig"
	"ts/adapters"
	"ts/productImport/mapping"
	"ts/productImport/product"
)

type ReportsHandler struct {
	Handler         adapters.HandlerInterface
	Header          *ReportLabels
	ColumnMapConfig *mapping.ColumnMapConfig
	FileManager     *adapters.FileManager
	productHandler  product.ProductHandlerInterface
}

type Deps struct {
	dig.In
	Handler        adapters.HandlerInterface
	FileManager    *adapters.FileManager
	Mapping        mapping.MappingHandlerInterface
	ProductHandler product.ProductHandlerInterface
}

func NewReportsHandler(deps Deps) *ReportsHandler {
	h := deps.Handler
	h.Init(adapters.CSV)
	m := deps.Mapping.GetColumnMapConfig()

	return &ReportsHandler{
		Handler:         h,
		FileManager:     deps.FileManager,
		productHandler:  deps.ProductHandler,
		ColumnMapConfig: m,
		Header:          initFailuresReportHeader(m),
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


func initFailuresReportHeader(m *mapping.ColumnMapConfig) *ReportLabels {
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