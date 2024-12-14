package handlers

import (
	"context"
	"email-signature-backend/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClickRequest struct {
	LinkID    string `json:"link_id"`
	IPAddress string `json:"ip_address"`
}

// TrackClick godoc
// @Summary Track a click event
// @Description Logs a click event for a specific link, including the user's IP address
// @Tags Clicks
// @Accept json
// @Produce json
// @Param request body ClickRequest true "Click tracking payload"
// @Success 200 {object} map[string]string "Click tracked successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Failed to track click"
// @Router /api/track [post]
func TrackClick(c *fiber.Ctx) error {
	// Parse request body
	req := new(ClickRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Insert click into database
	clickID := uuid.New()
	_, err := database.DB.Exec(
		context.Background(),
		"INSERT INTO clicks (id, link_id, ip_address) VALUES ($1, $2, $3)",
		clickID,
		req.LinkID,
		req.IPAddress,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track click",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Click tracked successfully",
	})
}
