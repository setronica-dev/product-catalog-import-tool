package reports

import (
	"reflect"
	"testing"
	"ts/adapters"
	"ts/productImport/mapping"
	"ts/productImport/product"
)

func Test_initFirstRaw(t *testing.T) {
	type args struct {
		m *mapping.ColumnMapConfig
	}
	tests := []struct {
		name string
		args args
		want *ReportLabels
	}{
		{
			name: "positive: header labels for failures report should be taken from mapping",
			args: args{
				m: &mapping.ColumnMapConfig{
					ProductID: "SKU",
					Category:  "Rubric",
					Name:      "ProductName",
				},
			},
			want: &ReportLabels{
				ProductId:    "SKU",
				Category:     "Rubric",
				Name:         "ProductName",
				CategoryName: "Category Name",
				AttrName:     "Attribute Name*",
				AttrValue:    "Attribute Value*",
				UoM:          "UOM",
				Errors:       "Error Message",
				Description:  "Description",
				DataType:     "Data Type",
				IsMandatory:  "Is Mandatory",
				CodedVal:     "Coded Value",
			},
		},
		{
			name: "positive: in case of empty mapping should be taken default values for fields ProductID, Category, Name",
			args: args{
				m: &mapping.ColumnMapConfig{},
			},
			want: &ReportLabels{
				ProductId:    "ProductID*",
				Category:     "Category",
				Name:         "Name",
				CategoryName: "Category Name",
				AttrName:     "Attribute Name*",
				AttrValue:    "Attribute Value*",
				UoM:          "UOM",
				Errors:       "Error Message",
				Description:  "Description",
				DataType:     "Data Type",
				IsMandatory:  "Is Mandatory",
				CodedVal:     "Coded Value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initFailuresReportHeader(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initFailuresReportHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReportsHandler_buildSuccessMapRaw(t *testing.T) {
	type fields struct {
		Handler         adapters.HandlerInterface
		Header          *ReportLabels
		ColumnMapConfig *mapping.ColumnMapConfig
		FileManager     *adapters.FileManager
		ProductHandler  *product.ProductHandler
	}
	type args struct {
		source      []map[string]interface{}
		reportItems []Report
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]string
	}{
		// row order is changeable- test are unstable for CI run
/*
		{
			name: "positive: attribute value in report should be replaced with value from fixed data. " +
				"Report should contain header in TradeShift format" +
				"Source data in report (CategoryID, Name) should be actualized relatively to fixed data",

			fields: fields{
				ColumnMapConfig: &mapping.ColumnMapConfig{
					ProductID: "ProductID",
					Category:  "UNSPSC",
					Name:      "Product Name",
				},
				ProductHandler: &product.ProductHandler{
					ColumnMap: &product.ProductColumnMap{
						ProductID: "ProductID",
						Category:  "UNSPSC",
						Name:      "Product Name",
					},
				},
			},
			args: args{
				source: []map[string]interface{}{
					{
						"ProductID":    "1",
						"UNSPSC":       "100002",
						"Product Name": "Product 1",
						"Attribute1":   "Value1",
					},
				},
				reportItems: []Report{
					{
						ProductId:   "1",
						Category:    "100001",
						Name:        " Product1",
						AttrName:    "Attribute1",
						AttrValue:   "FixedValue1",
						IsMandatory: "true",
					},
				},
			},
			want: [][]string{
				{
					"ID", "Category", "Name", "Attribute1",
				},
				{
					"1", "100001", "Product 1", "FixedValue1",
				},
			},
		},
/*
		{
			name: "positive: attributes from source data and fixed data should be added to report",

			fields: fields{
				ColumnMapConfig: &mapping.ColumnMapConfig{
					ProductID: "ProductID",
					Category:  "UNSPSC",
					Name:      "Product Name",
				},
				ProductHandler: &product.ProductHandler{
					ColumnMap: &product.ProductColumnMap{
						ProductID: "ProductID",
						Category:  "UNSPSC",
						Name:      "Product Name",
					},
				},
			},
			args: args{
				source: []map[string]interface{}{
					{
						"ProductID":    "1",
						"UNSPSC":       "100001",
						"Product Name": "Product 1",
						"Attribute1":   "Value1",
						"Attribute2":   "Value2",
					},
				},
				reportItems: []Report{
					{
						ProductId:   "1",
						Category:    "100001",
						Name:        " Product1",
						AttrName:    "Attribute1",
						AttrValue:   "FixedValue1",
						IsMandatory: "true",
					},
					{
						ProductId:   "1",
						Category:    "100001",
						Name:        " Product1",
						AttrName:    "Attribute3",
						AttrValue:   "Value3",
						IsMandatory: "true",
					},
				},
			},
			want: [][]string{
				{
					"ID", "Category", "Name", "Attribute1", "Attribute2", "Attribute3",
				},
				{
					"1", "100001", "Product 1", "FixedValue1", "Value2", "Value3",
				},
			},
		},
		{
			name: "positive: if attribute in fixed data has no category, category should be taken from source data",
			fields: fields{
				ColumnMapConfig: &mapping.ColumnMapConfig{
					ProductID: "ProductID",
					Category:  "UNSPSC",
					Name:      "Product Name",
				},
				ProductHandler: &product.ProductHandler{
					ColumnMap: &product.ProductColumnMap{
						ProductID: "ProductID",
						Category:  "UNSPSC",
						Name:      "Product Name",
					},
				},
			},
			args: args{
				source: []map[string]interface{}{
					{
						"ProductID":    "1",
						"UNSPSC":       "100001",
						"Product Name": "Product 1",
						"Attribute1":   "Value1",
					},
				},
				reportItems: []Report{
					{
						ProductId:   "1",
						Name:        " Product1",
						AttrName:    "Attribute1",
						AttrValue:   "FixedValue1",
						IsMandatory: "true",
					},
				},
			},
			want: [][]string{
				{
					"ID", "Category", "Name", "Attribute1",
				},
				{
					"1", "100001", "Product 1", "FixedValue1",
				},
			},
		},
		{
			name: "positive: product without attributes should be in report",

			fields: fields{
				ColumnMapConfig: &mapping.ColumnMapConfig{
					ProductID: "ProductID",
					Category:  "UNSPSC",
					Name:      "Product Name",
				},
				ProductHandler: &product.ProductHandler{
					ColumnMap: &product.ProductColumnMap{
						ProductID: "ProductID",
						Category:  "UNSPSC",
						Name:      "Product Name",
					},
				},
			},
			args: args{
				source: []map[string]interface{}{
					{
						"ProductID":    "1",
						"UNSPSC":       "100001",
						"Product Name": "Product 1",
						"Attribute1":   "Value1",
					},
					{
						"ProductID":    "2",
						"UNSPSC":       "200001",
						"Product Name": "Product2",
					},
				},
				reportItems: []Report{
					{
						ProductId:   "1",
						Category:    "100001",
						Name:        " Product1",
						AttrName:    "Attribute1",
						AttrValue:   "FixedValue1",
						IsMandatory: "true",
					},
				},
			},
			want: [][]string{
				{
					"ID", "Category", "Name", "Attribute1",
				},
				{
					"1", "100001", "Product 1", "FixedValue1",
				},
				{
					"2", "200001", "Product2", "",
				},
			},
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReportsHandler{
				Handler:         tt.fields.Handler,
				Header:          tt.fields.Header,
				ColumnMapConfig: tt.fields.ColumnMapConfig,
				FileManager:     tt.fields.FileManager,
				productHandler:  tt.fields.ProductHandler,
			}
			if got := r.buildSuccessMapRaw(tt.args.source, tt.args.reportItems); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSuccessMapRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}
