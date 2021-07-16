package mapping

import (
	"reflect"
	"testing"
)

func Test_mapping_Parse(t *testing.T) {
	type fields struct {
		rawMap    map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   *ColumnMap
	}{
		{
			name: "positive: map should be converted to object with ProductID, Category and Name from relative columns",
			fields: fields{
				rawMap: map[string]string{
					"ID":       "Label1",
					"Category": "Label2",
					"Name":     "Label3",
				},
			},
			want: &ColumnMap{
				ProductID: "Label1",
				Category:  "Label2",
				Name:      "Label3",
			},
		},
		{
			name: "positive: empty map should be converted to MAp Object with default values of ProductID, Category and Name",
			fields: fields{
				rawMap: nil,
			},
			want: &ColumnMap{
				ProductID: "ID",
				Category:  "Category",
				Name:      "Name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapping{
				rawMap:    tt.fields.rawMap,
			}
			if got := m.Parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
