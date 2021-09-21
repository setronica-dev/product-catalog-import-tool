package file

import (
	"reflect"
	"testing"
)

func Test_clearEmptyData(t *testing.T) {
	type args struct {
		data             [][]string
		headerLinesCount int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "positive: header rows should be skipped",
			args: args{
				data: [][]string{
					{"1", "2", "3"},
					{"a", "b", "c"},
					{"A", "B", "C"},
					{"AA", "BB", "CC"},
				},
				headerLinesCount: 1,
			},
			want: [][]string{
				{"a", "b", "c"},
				{"A", "B", "C"},
				{"AA", "BB", "CC"},
			},
		},

		{
			name: "positive: empty rows with non-empty header should be skipped",
			args: args{
				data: [][]string{
					{"", "", ""},
					{"1", "2", "3"},
					{"", "", ""},
					{"A", "B", "C"},
					{"AA", "BB", "CC"},
					{"", "", ""},
					{"", "", ""},
				},
				headerLinesCount: 1,
			},
			want: [][]string{
				{"1", "2", "3"},
				{"A", "B", "C"},
				{"AA", "BB", "CC"},
			},
		},
		{
			name: "positive: empty columns with non-empty header should be skipped",
			args: args{
				data: [][]string{
					{"01", "02", "03", "", "", "", ""},
					{"1", "2", "", "3", "", "", ""},
					{"A", "B", "", "C", "", "", ""},
					{"AA", "BB", "", "CC", "", ""},
				},
				headerLinesCount: 1,
			},
			want: [][]string{
				{"1", "2", "3"},
				{"A", "B", "C"},
				{"AA", "BB", "CC"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clearEmptyData(tt.args.data, tt.args.headerLinesCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clearEmptyData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValidColumnIndexes(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "positive: we should now indexes of all non-empty cells",
			args: args{
				data: []string{
					"0", "", "", "1", "", "", "",
				},
			},
			want: []int{
				0, 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValidColumnIndexes(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValidColumnIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}
