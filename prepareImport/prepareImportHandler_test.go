package prepareImport

import "testing"

func Test_isXLSX(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive: .xlsx file is considered as exel file",
			args: args{
				filePath: "test.XlSx",
			},
			want: true,
		},
		{
			name: "positive: .xls file is considered as exel file",
			args: args{
				filePath: "test.Xls",
			},
			want: true,
		},
		{
			name: "positive: not .xlsx or .xls file is considered as not exel file",
			args: args{
				filePath: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isXLSX(tt.args.filePath); got != tt.want {
				t.Errorf("isXLSX() = %v, want %v", got, tt.want)
			}
		})
	}
}
