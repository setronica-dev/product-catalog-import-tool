package reports

import (
	"go.uber.org/dig"
	"ts/adapters"
	"ts/config"
	"ts/productImport/mapping"
	"ts/productImport/product"
)

type ReportsHandler struct {
	Handler           adapters.HandlerInterface
	Header            *ReportLabels
	ColumnMapConfig   *mapping.ColumnMapConfig
	productHandler    product.ProductHandlerInterface
	SuccessResultPath string
	FailResultPath    string
}

type Deps struct {
	dig.In
	Config         *config.Config
	Handler        adapters.HandlerInterface
	Mapping        mapping.MappingHandlerInterface
	ProductHandler product.ProductHandlerInterface
}

func NewReportsHandler(deps Deps) *ReportsHandler {
	h := deps.Handler
	h.Init(adapters.CSV)
	m := deps.Mapping.GetColumnMapConfig()
	conf := deps.Config.ProductCatalog

	return &ReportsHandler{
		Handler:           h,
		productHandler:    deps.ProductHandler,
		ColumnMapConfig:   m,
		Header:            initFailedAttributesReportHeader(m),
		SuccessResultPath: conf.SuccessResultPath,
		FailResultPath:    conf.FailResultPath,
	}
}

func (r *ReportsHandler) WriteReport(
	feedPath string,
	isError bool,
	report []Report,
	sourceData []map[string]interface{},
) string {
	var path string
	if isError {
		path = r.writeFailedReport(report, feedPath)
	} else {
		path = r.writeSuccessReport(report, sourceData, feedPath)
	}

	return path
}

func initFailedAttributesReportHeader(m *mapping.ColumnMapConfig) *ReportLabels {
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
