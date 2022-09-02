package main

import (
	"log"

	"github.com/1995parham-teaching/students/internal/config"
	"github.com/1995parham-teaching/students/internal/db"
	"github.com/1995parham-teaching/students/internal/handler"
	"github.com/1995parham-teaching/students/internal/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	cfg := config.New()

	app := echo.New()

	var logger *zap.Logger
	var err error
	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.New(cfg.Database)
	if err != nil {
		logger.Named("db").Fatal("cannot create a db instance", zap.Error(err))
	}

	var studentStore store.Student

	{
		logger := logger.Named("store")

		studentStore = store.NewStudentMongoDB(
			db, logger.Named("student"),
		)
	}

	{
		logger := logger.Named("http")

		h := handler.Student{
			Store:  studentStore,
			Logger: logger.Named("student"),
		}

		h.Register(app.Group("/api/students"))
	}

	app.Debug = cfg.Debug

	if err := app.Start(":1234"); err != nil {
		log.Println(err)
	}
}
