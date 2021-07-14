package reports

import (
	"reflect"
	"testing"
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
			name: "positive: header labels should be taken from mapping",
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
