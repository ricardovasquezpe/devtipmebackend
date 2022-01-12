package controllers

import (
	"devtipmebackend/api/models"
	"devtipmebackend/api/responses"
	"devtipmebackend/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	err = json.Unmarshal(body, &solution)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	solution.Prepare()
	err = solution.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user := r.Context().Value("userID").(string)
	userId, _ := primitive.ObjectIDFromHex(user)
	solution.UserId = userId

	solutionCreated, err := solution.SaveSolution(a.MClient)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	for _, titleTopic := range solution.Topics {
		err := models.FindByTitleAndIncrease(a.MClient, titleTopic)
		if err != nil {
			topic := &models.Topic{
				Title: titleTopic,
				Total: 1,
			}

			topic.Prepare()
			topic.SaveTopic(a.MClient)
		}
	}

	resp["solution"] = solutionCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) GetSolutionById(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{}

	vars := mux.Vars(r)
	id, err := utils.Decrypt(vars["id"], os.Getenv("SECRET"))

	solutionFound, err := models.GetSolutionById(a.MClient, id)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	amount, err := models.GetTotalTipBySolutionId(a.MClient, id)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user, err := models.GetUserById(a.MClient, solutionFound.UserId.Hex())
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	user.Password = ""
	user.Email = ""

	resp["solution"] = solutionFound
	resp["amount"] = amount
	resp["user"] = user
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

	text := fmt.Sprint(result["text"])
	topic := fmt.Sprint(result["topic"])
	limit, _ := strconv.ParseInt(fmt.Sprint(result["limit"]), 10, 64)
	offset, _ := strconv.ParseInt(fmt.Sprint(result["offset"]), 10, 64)
	solutions, err := models.FindAllSolutions(a.MClient, text, limit, offset, topic)
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
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	defer file.Close()
	fileName := fmt.Sprint(utils.GetEpochTime())
	extension := filepath.Ext(header.Filename)

	resp["fileName"], err = a.S3.UploadImage(os.Getenv("AWS_BUCKET_NAME"), fileName+extension, file)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}

func (a *App) GetTrandingTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := models.GetTopicsLimited(a.MClient, 10)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, topics)
	return
}

func (a *App) GetMySolutions(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(string)
	solutions, err := models.GetSolutionsByUserId(a.MClient, userId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, solutions)
	return
}

func (a *App) UpdateSolutionStatus(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Solution updated successfully"}

	vars := mux.Vars(r)
	id := vars["id"]

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responses.ERROR(w, http.StatusOK, err)
		return
	}

	status := body["status"].(float64)

	err = models.UpdateSolutionStatus(a.MClient, id, status)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}
