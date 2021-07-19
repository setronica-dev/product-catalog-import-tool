package product

import (
	"reflect"
	"testing"
	"ts/adapters"
)

func Test_parse(t *testing.T) {
	type args struct {
		sourceData []map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []map[string]interface{}
	}{
		{
			name: "positive: empty rows in uploaded data should be skipped",
			args: args{
				sourceData: []map[string]interface{}{
					{
						"ke1":  "",
						"key2": "",
						"key3": nil,
					},
					{
						"ke1":  "1",
						"key2": "2",
						"key3": "3",
					},
				},
			},
			want: []map[string]interface{}{
				{
					"ke1":  "1",
					"key2": "2",
					"key3": "3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parse(tt.args.sourceData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductHandler_GetCurrentHeader(t *testing.T) {
	type fields struct {
		fileManager *adapters.FileManager
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
			name: "positive: should be selected compatible with mapped column names from source product data",
			fields: fields{
				columnMap: &ColumnMap{
					Category:  "Unspsc",
					ProductID: "Product ID",
					Name:      "NAME*",
				},
			},
			args: args{
				row: map[string]interface{}{
					"UNSPSC":        "1111",
					"ProductID*":    "22222",
					"Name ":         "33333",
					"Category Name": "44444",
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
			fields: fields{
				columnMap: &ColumnMap{
					Category:  "Unspsc",
					ProductID: "ProductID",
					Name:      "Name",
				},
			},
			args: args{
				row: map[string]interface{}{
					"Category":   "1111",
					"ProductID*": "22222",
					"Name":       "33333",
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
			p := &ProductHandler{
				fileManager: tt.fields.fileManager,
				handler:     tt.fields.handler,
				columnMap:   tt.fields.columnMap,
			}
			if got := p.GetCurrentHeader(tt.args.row); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCurrentHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
