package csvH

import "testing"

func Test_isValidRow(t *testing.T) {
	type args struct {
		row []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive: row with not empty and not formatted-keys values is considered as valid",
			args: args{
				row: []string{
					"sku*",
					"unspsc",
				},
			},
			want: true,
		},
		{
			name: "positive: row with header formatting should be considered as invalid",
			args: args{
				row: []string{
					"",
					"HEADER_V3_START",
				},
			},
			want: false,
		},
		{
			name: "positive: row with empty strings should be considered as empty",
			args: args{
				row: []string{
					"",
					"",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidRow(tt.args.row); got != tt.want {
				t.Errorf("isValidRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
