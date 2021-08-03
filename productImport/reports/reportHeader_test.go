package reports

import (
	"reflect"
	"testing"
	"ts/productImport/mapping"
)

func Test_buildHeader(t *testing.T) {
	type args struct {
		source      map[string]interface{}
		reportItems []Report
		columnMap   *mapping.ColumnMapConfig
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 map[string]int64
	}{
		// tests for manual run only
		/*	{
				name: "positive: success report header should be built in Tradeshift format(with default column values for productID and Category)",
				args: args{
					source: map[string]interface{}{
						"ProductID":  "1233",
						"UNSPSC":     "1321442",
						"PName":      "Test product",
						"Attribute1": "High availability",
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
					columnMap: &mapping.ColumnMapConfig{
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
				name: "positive: in success report column names in mapping and product header should be compatible regardless of *, spaces, tabs," +
					"header indexes should be oriented on source labels",
				args: args{
					source: map[string]interface{}{
						"ProductID*": "1233",
						"UNSPSC":     "1321442",
						"PName":      "Test product",
						"Attribute1": "High availability",
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
					columnMap: &mapping.ColumnMapConfig{
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
					"ProductID*":  0,
					"UNSPSC":   1,
					"PName":      2,
					"Attribute1": 3,
				},
			},
			{
				name: "positive: in success report columns in new header should be ordered: " +
					"first should be product ID, then-Category, then-all other fields",
				args: args{
					source: map[string]interface{}{
						"Attribute1": "AttrValue1",
						"Category":   "09876",
						"ID":         "12345",
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
					columnMap: &mapping.ColumnMapConfig{
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
			{
				name: "positive: attributes without category should be added to common report ",
				args: args{
					source: map[string]interface{}{
						"ProductID*": "1233",
						"UNSPSC":     "1321442",
						"PName":      "Test product",
						"Attribute1": "High availability",
					},
					reportItems: []Report{
						{
							ProductId:    "123",
							Name:         "Test product",
							Category:     "1321442",
							CategoryName: "Test Category Name",
							AttrName:     "Attribute2",
							AttrValue:    "Weight",
						},
					},
					columnMap: &mapping.ColumnMapConfig{
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
					"Attribute2",
				},
				want1: map[string]int64{
					"ProductID*":  0,
					"UNSPSC":   1,
					"PName":      2,
					"Attribute1": 3,
					"Attribute2": 4,
				},
			},*/
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

func Test_buildProductsPart(t *testing.T) {
	type args struct {
		sourceOrderedRow []string
		columnMapConfig  *mapping.ColumnMapConfig
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 map[string]int64
	}{
		{
			name: "Positive:",
			args: args{
				sourceOrderedRow: []string{
					"SKU*",
					"UNSPSC*",
					"ProductName",
					"Column1",
					"column 2",
					" Column 3",
				},
				columnMapConfig: &mapping.ColumnMapConfig{
					ProductID: "SKU",
					Category:  "UNSPSC",
					Name:      "Product Name",
					OtherColumns: []*mapping.ColumnItem{
						{
							DefaultKey: "Key1",
							MappedKey:  "Column 1",
						},
						{
							DefaultKey: "Key2",
							MappedKey:  "Column 2",
						},
					},
				},
			},
			want: []string{
				"ID", "Category", "Name", "Key1", "Key2", " Column 3",
			},
			want1: map[string]int64{
				"SKU*":        0,
				"UNSPSC*":     1,
				"ProductName": 2,
				"Column1":     3,
				"column 2":    4,
				" Column 3":   5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := buildProductsHeaderPart(tt.args.sourceOrderedRow, tt.args.columnMapConfig)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildProductsHeaderPart() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildProductsHeaderPart() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
