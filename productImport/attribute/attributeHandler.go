package attribute

import (
	"fmt"
	"os"
	"ts/adapters"
	"ts/utils"
)

type AttributeHandler struct {
	fileManager *adapters.FileManager
	handler     adapters.HandlerInterface
	columnMap   *ColumnMap
}

func NewAttributeHandler(deps Deps) AttributeHandlerInterface {
	h := deps.Report.Header
	return &AttributeHandler{
		fileManager: deps.FileManager,
		handler:     deps.Handler,
		columnMap: &ColumnMap{
			ProductId:    h.ProductId,
			Name:         h.Name,
			Category:     h.Category,
			CategoryName: h.CategoryName,
			AttrName:     h.AttrName,
			AttrValue:    h.AttrValue,
			UoM:          h.UoM,
			DataType:     h.DataType,
			Description:  h.Description,
			IsMandatory:  h.IsMandatory,
			CodedVal:     h.CodedVal,
		},
	}
}

func (ah *AttributeHandler) InitAttributeData(filePath string) ([]*Attribute, error) {
	reportData := make([]*Attribute, 0)
	reportDataSource, err := ah.readData(filePath)
	if err != nil {
		return nil, err
	}
	currentLabels := ah.getCurrentHeader(reportDataSource[0])
	reportData = ah.parseData(reportDataSource, currentLabels)
	return reportData, nil
}

func (ah *AttributeHandler) readData(filePath string) ([]map[string]interface{}, error) {
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		ah.handler.Init(ah.fileManager.GetFileType(filePath))
		reportDataSource := ah.handler.Parse(filePath)
		return reportDataSource, nil
	}
	return nil, err
}

func (ah *AttributeHandler) parseData(reportDataSource []map[string]interface{}, currentLabels *ColumnMap) []*Attribute {
	reportData := make([]*Attribute, 0)

	for _, line := range reportDataSource {
		if !utils.IsEmptyMap(line) {
			r := Attribute{
				ProductId: fmt.Sprintf("%v", line[currentLabels.ProductId]),
			}
			if line[currentLabels.Name] != nil {
				r.Name = fmt.Sprintf("%v", line[currentLabels.Name])
			}
			if line[currentLabels.Category] != nil {
				r.Category = fmt.Sprintf("%v", line[currentLabels.Category])
			}
			if line[currentLabels.CategoryName] != nil {
				r.CategoryName = fmt.Sprintf("%v", line[currentLabels.CategoryName])
			}
			if line[currentLabels.AttrName] != nil {
				r.AttrName = fmt.Sprintf("%v", line[currentLabels.AttrName])
			}
			if line[currentLabels.AttrValue] != nil {
				r.AttrValue = fmt.Sprintf("%v", line[currentLabels.AttrValue])
			}
			if line[currentLabels.UoM] != "" {
				r.UoM = fmt.Sprintf("%v", line[currentLabels.UoM])
			}
			if line[currentLabels.Description] != "" {
				r.Description = fmt.Sprintf("%v", line[currentLabels.Description])
			}
			if line[currentLabels.DataType] != "" {
				r.DataType = fmt.Sprintf("%v", line[currentLabels.DataType])
			}
			if line[currentLabels.IsMandatory] != "" {
				r.IsMandatory = fmt.Sprintf("%v", line[currentLabels.IsMandatory])
			}
			if line[currentLabels.CodedVal] != "" {
				r.CodedVal = fmt.Sprintf("%v", line[currentLabels.CodedVal])
			}

			reportData = append(reportData, &r)
		}
	}
	return reportData
}
