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
	api.Get("/signatures", middleware.Authenticate, handlers.GetAllSignatures)           // Get all signatures
	api.Delete("/signature/:id", middleware.Authenticate, handlers.DeleteSignature)      // Delete a specific signature
	api.Get("/analytics/count", middleware.Authenticate, handlers.CountAnalyticsEntries) // Total analytics entries
	api.Get("/signatures/count", middleware.Authenticate, handlers.CountSignatures)      // Total signatures
	api.Get("/links/count", middleware.Authenticate, handlers.CountLinks)                // Total links

	api.Post("/links", middleware.Authenticate, handlers.CreateLink)
	api.Post("/track", middleware.Authenticate, handlers.TrackClick)
	api.Get("/analytics", middleware.Authenticate, handlers.GetAnalytics)

	// Swagger routes
	SetupSwaggerRoutes(app)
}
