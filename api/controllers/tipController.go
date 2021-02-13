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
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &tip)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tip.Prepare()
	err = tip.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := r.Context().Value("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(user)
	tip.UserId = userId

	tipCreated, err := tip.SaveTip(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["tip"] = tipCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}
