package tradeshiftAPI

import (
	"testing"
)

func Test_buildAdvancedSearchValue(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should be build json-formated value for advancedSearch parameter in search offers API request",
			args: args{
				name: "\"test\"",
			},
			want: "{\"name\":\"\\\"test\\\"\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildAdvancedSearchValue(tt.args.name); got != tt.want {
				t.Errorf("buildAdvancedSearchValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
