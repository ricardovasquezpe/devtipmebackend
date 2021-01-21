package controllers

import (
	"devtipmebackend/api/responses"
	"encoding/json"
	"net/http"
)

func (a *App) GetExternalData(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, result)
		return
	}
}
