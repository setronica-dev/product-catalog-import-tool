package importHandler

import "testing"

func Test_isId(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "id consists of 5 hex-formatted parts connected with '-'",
			args: args{
				input: "2addd7f5-5633-4c74-b3d7-760d04f1e4cc",
			},
			want: true,
		},
		{
			name: "nonHexFormated string is not considered as id",
			args: args{
				input: "2add7f5-5633-4c74-bRd7-760d04f1e4cc",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isId(tt.args.input); got != tt.want {
				t.Errorf("isId() = %v, want %v", got, tt.want)
			}
		})
	}
}
