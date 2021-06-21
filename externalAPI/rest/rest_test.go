package rest

import (
	"net/http"
	"testing"
)

func TestAPIClient_buildUrl(t *testing.T) {
	type fields struct {
		BaseURL    string
		Auth       Auth
		HTTPClient *http.Client
	}
	type args struct {
		method string
		params []UrlParam
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "check that special symbols in URL parameter values are URL-encoded",
			fields: fields{
				BaseURL: "http://test.test",
				Auth:    Auth{},
			},
			args: args{
				method: "/method",
				params: []UrlParam{
					{
						Key:   "url&Key",
						Value: "quoted\"text",
					},
				},
			},
			want: "http://test.test/method?url%26Key=quoted%22text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &APIClient{
				BaseURL:    tt.fields.BaseURL,
				Auth:       tt.fields.Auth,
				HTTPClient: tt.fields.HTTPClient,
			}
			if got := c.buildUrl(tt.args.method, tt.args.params); got != tt.want {
				t.Errorf("buildUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
