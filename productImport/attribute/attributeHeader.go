package attribute

import (
	"ts/utils"
)

type ColumnMap struct {
	ProductId    string
	Name         string
	Category     string
	CategoryName string
	AttrName     string
	AttrValue    string
	UoM          string
	DataType     string
	Description  string
	IsMandatory  string
	CodedVal     string
}

func (ah *AttributeHandler) getCurrentHeader(row map[string]interface{}) *ColumnMap {
	labels := ah.columnMap
	var res ColumnMap
	for key, _ := range row {
		switch utils.TrimAll(key) {
		case utils.TrimAll(labels.ProductId):
			res.ProductId = key
		case utils.TrimAll(labels.Name):
			res.Name = key
		case utils.TrimAll(labels.Category):
			res.Category = key
		case utils.TrimAll(labels.CategoryName):
			res.CategoryName = key
		case utils.TrimAll(labels.AttrName):
			res.AttrName = key
		case utils.TrimAll(labels.AttrValue):
			res.AttrValue = key
		case utils.TrimAll(labels.IsMandatory):
			res.IsMandatory = key
		case utils.TrimAll(labels.DataType):
			res.DataType = key
		case utils.TrimAll(labels.CodedVal):
			res.CodedVal = key
		case utils.TrimAll(labels.Description):
			res.Description = key
		case utils.TrimAll(labels.UoM):
			res.UoM = key
		}
	}
	return &res
}
