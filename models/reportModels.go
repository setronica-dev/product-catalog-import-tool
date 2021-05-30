package models

type Report struct {
	ProductId    string
	Name         string
	Category     string
	CategoryName string
	AttrName     string
	AttrValue    string
	UoM          string
	Errors       []string
	DataType     string
	Description  string
	IsMandatory  string
	CodedVal     string
}

type ReportLabels struct {
	ProductId    string
	Name         string
	Category     string
	CategoryName string
	AttrName     string
	AttrValue    string
	UoM          string
	Errors       string
	DataType     string
	Description  string
	IsMandatory  string
	CodedVal     string
}
