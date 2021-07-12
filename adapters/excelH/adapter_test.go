package excelH

import "testing"

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
		}, /**/
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
