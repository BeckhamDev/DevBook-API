package controllers

import (
	"api/src/auth"
	"api/src/db"
	"api/src/models"
	"api/src/repositories"
	"api/src/response"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func NewPost(w http.ResponseWriter, r*http.Request){
	userID, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	params, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	if err = json.Unmarshal(params, &post); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorID = userID

	if err = post.Prepare(); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	post.ID, err = rep.CreatePost(post)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

func GetPosts(w http.ResponseWriter, r*http.Request){
	userID, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	posts, err := rep.SearchPosts(userID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func GetOnePost(w http.ResponseWriter, r*http.Request){
	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["idPost"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	post, err := rep.GetOnePost(postID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func UpdatePost(w http.ResponseWriter, r*http.Request){
	userID, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["idPost"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	postFromDB, err := rep.GetOnePost(postID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if postFromDB.AuthorID != userID {
		response.Erro(w, http.StatusForbidden, errors.New("não é possivel atualizar uma publicação que não seja sua"))
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	if err := json.Unmarshal(request, &post); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err := post.Prepare(); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err := rep.Update(postID, post); err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DeletePost(w http.ResponseWriter, r*http.Request){
	userID, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["idPost"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	postFromDB, err := rep.GetOnePost(postID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if postFromDB.AuthorID != userID {
		response.Erro(w, http.StatusForbidden, errors.New("não é possivel deletar uma publicação que não seja sua"))
		return
	}

	if err := rep.DeletePost(postID); err != nil {
	response.Erro(w, http.StatusInternalServerError, err)
	return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	posts, err := rep.GetUserPosts(userID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func LikePost(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["postID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	if err := rep.Like(postID); err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DislikePost(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["postID"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewPostRep(db)
	if err := rep.Dislike(postID); err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}