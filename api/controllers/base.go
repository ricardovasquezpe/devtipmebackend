package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	"devtipmebackend/api/config"
	"devtipmebackend/api/middlewares"
	"devtipmebackend/api/responses"
)

type App struct {
	Router  *mux.Router
	MClient *mongo.Database
	S3      config.S3
	Mailer  config.Mailer
}

func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	//DBURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", DbUser, DbPassword, DbHost, DbName)
	//a.MClient = GetClient("mongodb://localhost:27017")
	a.MClient = config.GetDatabase("mongodb://localhost:27017", DbName)
	a.Router = mux.NewRouter()
	a.setVersionApi("v1")
	a.initializeRoutes()
}

func (a *App) InitializeS3Bucket(region, accessKeyId, accessKeySecret string) {
	a.S3 = config.NewS3(accessKeyId, accessKeySecret, region)
	a.S3.ConnectAws()
}

func (a *App) InitializeMailer(port, server, email, password string) {
	a.Mailer = config.NewMailer(port, server, email, password)
	a.Mailer.SetUpMailer()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middlewares.SetContentTypeMiddleware)

	a.Router.HandleFunc("/", home).Methods("GET")
	a.Router.HandleFunc("/user/getusers", a.GetAllUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.SaveUser).Methods("POST")
	a.Router.HandleFunc("/user/{id}", a.DeleteUser).Methods("DELETE")
	a.Router.HandleFunc("/user/{id}", a.UpdateUser).Methods("PUT")
	a.Router.HandleFunc("/user/login", a.Login).Methods("POST")
	a.Router.HandleFunc("/user/verify/{id}", a.VerifyUser).Methods("PUT")

	a.Router.HandleFunc("/callexternalapi", a.GetExternalData).Methods("GET")

	a.Router.HandleFunc("/solution/{id}", a.GetSolutionById).Methods("GET")
	a.Router.HandleFunc("/solution", a.GetAllSolutions).Methods("GET")
	a.Router.HandleFunc("/solution/find", a.FindAllSolutions).Methods("POST")

	a.Router.HandleFunc("/comment/find", a.FindAllComments).Methods("GET")

	s := a.Router.PathPrefix("/v1").Subrouter()
	s.Use(middlewares.AuthJwtVerify)
	s.HandleFunc("/solution", a.SaveSolution).Methods("POST")
	s.HandleFunc("/solution/uploadfile", a.uploadFile).Methods("POST")
	s.HandleFunc("/tip", a.SaveTip).Methods("POST")
	s.HandleFunc("/comment", a.SaveComment).Methods("POST")
	s.HandleFunc("/paypal/authorize", a.Authorize).Methods("POST")
}

func (a *App) RunServer() {
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Printf("\nServer starting on port 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(header, methods, origins)(a.Router)))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Ivents")
}

func (a *App) setVersionApi(v string) {
	a.Router = a.Router.PathPrefix("/api").Subrouter()
	//a.Router.Use(middle.MiddlewareOne)
}
