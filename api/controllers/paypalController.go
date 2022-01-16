package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"devtipmebackend/api/services"
	"devtipmebackend/utils"
	"encoding/json"
	"errors"
	"net/http"
	"os"

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

	if body["orderId"] == nil || body["solutionId"] == nil || body["amount"] == nil || body["userIdTipped"] == nil {
		responses.ERROR(w, http.StatusOK, errors.New("Missing field"))
		return
	}

	userIdTipped, err := utils.Decrypt(body["userIdTipped"].(string), os.Getenv("SECRET"))
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}
	solutionId, err := utils.Decrypt(body["solutionId"].(string), os.Getenv("SECRET"))
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}
	orderId := body["orderId"].(string)
	amount := body["amount"].(float64)
	//amountString := fmt.Sprintf("%v", body["amount"].(float64))

	client, err := services.NewClient()
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	err = services.Authorize(client, orderId)
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

	userIdTippedObj, _ := primitive.ObjectIDFromHex(userIdTipped)
	tip.UserIdTipped = userIdTippedObj

	tip.PaypalId = orderId
	tip.Amount = amount

	tip.Prepare()

	err = tip.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	_, err = tip.SaveTip(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	//resp["tip"] = tipCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}
