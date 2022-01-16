package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"net/http"
)

func (a *App) GetTrandingTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := models.GetTopicsLimited(a.MClient, 10)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, topics)
	return
}
