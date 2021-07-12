package adapters

import "testing"


func TestGetFileName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive: file name should be selected from path",
			args: args{
				path: "./data/source/input.xlsx::Products",
			},
			want: "input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileName(tt.args.path); got != tt.want {
				t.Errorf("GetFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
