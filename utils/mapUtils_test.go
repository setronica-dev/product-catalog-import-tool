package utils

import "testing"

func TestIsEmptyMap(t *testing.T) {
	type args struct {
		row map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive: row without values should be considered as empty",
			args: args{
				row: map[string]interface{}{
					"Offer":    "",
					"Receiver": nil,
				},
			},
			want: true,
		},
		{
			name: "positive: row with part of values should be considered as not empty",
			args: args{
				row: map[string]interface{}{
					"Offer":      "",
					"Receiver":   nil,
					"ContractID": "123",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyMap(tt.args.row); got != tt.want {
				t.Errorf("IsEmptyMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
