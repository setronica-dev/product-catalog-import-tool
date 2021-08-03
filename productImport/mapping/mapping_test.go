package mapping

import (
	"reflect"
	"testing"
)

func Test_mapping_Parse(t *testing.T) {
	type fields struct {
		rawMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   *ColumnMapConfig
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
			want: &ColumnMapConfig{
				ProductID:    "Label1",
				Category:     "Label2",
				Name:         "Label3",
				OtherColumns: []*ColumnItem{},
			},
		},
		{
			name: "positive: empty map should be converted to MAp Object with default values of ProductID, Category and Name",
			fields: fields{
				rawMap: nil,
			},
			want: &ColumnMapConfig{
				ProductID:    "ID",
				Category:     "Category",
				Name:         "Name",
				OtherColumns: []*ColumnItem{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapping{
				rawMap: tt.fields.rawMap,
			}
			if got := m.NewColumnMap(m.rawMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewColumnMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
