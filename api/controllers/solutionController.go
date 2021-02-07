package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) SaveSolution(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Solution created"}

	solution := &models.Solution{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &solution)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	solution.Prepare()
	err = solution.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	solutionCreated, err := solution.SaveSolution(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["solution"] = solutionCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) GetSolutionById(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{}

	vars := mux.Vars(r)
	id := vars["id"]

	solutionFound, err := models.GetSolutionById(a.MClient, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["solution"] = solutionFound
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) GetAllSolutions(w http.ResponseWriter, r *http.Request) {
	solutions, err := models.GetAllSolutions(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, solutions)
	return
}
