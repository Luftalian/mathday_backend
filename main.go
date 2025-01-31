package main

import (
	"os"
	"strings"

	"github.com/ras0q/go-backend-template/internal/handler"
	"github.com/ras0q/go-backend-template/internal/migration"
	"github.com/ras0q/go-backend-template/internal/pkg/config"
	"github.com/ras0q/go-backend-template/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	allowOrigins := strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")

	// config.CORE_FRONTEND_URL
	allowOrigins = append(allowOrigins, config.CORE_FRONTEND_URL)

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
	}))

	// connect to database
	db, err := sqlx.Connect("mysql", config.MySQL().FormatDSN())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	// migrate tables
	if err := migration.MigrateTables(db.DB); err != nil {
		e.Logger.Fatal(err)
	}

	// setup repository
	repo := repository.New(db)

	// setup routes
	h := handler.New(repo)
	v1API := e.Group("/api/v1")
	h.SetupRoutes(v1API)

	e.Logger.Fatal(e.Start(config.AppAddr()))
}
