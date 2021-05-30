package handlers

import (
	"encoding/json"
	"net/http"
	"ts/errors"
)

const (
	charsetUTF8                    = "charset=UTF-8"
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	HeaderContentType              = "Content-Type"
)

type Context struct {
	writer  http.ResponseWriter
	request *http.Request
}

func (c *Context) Bind(i interface{}) error {
	return json.NewDecoder(c.request.Body).Decode(i)
}

func (c *Context) NoContent(code int) error {
	c.writer.WriteHeader(code)
	return nil
}

func (c *Context) JSON(code int, i interface{}) error {
	c.writer.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.writer.WriteHeader(code)
	return json.NewEncoder(c.writer).Encode(i)
}

func (c *Context) Error(code int, err error) error {
	c.writer.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.writer.WriteHeader(code)
	return json.NewEncoder(c.writer).Encode(err)
}

type HandlerFunc func(ctx *Context) error

func Handle(w http.ResponseWriter, r *http.Request, handlerFunc HandlerFunc) {
	ctx := &Context{
		writer:  w,
		request: r,
	}

	if err := handlerFunc(ctx); err != nil {
		httpErrorHandler(w, r, err)
	}
}

func httpErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var he errors.HttpError
	he, ok := err.(errors.HttpError)
	if !ok {
		he = errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	r.Response.Header.Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	r.Response.StatusCode = he.Code
	_ = json.NewEncoder(w).Encode(he)
}
