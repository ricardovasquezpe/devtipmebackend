package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"devtipmebackend/utils"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *App) SaveComment(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Comment created"}

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	if body["comment"] == nil || body["solutionId"] == nil {
		responses.ERROR(w, http.StatusOK, errors.New("Missing fields"))
		return
	}

	comment := &models.Comment{}
	solutionId, err := utils.Decrypt(body["solutionId"].(string), os.Getenv("SECRET"))
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	comment.SolutionId, _ = primitive.ObjectIDFromHex(solutionId)
	comment.Comment = body["comment"].(string)

	comment.Prepare()
	err = comment.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user := r.Context().Value("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(user)
	comment.UserId = userId

	_, err = comment.SaveComment(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["comment"] = comment
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) FindAllComments(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	solutionIdEncrypted := urlParams["solutionId"][0]
	solutionId, err := utils.Decrypt(solutionIdEncrypted, os.Getenv("SECRET"))
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	comments, err := models.FindAllComments(a.MClient, solutionId)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, comments)
	return
}
