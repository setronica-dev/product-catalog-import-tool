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
			r := &Attribute{
				ProductId:    fmt.Sprintf("%v", line[currentLabels.ProductId]),
				Name:         fmt.Sprintf("%v", line[currentLabels.Name]),
				Category:     fmt.Sprintf("%v", line[currentLabels.Category]),
				CategoryName: fmt.Sprintf("%v", line[currentLabels.CategoryName]),
				AttrName:     fmt.Sprintf("%v", line[currentLabels.AttrName]),
				AttrValue:    fmt.Sprintf("%v", line[currentLabels.AttrValue]),
				UoM:          fmt.Sprintf("%v", line[currentLabels.UoM]),
				Description:  fmt.Sprintf("%v", line[currentLabels.Description]),
				DataType:     fmt.Sprintf("%v", line[currentLabels.DataType]),
				IsMandatory:  fmt.Sprintf("%v", line[currentLabels.IsMandatory]),
				CodedVal:     fmt.Sprintf("%v", line[currentLabels.CodedVal]),
			}
			reportData = append(reportData, r)
		}
	}
	return reportData
}
