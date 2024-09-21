package api

import (
	"boilerplate/api/middleware"
	"boilerplate/config"
	"boilerplate/gen"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

const (
	apiV1 = "/v1"
)

func NewServer(mainCtx context.Context, s config.Settings, si gen.ServerInterface) *http.Server {
	r := mux.NewRouter()

	r.PathPrefix(apiV1 + "/swaggerui").Handler(http.StripPrefix(apiV1+"/swaggerui", http.FileServer(http.Dir("./dist"))))
	gen.HandlerWithOptions(si, gen.GorillaServerOptions{
		BaseURL:          apiV1,
		BaseRouter:       r,
		Middlewares:      []gen.MiddlewareFunc{middleware.GetCheckAuth(s.JwtSecret)},
		ErrorHandlerFunc: nil,
	})

	fmt.Printf("--> localhost:%d%s\n", s.Port, "/v1/swaggerui/ - swagger")
	fmt.Printf("--> localhost:%d%s\n", s.Port, "/v1/ - api")
	return &http.Server{
		Addr: fmt.Sprintf(":%d", s.Port),
		BaseContext: func(listener net.Listener) context.Context {
			return mainCtx
		},
		Handler: r,
	}
}
