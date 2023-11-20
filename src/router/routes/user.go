package routes

import (
	"api/src/controllers"
	"net/http"
)

var usersRoute = []Route{
	{
		URI:    "/users",
		Method: http.MethodPost,
		Funcao: controllers.CreateUser,
		NeedAuth: false,
	},
	{
		URI:    "/users",
		Method: http.MethodGet,
		Funcao: controllers.GetAllUsers,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}",
		Method: http.MethodGet,
		Funcao: controllers.GetOneUser,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}",
		Method: http.MethodPut,
		Funcao: controllers.UpdateUser,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}",
		Method: http.MethodDelete,
		Funcao: controllers.DeleteUser,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}/Follow",
		Method: http.MethodPost,
		Funcao: controllers.NewFollow,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}/StopFollowing",
		Method: http.MethodPost,
		Funcao: controllers.StopFollowing,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}/Followers",
		Method: http.MethodGet,
		Funcao: controllers.Followers,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}/Following",
		Method: http.MethodGet,
		Funcao: controllers.Following,
		NeedAuth: true,
	},
	{
		URI:    "/users/{userId}/NewPassword",
		Method: http.MethodPost,
		Funcao: controllers.NewPassword,
		NeedAuth: true,
	},
}