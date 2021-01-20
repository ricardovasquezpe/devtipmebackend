package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/api/middlewares"
	"backend/api/responses"

	"backend/api/config"
)

type App struct {
	Router  *mux.Router
	MClient *mongo.Database
}

func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	//DBURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", DbUser, DbPassword, DbHost, DbName)
	//a.MClient = GetClient("mongodb://localhost:27017")
	a.MClient = config.GetDatabase("mongodb://localhost:27017", DbName)
	a.Router = mux.NewRouter()
	a.setVersionApi("v1")
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middlewares.SetContentTypeMiddleware)

	a.Router.HandleFunc("/", home).Methods("GET")
	a.Router.HandleFunc("/getusers", a.GetUsers).Methods("GET")
	a.Router.HandleFunc("/saveuser", a.SaveUser).Methods("POST")
	a.Router.HandleFunc("/deleteuser/{id}", a.DeleteUser).Methods("DELETE")
	a.Router.HandleFunc("/updateuser/{id}", a.UpdateUser).Methods("PUT")
}

func (a *App) RunServer() {
	log.Printf("\nServer starting on port 5000")
	log.Fatal(http.ListenAndServe(":5000", a.Router))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Ivents")
}

func (a *App) setVersionApi(v string) {
	a.Router = a.Router.PathPrefix("/api/" + v).Subrouter()
	//a.Router.Use(middle.MiddlewareOne)
}
