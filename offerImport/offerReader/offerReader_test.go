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
