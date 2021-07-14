package xlsxFile

import (
	"github.com/tealeg/xlsx"
	"reflect"
	"testing"
)

func Test_getSheetData(t *testing.T) {
	type args struct {
		sheets []*xlsx.Sheet
		name   string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "positive: should be selected sheet with defined name, case-independent",
			args: args{
				sheets: []*xlsx.Sheet{
					{
						Name: "test1",
						Rows: []*xlsx.Row{
							{
								Cells: []*xlsx.Cell{
									{
										Value: "value1",
									},
								},
							},
						},
					},
					{
						Name: "Test2",
						Rows: []*xlsx.Row{
							{
								Cells: []*xlsx.Cell{
									{
										Value: "value2",
									},
								},
							},
						},
					},
				},
				name: "test2",
			},
			want: [][]string{
				{
					"value2",
				},
			},
		},
		{
			name: "positive:empty file should be processed without failures",
			args: args{
				sheets: []*xlsx.Sheet{
					{
						Name: "test1",
						Rows: []*xlsx.Row{
							{
								Cells: []*xlsx.Cell{{}},
							},
						},
					},
				},
				name: "test2",
			},
			want: [][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSheetData(tt.args.sheets, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSheetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processSheetData(t *testing.T) {
	type args struct {
		sheet *xlsx.Sheet
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "positive: xlsx file content should be stored as list of rows",
			args: args{
				sheet: &xlsx.Sheet{
					Rows: []*xlsx.Row{
						{
							Cells: []*xlsx.Cell{
								{
									Value: "a1",
								},
								{
									Value: "b1",
								},
								{
									Value: "c1",
								},
							},
						},
						{
							Cells: []*xlsx.Cell{
								{
									Value: "a2",
								},
								{
									Value: "b2",
								},
							},
						},
					},
				},
			},
			want: [][]string{
				{
					"a1", "b1", "c1",
				},
				{
					"a2", "b2",
				},
			},
		},
		{
			name: "positive: file without data should be processed without failures",
			args: args{
				sheet: &xlsx.Sheet{
					Rows: []*xlsx.Row{
						{
							Cells: []*xlsx.Cell{},
						},
					},
				},
			},
			want: [][]string{{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processSheetData(tt.args.sheet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processSheetData() = %v, want %v", got, tt.want)
			}
		})
	}
}
