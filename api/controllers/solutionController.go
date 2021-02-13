package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"devtipmebackend/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	user := r.Context().Value("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(user)
	solution.UserId = userId

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

	amount, err := models.GetTotalTipBySolutionId(a.MClient, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["solution"] = solutionFound
	resp["amount"] = amount
	responses.JSON(w, http.StatusOK, resp)
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

func (a *App) FindAllSolutions(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	str := fmt.Sprint(result["text"])
	limit, _ := strconv.ParseInt(fmt.Sprint(result["limit"]), 10, 64)
	offset, _ := strconv.ParseInt(fmt.Sprint(result["offset"]), 10, 64)
	solutions, err := models.FindAllSolutions(a.MClient, str, limit, offset)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, solutions)
	return
}

func (a *App) uploadFile(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "File uploaded successfully"}
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	defer file.Close()
	fileName := fmt.Sprint(utils.GetEpochTime())
	extension := filepath.Ext(header.Filename)

	resp["fileName"], err = a.S3.UploadImage("devtipme", fileName+extension, file)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}
