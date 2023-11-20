package routes

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI      string              `json:"uri"`
	Method   string              `json:"method"`
	Funcao   func(http.ResponseWriter, *http.Request) 
	NeedAuth bool              `json:"need_auth"`
}

func RouteConfig(r *mux.Router) *mux.Router {
	routes := usersRoute
	routes = append(routes, login)
	routes = append(routes, postsRoutes...)

	for _, route := range routes {

		if route.NeedAuth {
			r.HandleFunc(route.URI, middlewares.Logger(middlewares.Authenticate(route.Funcao))).Methods(route.Method)
		}
		r.HandleFunc(route.URI, middlewares.Logger(route.Funcao)).Methods(route.Method)

	}

	return r
}