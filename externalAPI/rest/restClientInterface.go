package rest

import (
	"io"
	"net/http"
)

type RestClientInterface interface {
	Post(method string, body io.Reader, params []UrlParam) (*http.Response, error)
	PostFile(method string, filePath string) (*http.Response, error)
	Get(method string, params []UrlParam) (*http.Response, error)
}
