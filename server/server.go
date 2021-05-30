package server

import (
	"encoding/json"
	rice "github.com/GeertJohan/go.rice"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/dig"
	"log"
	"net/http"
	"ts/api"
	"ts/config"
)

const (
	BaseApiPath       = "/api"
	BaseSwaggerUIPath = "/swagger"
)

type Deps struct {
	dig.In
	Config     *config.Config
	ApiHandler api.ServerInterface
}

func New(deps Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount(BaseApiPath, api.Handler(deps.ApiHandler))

	applySwaggerEmbedded(r)
	return r
}

// This uses embedded swagger static files
func applySwaggerEmbedded(r chi.Router) {
	s, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n%s", err)
	}

	box := rice.MustFindBox("swagger").HTTPBox()
	embeddedFs := http.FileServer(box)

	r.Get(BaseSwaggerUIPath, func(w http.ResponseWriter, r *http.Request) {
		str, _ := box.String("index.html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(str))
	})
	r.Get(BaseSwaggerUIPath+"/api.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s)
	})
	r.Method(http.MethodGet, BaseSwaggerUIPath+"/*", http.StripPrefix("/swagger/", embeddedFs))
}
