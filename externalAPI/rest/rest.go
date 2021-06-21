package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	oauth1 "github.com/klaidas/go-oauth1"
	"go.uber.org/dig"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"ts/config"
)

type APIClient struct {
	BaseURL    string
	Auth       Auth
	HTTPClient *http.Client
}

type Auth struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

type UrlParam struct {
	Key   string
	Value string
}

type Deps struct {
	dig.In
	Config *config.Config
}

func NewRestClient(deps Deps) RestClientInterface {
	tsConfig := deps.Config.TradeshiftAPI
	c := APIClient{
		BaseURL:    tsConfig.APIBaseURL,
		HTTPClient: &http.Client{},
		Auth: Auth{
			ConsumerKey:    tsConfig.ConsumerKey,
			ConsumerSecret: tsConfig.ConsumerSecret,
			Token:          tsConfig.Token,
			TokenSecret:    tsConfig.TokenSecret,
		},
	}
	return &c
}

func (c *APIClient) Post(method string, body io.Reader, params []UrlParam) (*http.Response, error) {
	contentType := "application/json"
	req := c.buildRequest(
		http.MethodPost,
		c.buildUrl(method, params),
		body,
		buildParams(params))
	req.Header.Set("Content-Type", contentType)

	resp, err := c.executeRequest(req)

	return resp, err
}

func (c *APIClient) PostFile(method string, filePath string) (*http.Response, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, f)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req := c.buildRequest(http.MethodPost, c.buildUrl(method, nil), body, nil)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := c.executeRequest(req)
	return res, err
}

func (c *APIClient) Get(method string, params []UrlParam) (*http.Response, error) {
	req := c.buildRequest(
		http.MethodGet,
		c.buildUrl(method, params),
		nil,
		buildParams(params))
	resp, err := c.executeRequest(req)
	return resp, err
}

func (c *APIClient) buildRequest(method string, path string, body io.Reader, params map[string]string) *http.Request {
	auth := oauth1.OAuth1{
		ConsumerKey:    c.Auth.ConsumerKey,
		ConsumerSecret: c.Auth.ConsumerSecret,
		AccessToken:    c.Auth.Token,
		AccessSecret:   c.Auth.TokenSecret,
	}
	authHeader := auth.BuildOAuth1Header(method, path, params)
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Authorization", authHeader)
	return req
}

func BuildBody(data interface{}) io.Reader {
	body, _ := json.Marshal(&data)
	return bytes.NewReader(body)
}

func (c *APIClient) buildUrl(method string, params []UrlParam) string {
	endpoint := fmt.Sprintf("%v%v", c.BaseURL, method)
	u, _ := url.Parse(endpoint)
	q := u.Query()
	for _, v := range params {
		q.Add(v.Key, v.Value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func buildParams(params []UrlParam) map[string]string {
	res := map[string]string{}

	for _, item := range params {
		res[item.Key] = item.Value
	}
	return res
}

func (c *APIClient) executeRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		r, _ := ParseResponseToString(resp)
		return nil, fmt.Errorf("status code is %v: %v", resp.StatusCode, r)
	}
	return resp, nil
}

func ParseResponse(r *http.Response) (map[string]interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("http response is empty")
	}
	var resp map[string]interface{}

	if r.Body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("failed to read response body %v", err)
		return nil, err
	}
	defer r.Body.Close()

	if !json.Valid(body) {
		return nil, fmt.Errorf("caught invalid json")
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body to object %v", err)
	}

	return resp, nil
}

func ParseResponseToString(r *http.Response) (string, error) {
	if r == nil {
		return "", fmt.Errorf("http response is empty")
	}
	if r.Body == nil {
		return "", fmt.Errorf("empty response body")
	}

	body, err := ioutil.ReadAll(r.Body)
	return string(body), err
}
