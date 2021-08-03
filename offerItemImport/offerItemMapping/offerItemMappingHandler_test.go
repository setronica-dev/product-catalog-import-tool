package offerItemMapping

import (
	"reflect"
	"testing"
	"ts/productImport/mapping"
)

func TestOfferItemReader_buildHeader(t *testing.T) {
	type fields struct {
		columnMap         *mapping.ColumnMapConfig
		sourcePath        string
		successReportPath string
	}
	type args struct {
		row []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "positive",
			fields: fields{
				columnMap: &mapping.ColumnMapConfig{
					ProductID: "Product ID",
				},
			},
			args: args{
				row: []string{
					"Name",
					"Offer",
					"productID *",
					"Price",
				},
			},
			want: []string{
				"Name",
				"Offer",
				"ID",
				"Price",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oi := &OfferItemMappingHandler{
				columnMap:         tt.fields.columnMap,
				sourcePath:        tt.fields.sourcePath,
				successReportPath: tt.fields.successReportPath,
			}
			if got := oi.buildHeader(tt.args.row); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
