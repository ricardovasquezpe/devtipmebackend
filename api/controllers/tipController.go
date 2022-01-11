package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *App) SaveTip(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Tip created"}

	tip := &models.Tip{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	err = json.Unmarshal(body, &tip)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	tip.Prepare()
	err = tip.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user := r.Context().Value("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(user)
	tip.UserId = userId

	tipCreated, err := tip.SaveTip(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	resp["tip"] = tipCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) GetMyTotalTips(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(string)
	totalTip, err := models.GetTotalTipByUserId(a.MClient, userId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, totalTip)
	return
}
