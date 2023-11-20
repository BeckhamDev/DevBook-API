package router

import (
	"api/src/router/routes"

	"github.com/gorilla/mux"
)

//Router retorna uma instancia de router com todas as rotas configuradas
func Router() *mux.Router {
	r:= mux.NewRouter()

	return routes.RouteConfig(r) 
}