package reports

import (
	"fmt"
	"path/filepath"
	"strings"
	"ts/adapters"
)

func (r *ReportsHandler) writeFailedReport(report []Report, feedFilePath string) string {
	filePath := filepath.Join(r.FailResultPath, buildAttributesFileName(feedFilePath))
	var data [][]string
	data = r.buildAttributesReportData(report)
	r.Handler.Write(filePath, data)
	return filePath

}

func buildAttributesFileName(feedPath string) string {
	ext := filepath.Ext(feedPath)
	sourceFileName := adapters.GetFileName(feedPath)
	return fmt.Sprintf("%v_attributes%v", sourceFileName, ext)
}

func (r *ReportsHandler) buildAttributesReportData(report []Report) [][]string {
	var res [][]string
	res = append(res, r.buildAttributesReportHeader())

	for _, item := range report {
		recordItem := r.buildAttributesRaw(item)
		res = append(res, recordItem)
	}
	return res
}

func (r *ReportsHandler) buildAttributesReportHeader() []string {
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

func (r *ReportsHandler) buildAttributesRaw(item Report) []string {
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
