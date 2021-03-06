package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"devtipmebackend/api/services"
	"encoding/json"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *App) Authorize(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Authorize Completed"}

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	if body["orderId"] == nil || body["solutionId"] == nil || body["amount"] == nil{
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	orderId := body["orderId"].(string)
	solutionId := body["solutionId"].(string)
	amount := body["amount"].(float64)
	amountString := body["amount"].(string)

	client, err := services.NewClient()
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	err = services.Authorize(client, orderId, amountString)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	tip := &models.Tip{}

	solutionIdHex, _ := primitive.ObjectIDFromHex(solutionId)
	tip.SolutionId = solutionIdHex

	user := r.Context().Value("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(user)
	tip.UserId = userId
	tip.PaypalId = orderId
	tip.Amount = amount

	tip.Prepare()
	
	err = tip.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	tipCreated, err := tip.SaveTip(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	resp["tip"] = tipCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}
