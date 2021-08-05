package reports

import "strings"

func (r *ReportsHandler) buildFailuresReportData(report []Report) [][]string {
	var res [][]string
	res = append(res, r.buildFailureReportHeader())

	for _, item := range report {
		recordItem := r.buildRaw(item)
		res = append(res, recordItem)
	}
	return res
}

func (r *ReportsHandler) buildFailureReportHeader() []string {
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
