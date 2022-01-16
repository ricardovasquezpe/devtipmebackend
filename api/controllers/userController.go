package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"devtipmebackend/utils"
	"encoding/json"
	"errors"
	"fmt"
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
		fmt.Println(err)
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	/*err = sendEmailVerify(a, user.Email, userCreated.ID.Hex())
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}*/

	token, err := utils.EncodeAuthToken(userCreated.ID)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

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
	resp["name"] = usr.Name
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
	_, err = user.UpdateUser(idDecrypt, a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}

func sendEmailVerify(a *App, email, id string) error {
	encriptedIdUser, err := utils.Encrypt([]byte(id), os.Getenv("SECRET"))

	if err != nil {
		return err
	}

	var data []models.TemplateData = []models.TemplateData{}
	data = append(data, models.TemplateData{
		Key:   "name",
		Value: "Ricardo",
	})

	url := fmt.Sprintf("%s/verifyme/%s", os.Getenv("SERVER_URL"), encriptedIdUser)
	//fmt.Println(url)
	data = append(data, models.TemplateData{
		Key:   "url",
		Value: url,
	})

	err = a.SendGridMailer.SendEmail([]string{"devtipmedeveloper@gmail.com"}, "SubJect Test", "d-52316f68e993473ba040673c6c8149c1", data)
	if err != nil {
		return err
	}

	return nil
}

func validateUserVerified(a *App, userId string) error {
	user, err := models.GetUserById(a.MClient, userId)
	if err != nil {
		return err
	}

	if user.Status == 0 {
		return errors.New("User is not verified")
	}

	return nil
}
