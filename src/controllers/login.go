package controllers

import (
	"api/src/auth"
	"api/src/db"
	"api/src/models"
	"api/src/repositories"
	"api/src/response"
	"api/src/security"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request){
	request, err := io.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(request, &user); err != nil {
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
	//valor recuperado do banco
	hashedUser, err := rep.SearchByEmail(user.Email)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.ValidatePassword(user.Password, hashedUser.Password); err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.CreateToken(hashedUser.ID)
	if err != nil {
		response.Erro(w, http.StatusInternalServerError, err)
		return
	}
	
	userID := strconv.FormatUint(hashedUser.ID, 10)

	response.JSON(w, http.StatusAccepted, models.AuthData{ID:userID,Token: token})
}