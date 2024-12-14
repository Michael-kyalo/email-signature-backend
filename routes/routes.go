package routes

import (
	"email-signature-backend/handlers"
	"email-signature-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Authentication routes
	api.Post("/register", handlers.RegisterUser)
	api.Post("/login", handlers.LoginUser)

	// Protected routes
	api.Post("/signature", middleware.Authenticate, handlers.CreateSignature)
	api.Get("/signature/:id/preview", middleware.Authenticate, handlers.PreviewSignature)
	api.Get("/signature/:id/export", middleware.Authenticate, handlers.ExportSignature)
	api.Post("/links", middleware.Authenticate, handlers.CreateLink)
	api.Post("/track", middleware.Authenticate, handlers.TrackClick)
	api.Get("/analytics", middleware.Authenticate, handlers.GetAnalytics)

	// Swagger routes
	SetupSwaggerRoutes(app)
}
