package utils

import (
	"reflect"
	"testing"
)

func TestRowsToMapRows(t *testing.T) {
	type args struct {
		data   [][]string
		header []string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "positive: data from cells should be linked with columns",
			args: args{
				data: [][]string{
					{
						"title1", "title2",
					}, {
						"b1", "b2",
					}, {
						"c1", "c2",
					},
					{
						"d1",
					},
				},
				header: []string{
					"title1",
					"title2",
				},
			},
			want: []map[string]interface{}{
				{
					"title1": "b1",
					"title2": "b2",
				},
				{
					"title1": "c1",
					"title2": "c2",
				},
				{
					"title1": "d1",
				},
			},
			wantErr: false,
		},
		{
			name: "negative: data without defined header can not be parsed",
			args: args{
				data: [][]string{
					{
						"test1", "test2",
					},
				},
				header: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RowsToMapRows(tt.args.data, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("RowsToMapRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RowsToMapRows() got = %v, want %v", got, tt.want)
			}
		})
	}
}
