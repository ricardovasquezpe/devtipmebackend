package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"devtipmebackend/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func (a *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
	return
}

func (a *App) SaveUser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "User created"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user.Prepare()
	err = user.Validate("")

	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user.BeforeSave()
	userCreated, err := user.SaveUser(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	/*encriptedIdUser, _ := utils.Encrypt([]byte(userCreated.ID.Hex()), os.Getenv("SECRET"))
	//plaintext, _ := utils.Decrypt(encriptedIdUser, os.Getenv("SECRET"))

	err = a.Mailer.SendEmail(
		[]string{userCreated.Email},
		"Verify your account",
		"templates/template.html",
		map[string]string{"username": userCreated.Name, "url": "http://www.devoti.me/verifyaccount/" + encriptedIdUser})
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}*/

	token, err := utils.EncodeAuthToken(userCreated.ID)

	resp["token"] = token
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "User deleted successfully"}
	vars := mux.Vars(r)
	id := vars["id"]

	err := models.DeleteUser(id, a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "User updated successfully"}

	vars := mux.Vars(r)
	id := vars["id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userUpdate := models.User{}
	if err = json.Unmarshal(body, &userUpdate); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userUpdate.Prepare()

	_, err = userUpdate.UpdateUser(id, a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user.Prepare()

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, err := user.GetUserByEmail(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if usr == nil {
		resp["status"] = "failed"
		resp["error"] = "Login failed, please signup"
		responses.JSON(w, http.StatusOK, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, usr.Password)
	if err != nil {
		resp["status"] = "failed"
		resp["error"] = "Login failed, please try again"
		responses.JSON(w, http.StatusOK, resp)
		return
	}

	token, err := utils.EncodeAuthToken(usr.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
	return
}

func (a *App) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "User verify successfully"}
	vars := mux.Vars(r)
	id := vars["id"]

	idDecrypt, err := utils.Decrypt(id, os.Getenv("SECRET"))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	user, err := models.GetUserById(a.MClient, idDecrypt)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if user.Status != 0 {
		resp = map[string]interface{}{"status": "error", "message": "The user was already verified"}
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	user.Status = 1
	user.UpdatedAt = time.Now()
	user.UpdateUser(idDecrypt, a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}
