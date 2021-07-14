package prepareImport

import "testing"

func TestXLSXSheetToCSVConverter_buildPath(t *testing.T) {
	type fields struct {
		sheet                 string
		destinationPath       string
		destinationFileSuffix string
	}
	type args struct {
		sourceFilePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "positive: path to source file should be converted to csv file with new destination and suffix",
			fields: fields{
				sheet: "Sheet1",
				destinationPath: "./data/test/",
				destinationFileSuffix: "-new",
			},
			args: args{
				sourceFilePath: "./input/hello.txt",
			},
			want: "data/test/hello-new.csv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &XLSXSheetToCSVConverter{
				sheet:                 tt.fields.sheet,
				destinationPath:       tt.fields.destinationPath,
				destinationFileSuffix: tt.fields.destinationFileSuffix,
			}
			if got := c.buildPath(tt.args.sourceFilePath); got != tt.want {
				t.Errorf("buildPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
