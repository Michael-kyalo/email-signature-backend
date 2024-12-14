package handlers

import (
	"context"
	"email-signature-backend/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
)

type LinkRequest struct {
	SignatureID string `json:"signature_id"`
	URL         string `json:"url"`
}

// CreateLink godoc
// @Summary Create a new link for a signature
// @Description Adds a new link to an existing signature, ensuring the signature belongs to the authenticated user
// @Tags Links
// @Accept json
// @Produce json
// @Param request body LinkRequest true "Link creation payload"
// @Success 201 {object} map[string]string "Link created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized to add links to this signature"
// @Failure 500 {object} map[string]interface{} "Failed to create link"
// @Security BearerAuth
// @Router /api/links [post]
func CreateLink(c *fiber.Ctx) error {
	// Get user_id from context
	userID := c.Locals("user_id").(string)

	req := new(LinkRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Ensure the signature belongs to the user
	var count int
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM signatures WHERE id = $1 AND user_id = $2",
		req.SignatureID,
		userID,
	).Scan(&count)

	if err != nil || count == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized to add links to this signature",
		})
	}

	// Generate a new link ID
	linkID := uuid.New()

	// Insert the link into the database
	_, err = database.DB.Exec(
		context.Background(),
		"INSERT INTO links (id, signature_id, url) VALUES ($1, $2, $3)",
		linkID,
		req.SignatureID,
		req.URL,
	)
	if err != nil {
		log.Printf("Failed to insert link: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create link",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Link created successfully",
		"link_id": linkID,
	})
}
