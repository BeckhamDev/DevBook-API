package controllers

import (
	"api/src/auth"
	"api/src/db"
	"api/src/models"
	"api/src/repositories"
	"api/src/response"
	"api/src/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	request, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err := json.Unmarshal(request,&user); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err := user.Prepare("register"); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewUserRep(db)
	user.ID, err = rep.Create(user)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	
	response.JSON(w, http.StatusCreated, user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	//value é o parametro vindo da requisição que será usado na busca
	value := strings.ToLower(r.URL.Query().Get("user"))

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewUserRep(db)
	user, err := rep.Search(value)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func GetOneUser(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	userID, err := strconv.ParseUint(value["userId"], 10, 64)
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

	rep := repositories.NewUserRep(db)
	user, err := rep.GetById(userID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	userID, err := strconv.ParseUint(value["userId"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	//Id do usuario que está salvo no token jwt
	userIdToken, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}
	
	if userID != userIdToken {
		response.Erro(w, http.StatusForbidden, errors.New("você só pode atualizar seu usuário"))
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err := json.Unmarshal(request, &user); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err := user.Prepare("edit"); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewUserRep(db)
	if err := rep.Update(userID, user); err != nil {
		response.Erro(w, http.StatusNoContent, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	userID, err := strconv.ParseUint(value["userId"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	//Id do usuario que está salvo no token jwt
	userIdToken, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}
	
	if userID != userIdToken {
		response.Erro(w, http.StatusForbidden, errors.New("você só pode excluir seu usuário"))
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewUserRep(db)
	err = rep.Delete(userID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func NewFollow(w http.ResponseWriter, r *http.Request) {
	followerID, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if userID == followerID {
		response.Erro(w, http.StatusForbidden, errors.New("você não pode seguir a si mesmo"))
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewUserRep(db)
	if err := rep.Follow(userID, followerID); err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func StopFollowing(w http.ResponseWriter, r *http.Request) {
	followerID, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if userID == followerID {
		response.Erro(w, http.StatusForbidden, errors.New("você não pode deixar de seguir a si mesmo"))
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewUserRep(db)
	if err := rep.StopFollowing(userID, followerID); err != nil{
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func Followers(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
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

	rep := repositories.NewUserRep(db)
	followers, err := rep.GetFollowersById(userID);
	if err != nil{
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}

func Following(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
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

	rep := repositories.NewUserRep(db)
	followers, err := rep.GetFollowing(userID);
	if err != nil{
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}

func NewPassword(w http.ResponseWriter, r *http.Request){
	followerID, err := auth.GetUserID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if userID != followerID {
		response.Erro(w, http.StatusForbidden, errors.New("você só pode alterar a sua própria senha"))
		return
	}

	corpo, err :=  io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var password models.UpdatePassword
	if err = json.Unmarshal(corpo, &password); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rep := repositories.NewUserRep(db)
	currentPassword, err := rep.GetCurrentPassword(userID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	//Aqui estou comparando a senha recuperada no banco (currentPassword) com a senha recebida na api que (OldPassword)
	if err := security.ValidatePassword(password.OldPassword,currentPassword); err != nil {
		response.Erro(w, http.StatusUnauthorized, errors.New("Teste"))
		return
	}

	newHashedPassword, err := security.Hash(password.NewPassword)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err := rep.UpdatePassword(userID,string(newHashedPassword) ); err != nil{
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}