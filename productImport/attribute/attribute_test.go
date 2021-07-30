package attribute

import (
	"reflect"
	"testing"
	"ts/adapters"
)

func TestAttributeHandler_getCurrentHeader(t *testing.T) {
	type fields struct {
		handler     adapters.HandlerInterface
		columnMap   *ColumnMap
	}
	type args struct {
		row map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ColumnMap
	}{
		{
			name: "positive: should be selected compatible with mapped column names from attributes data",
			fields: fields{
				columnMap: &ColumnMap{
					ProductId:    "SKU*",
					Name:         "Product NAme",
					Category:     "UNSPSC* ",
					CategoryName: "UNSPSC NAME",
					AttrName:     "Attribute Name*",
					AttrValue:    "Attribute Value",
					UoM:          "uom",
					DataType:     "Data type",
					Description:  " Description",
					IsMandatory:  "is mandatory",
					CodedVal:     "coded value",
				},
			},
			args: args{
				row: map[string]interface{}{
					"SKU":              "11",
					"PRODUCT NAME":     "22",
					"UNSPSC":           "33",
					"UNSPSCName":       "44",
					"AttributeName *":  "55",
					"Attribute Value*": "66",
					"UoM*":             "77",
					"Data Type":        "88",
					"Description":      "99",
					"Is Mandatory":     "00",
					"CodedVal":         "aa",
				},
			},
			want: &ColumnMap{
				ProductId:    "SKU",
				Name:         "PRODUCT NAME",
				Category:     "UNSPSC",
				CategoryName: "UNSPSCName",
				AttrName:     "AttributeName *",
				AttrValue:    "Attribute Value*",
				UoM:          "UoM*",
				DataType:     "Data Type",
				Description:  "Description",
				IsMandatory:  "Is Mandatory",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ah := &AttributeHandler{
				handler:     tt.fields.handler,
				columnMap:   tt.fields.columnMap,
			}
			if got := ah.getCurrentHeader(tt.args.row); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCurrentHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeHandler_parseData(t *testing.T) {
	type fields struct {
		fileManager *adapters.FileManager
		handler     adapters.HandlerInterface
		columnMap   *ColumnMap
	}
	type args struct {
		reportDataSource []map[string]interface{}
		currentLabels    *ColumnMap
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Attribute
	}{
		{
			name: "positive: Raw attribute data should be parsed to structured according to mapped column names," +
				"compatible with names, defined in product mapping",
			fields: fields{
				columnMap: &ColumnMap{
					ProductId:    "SKU*",
					Name:         "Product NAme",
					Category:     "UNSPSC* ",
					CategoryName: "UNSPSC NAME",
					AttrName:     "Attribute Name*",
					AttrValue:    "Attribute Value",
					UoM:          "uom",
					DataType:     "Data type",
					Description:  " Description",
					IsMandatory:  "is mandatory",
					CodedVal:     "coded val",
				},
			},
			args: args{
				reportDataSource: []map[string]interface{}{
					{
						"SKU":              "11",
						"PRODUCT NAME":     "22",
						"UNSPSC":           "33",
						"UNSPSCName":       "44",
						"AttributeName *":  "55",
						"Attribute Value*": "66",
						"UoM*":             "77",
						"Data Type":        "88",
						"Description":      "99",
						"Is Mandatory":     "00",
						"CodedVal":         "aa"},
				},
				currentLabels: &ColumnMap{
					ProductId:    "SKU",
					Name:         "PRODUCT NAME",
					Category:     "UNSPSC",
					CategoryName: "UNSPSCName",
					AttrName:     "AttributeName *",
					AttrValue:    "Attribute Value*",
					UoM:          "UoM*",
					DataType:     "Data Type",
					Description:  "Description",
					IsMandatory:  "Is Mandatory",
					CodedVal:     "CodedVal",
				},
			},
			want: []*Attribute{
				{
					ProductId:    "11",
					Name:         "22",
					Category:     "33",
					CategoryName: "44",
					AttrName:     "55",
					AttrValue:    "66",
					UoM:          "77",
					DataType:     "88",
					Description:  "99",
					IsMandatory:  "00",
					CodedVal:     "aa",
				},
			},
		},
		{
			name: "positive: empty row should be skipped",
			fields: fields{
				columnMap: &ColumnMap{
					ProductId:    "SKU*",
					Name:         "Product NAme",
					Category:     "UNSPSC* ",
					CategoryName: "UNSPSC NAME",
					AttrName:     "Attribute Name*",
					AttrValue:    "Attribute Value",
					UoM:          "uom",
					DataType:     "Data type",
					Description:  " Description",
					IsMandatory:  "is mandatory",
					CodedVal:     "coded val",
				},
			},
			args: args{
				reportDataSource: []map[string]interface{}{
					{
						"SKU":              "",
						"PRODUCT NAME":     "",
						"UNSPSC":           "",
						"UNSPSCName":       "",
						"AttributeName *":  "",
						"Attribute Value*": "",
						"UoM*":             "",
						"Data Type":        "",
						"Description":      "",
						"Is Mandatory":     "",
						"CodedVal":         "",
					},
				},
				currentLabels: &ColumnMap{
					ProductId:    "SKU",
					Name:         "PRODUCT NAME",
					Category:     "UNSPSC",
					CategoryName: "UNSPSCName",
					AttrName:     "AttributeName *",
					AttrValue:    "Attribute Value*",
					UoM:          "UoM*",
					DataType:     "Data Type",
					Description:  "Description",
					IsMandatory:  "Is Mandatory",
				},
			},
			want: []*Attribute{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ah := &AttributeHandler{
				handler:     tt.fields.handler,
				columnMap:   tt.fields.columnMap,
			}
			if got := ah.parseData(tt.args.reportDataSource, tt.args.currentLabels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
