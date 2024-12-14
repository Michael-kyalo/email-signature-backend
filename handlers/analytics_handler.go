package handlers

import (
	"context"
	"email-signature-backend/database"
	"github.com/gofiber/fiber/v2"
	"log"
)

type AnalyticsResponse struct {
	LinkID      string `json:"link_id"`
	TotalClicks int    `json:"total_clicks"`
	LastClicked string `json:"last_clicked"`
}

type CountResponse struct {
	Count int `json:"count"`
}

// GetAnalytics godoc
// @Summary Retrieve analytics for user links
// @Description Fetches the total clicks and last clicked timestamps for all links belonging to the user's signatures
// @Tags Analytics
// @Produce json
// @Success 200 {object} map[string][]AnalyticsResponse
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/analytics [get]
func GetAnalytics(c *fiber.Ctx) error {
	// Get user_id from context
	userID := c.Locals("user_id").(string)

	// Query to fetch click analytics
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT 
            links.id AS link_id, 
            COUNT(clicks.id) AS total_clicks,
            MAX(clicks.timestamp) AS last_clicked
         FROM links
         LEFT JOIN clicks ON clicks.link_id = links.id
         WHERE links.signature_id IN (
             SELECT id FROM signatures WHERE user_id = $1
         )
         GROUP BY links.id`,
		userID,
	)
	if err != nil {
		log.Printf("Failed to fetch analytics: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve analytics",
		})
	}
	defer rows.Close()

	// Parse the results
	analytics := []AnalyticsResponse{}
	for rows.Next() {
		var response AnalyticsResponse
		if err := rows.Scan(&response.LinkID, &response.TotalClicks, &response.LastClicked); err != nil {
			log.Printf("Failed to parse row: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse analytics data",
			})
		}
		analytics = append(analytics, response)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"analytics": analytics,
	})
}

// CountAnalyticsEntries godoc
// @Summary Get total analytics entries
// @Description Retrieve the total number of analytics entries (clicks).
// @Tags Analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CountResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /analytics/count [get]
func CountAnalyticsEntries(c *fiber.Ctx) error {
	// Query to count all analytics entries (clicks)
	var count int
	err := database.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM clicks").Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Failed to count analytics entries"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"count": count})
}
