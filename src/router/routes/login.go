package routes

import (
	"api/src/controllers"
	"net/http"
)

var login = Route{

	URI:      "/login",
	Method:   http.MethodPost,
	Funcao:   controllers.Login,
	NeedAuth: false,
}