package offerReader

import (
	"reflect"
	"testing"
)

func Test_processOffer(t *testing.T) {
	type args struct {
		header *RawHeader
		item   map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *RawOffer
	}{
		{
			name: "positive: row without required fields should be skipped",
			args: args{
				header: &RawHeader{
					Offer:    "Offers",
					Receiver: "Countries",
				},
				item: map[string]interface{}{
					"Offers":    nil,
					"Receiver":  "",
					"Countries": "US",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processOffer(tt.args.header, tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processOffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEmptyRow(t *testing.T) {
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
			if got := isEmptyRow(tt.args.row); got != tt.want {
				t.Errorf("isEmptyRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
