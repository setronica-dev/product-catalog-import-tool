package utils

import "testing"

func TestTrimAll(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive: column names should be transformed to spaces-, tabs-, *- free and converted to low-case " +
				"format to be compatible with default column names",
			args: args{
				input: "This *is	column NaMe",
			},
			want: "thisiscolumnname",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimAll(tt.args.input); got != tt.want {
				t.Errorf("TrimAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
