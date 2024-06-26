package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"healthchecker-app/handlers"
	"os"

	_ "healthchecker-app/docs"
)

// @title Health Checker API
// @version 1.0
// @description Це API для проекту Health Checker.
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /
func main() {
	e := echo.New()
	logrus.SetFormatter(&logrus.JSONFormatter{})

	db, err := connectDB()
	if err != nil {
		logrus.Fatal("Could not connect to database")
	}

	h := handlers.NewHandler(db)

	e.GET("/healthcheck", h.HealthCheck)
	e.GET("/questions", h.GetQuestions)
	e.POST("/submit", h.SubmitAnswers)
	e.Static("/", "static")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func connectDB() (*sqlx.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
