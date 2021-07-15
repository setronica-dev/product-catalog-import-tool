package reports

import (
	"reflect"
	"testing"
	"ts/adapters"
	"ts/productImport/mapping"
)

func Test_initFirstRaw(t *testing.T) {
	type args struct {
		m *mapping.ColumnMap
	}
	tests := []struct {
		name string
		args args
		want *ReportLabels
	}{
		{
			name: "positive: header labels for failures report should be taken from mapping",
			args: args{
				m: &mapping.ColumnMap{
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
				m: &mapping.ColumnMap{},
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
			if got := initFirstRaw(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initFirstRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReportsHandler_buildSuccessMapRaw(t *testing.T) {
	type fields struct {
		Handler     adapters.HandlerInterface
		Header      *ReportLabels
		ColumnMap   *ColumnMap
		FileManager *adapters.FileManager
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
		{
			name: "positive: attribute value in report should be replaced with value from fixed data. " +
				"Report should contain header in TradeShift format" +
				"Source data in report (CategoryID, Name) should be actualized relatively to fixed data",

			fields: fields{
				ColumnMap: &ColumnMap{
					ProductID: "ProductID",
					Category:  "UNSPSC",
					Name:      "Product Name",
				},
			},
			args: args{
				source: []map[string]interface{}{
					{
						"ProductID":    "1",
						"UNSPSC":       "1000011",
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
		/* row order is changeble- test are unstable for CI run
		{
			name: "positive: attributes from source data and fixed data should be added to report",

			fields: fields{
				ColumnMap: &ColumnMap{
					ProductID: "ProductID",
					Category:  "UNSPSC",
					Name:      "Product Name",
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
				ColumnMap: &ColumnMap{
					ProductID: "ProductID",
					Category:  "UNSPSC",
					Name:      "Product Name",
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
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReportsHandler{
				Handler:     tt.fields.Handler,
				Header:      tt.fields.Header,
				ColumnMap:   tt.fields.ColumnMap,
				FileManager: tt.fields.FileManager,
			}
			if got := r.buildSuccessMapRaw(tt.args.source, tt.args.reportItems); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSuccessMapRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildHeader(t *testing.T) {
	type args struct {
		source      []map[string]interface{}
		reportItems []Report
		columnMap   *ColumnMap
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 map[string]int64
	}{
		{
			name: "positive: success report header should be built in Tradeshift format(with default column values for productID and Category)",
			args: args{
				source: []map[string]interface{}{
					{
						"ProductID":  "1233",
						"UNSPSC":     "1321442",
						"PName":      "Test product",
						"Attribute1": "High availability",
					},
				},
				reportItems: []Report{
					{
						ProductId:    "123",
						Name:         "Test product",
						Category:     "1321442",
						CategoryName: "Test Category Name",
						AttrName:     "Attribute1",
						AttrValue:    "High availability",
					},
				},
				columnMap: &ColumnMap{
					Category:  "UNSPSC",
					ProductID: "ProductID",
					Name:      "PName",
				},
			},
			want: []string{
				"ID",
				"Category",
				"Name",
				"Attribute1",
			},
			want1: map[string]int64{
				"ProductID":  0,
				"UNSPSC":     1,
				"PName":      2,
				"Attribute1": 3,
			},
		},
		{
			name: "positive: in success report column names in mapping and product header should be compatible regardless of *, spaces, tabs",
			args: args{
				source: []map[string]interface{}{
					{
						"ProductID*": "1233",
						"UNSPSC":     "1321442",
						"PName":      "Test product",
						"Attribute1": "High availability",
					},
				},
				reportItems: []Report{
					{
						ProductId:    "123",
						Name:         "Test product",
						Category:     "1321442",
						CategoryName: "Test Category Name",
						AttrName:     "Attribute1",
						AttrValue:    "High availability",
					},
				},
				columnMap: &ColumnMap{
					Category:  "UNSPSC *",
					ProductID: "ProductID",
					Name:      "PName",
				},
			},
			want: []string{
				"ID",
				"Category",
				"Name",
				"Attribute1",
			},
			want1: map[string]int64{
				"ProductID":  0,
				"UNSPSC *":   1,
				"PName":      2,
				"Attribute1": 3,
			},
		},
		{
			name: "positive: in success report columns in new header should be ordered: " +
				"first should be product ID, then-Category, then-all other fields",
			args: args{
				source: []map[string]interface{}{
					{
						"Attribute1": "AttrValue1",
						"Category":   "09876",
						"ID":         "12345",
					},
				},
				reportItems: []Report{
					{
						ProductId:    "123",
						Name:         "Test product",
						Category:     "1321442",
						CategoryName: "Test Category Name",
						AttrName:     "Attribute1",
						AttrValue:    "High availability",
					},
				},
				columnMap: &ColumnMap{
					Category:  "Category",
					ProductID: "ID",
					Name:      "Name",
				},
			},
			want: []string{
				"ID",
				"Category",
				"Attribute1",
			},
			want1: map[string]int64{
				"ID":         0,
				"Category":   1,
				"Attribute1": 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := buildSuccessReportHeader(tt.args.source, tt.args.reportItems, tt.args.columnMap)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSuccessReportHeader() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildSuccessReportHeader() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getSourceKeys(t *testing.T) {
	type args struct {
		sourceRow map[string]interface{}
		columnMap *ColumnMap
	}
	tests := []struct {
		name string
		args args
		want *ColumnMap
	}{
		{
			name: "positive: should be selected compatible with mapped column names from source product data",
			args: args{
				sourceRow: map[string]interface{}{
					"UNSPSC":        "1111",
					"ProductID*":    "22222",
					"Name ":         "33333",
					"Category Name": "44444",
				},
				columnMap: &ColumnMap{
					Category:  "Unspsc",
					ProductID: "Product ID",
					Name:      "NAME*",
				},
			},
			want: &ColumnMap{
				Category:  "UNSPSC",
				ProductID: "ProductID*",
				Name:      "Name ",
			},
		},
		{
			name: "positive: should not be selected column names from source products data which are incompatible " +
				"with mapping column names",
			args: args{
				sourceRow: map[string]interface{}{
					"Category":   "1111",
					"ProductID*": "22222",
					"Name":       "33333",
				},
				columnMap: &ColumnMap{
					Category:  "Unspsc",
					ProductID: "ProductID",
					Name:      "Name",
				},
			},
			want: &ColumnMap{
				ProductID: "ProductID*",
				Name:      "Name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSourceKeys(tt.args.sourceRow, tt.args.columnMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSourceKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
