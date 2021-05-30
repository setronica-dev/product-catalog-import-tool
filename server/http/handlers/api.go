package handlers

import (
	"go.uber.org/dig"
	"net/http"
	"ts/api"
)

type Deps struct {
	dig.In
}

func New(deps Deps) api.ServerInterface {
	return &ApiHandler{deps: deps}
}

type ApiHandler struct {
	deps Deps
}

func (a ApiHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	Handle(w, r, func(ctx *Context) error {
		return ctx.NoContent(http.StatusNoContent)
	})
}
