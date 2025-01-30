package handler

import (
	"github.com/ras0q/go-backend-template/internal/repository"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) SetupRoutes(api *echo.Group) {
	// ping API
	pingAPI := api.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// // user API
	// userAPI := api.Group("/users")
	// {
	// 	userAPI.GET("", h.GetUsers)
	// 	userAPI.POST("", h.CreateUser)
	// 	userAPI.GET("/:userID", h.GetUser)
	// }

	// event API
	eventAPI := api.Group("/event")
	{
		eventAPI.GET("/all", h.GetEvents)
		eventAPI.POST("/new", h.CreateEvent)
		eventAPI.GET("/:id", h.GetEvent)
		eventAPI.GET("/update/:id", h.UpdateEvent)
		eventAPI.POST("/:id/approve", h.ApproveEvent)
		eventAPI.POST("/:id/reject", h.RejectEvent)
	}

	contactAPI := api.Group("/contact")
	{
		contactAPI.POST("", h.CreateContact)
	}
}
