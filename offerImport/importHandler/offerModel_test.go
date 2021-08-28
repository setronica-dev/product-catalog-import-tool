package importHandler

import (
	"reflect"
	"testing"
)

func Test_getCountries(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "positive: list of countries from API should be processed as list",
			args: args{
				input: "[de es fr]",
			},
			want: []string{
				"de",
				"es",
				"fr",
			},
		},
		{
			name: "positive: empty list of countries should be processed as empty value",
			args: args{
				input: "[]",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCountriesAsArray(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCountriesAsArray() = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}
