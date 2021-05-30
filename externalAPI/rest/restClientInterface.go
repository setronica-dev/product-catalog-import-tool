package rest

import (
	"net/http"
)

type RestClientInterface interface {
	Post(method string, body map[string]interface{}, params []UrlAttributes) (*http.Response, error)
	PostFile(method string, filePath string) (*http.Response, error)
	Get(method string) (*http.Response, error)
}
