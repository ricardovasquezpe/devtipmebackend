// main.go
package main

import (
	"devtipmebackend/api/controllers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := controllers.App{}
	app.Initialize(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	app.InitializeS3Bucket(os.Getenv("AWS_REGION"), os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"))
	//app.InitializeMailer(os.Getenv("MAILER_PORT"), os.Getenv("MAILER_SERVER"), os.Getenv("MAILER_EMAIL"), os.Getenv("MAILER_PASSWORD"))
	//app.InitializeGoMailer(os.Getenv("MAILER_PORT"), os.Getenv("MAILER_SERVER"), os.Getenv("MAILER_EMAIL"), os.Getenv("MAILER_PASSWORD"))
	app.InitializeSendgridMailer(os.Getenv("MAILER_EMAIL"), os.Getenv("MAILER_NAME"), os.Getenv("SENDGRID_APIKEY"), os.Getenv("SENDGRID_URL"), os.Getenv("SENDGRID_API"))
	app.RunServer()
}
