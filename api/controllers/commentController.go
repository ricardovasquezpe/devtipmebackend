package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *App) SaveComment(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Comment created"}

	comment := &models.Comment{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	err = json.Unmarshal(body, &comment)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	comment.Prepare()
	err = comment.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user := r.Context().Value("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(user)
	comment.UserId = userId

	commentCreated, err := comment.SaveComment(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["comment"] = commentCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) FindAllComments(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	comments, err := models.FindAllComments(a.MClient, urlParams["solutionId"][0])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, comments)
	return
}
