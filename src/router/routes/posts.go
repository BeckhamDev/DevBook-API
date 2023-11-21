package routes

import (
	"api/src/controllers"
	"net/http"
)

var postsRoutes = []Route{
	{
		URI:      "/Posts",
		Method:   http.MethodPost,
		Funcao:   controllers.NewPost,
		NeedAuth: true,
	},
	{
		URI:      "/Posts",
		Method:   http.MethodGet,
		Funcao:   controllers.GetPosts,
		NeedAuth: true,
	},
	{
		URI:      "/Posts/{idPost}",
		Method:   http.MethodGet,
		Funcao:   controllers.GetOnePost,
		NeedAuth: true,
	},
	{
		URI:      "/Posts/{idPost}",
		Method:   http.MethodPut,
		Funcao:   controllers.UpdatePost,
		NeedAuth: true,
	},
	{
		URI:      "/Posts/{idPost}",
		Method:   http.MethodDelete,
		Funcao:   controllers.DeletePost,
		NeedAuth: true,
	},
	{
		URI:      "/Users/{userID}/Posts",
		Method:   http.MethodGet,
		Funcao:   controllers.GetPostsByUser,
		NeedAuth: true,
	},
	{
		URI:      "/Posts/{postID}/Like",
		Method:   http.MethodPost,
		Funcao:   controllers.LikePost,
		NeedAuth: true,
	},
	{
		URI:      "/Posts/{postID}/Unlike",
		Method:   http.MethodPost,
		Funcao:   controllers.UnlikePost,
		NeedAuth: true,
	},
}