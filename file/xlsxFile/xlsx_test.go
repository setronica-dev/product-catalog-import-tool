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

func Test_parsePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "positive: path should be split to filePath and sheetName by first \"::\"delimiter",
			args: args{
				path: "test::test1::test2",
			},
			want:    "test",
			want1:   "test1::test2",
			wantErr: false,
		},
		{
			name: "negative: path without filePath is invalid",
			args: args{
				path: "::test1",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "negative: path without sheetName is invalid",
			args: args{
				path: "test1::",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "negative: path only with filePath is invalid",
			args: args{
				path: "test1",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parsePath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parsePath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
