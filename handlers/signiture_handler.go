package handlers

import (
	"context"
	"email-signature-backend/database"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SignatureRequest struct {
	TemplateData map[string]interface{} `json:"template_data"`
}

// CreateSignature godoc
// @Summary Create a new email signature
// @Description Allows an authenticated user to create a new email signature
// @Tags Signatures
// @Accept json
// @Produce json
// @Param request body SignatureRequest true "Signature creation payload"
// @Success 201 {object} map[string]interface{} "Signature created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 500 {object} map[string]interface{} "Failed to create signature"
// @Security BearerAuth
// @Router /api/signature [post]
func CreateSignature(c *fiber.Ctx) error {
	// Get user_id from context
	userID := c.Locals("user_id").(string)

	req := new(SignatureRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Generate a new signature ID
	signatureID := uuid.New()

	// Insert the signature into the database
	_, err := database.DB.Exec(
		context.Background(),
		"INSERT INTO signatures (id, user_id, template_data) VALUES ($1, $2, $3)",
		signatureID,
		userID,
		req.TemplateData,
	)
	if err != nil {
		log.Printf("Failed to insert signature: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create signature",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Signature created successfully",
		"signature_id": signatureID,
	})
}

// ExportSignature godoc
// @Summary Export an email signature as HTML
// @Description Generates an HTML version of the specified signature for email clients
// @Tags Signatures
// @Param id path string true "Signature ID"
// @Param template query string false "Template type (basic or modern)"
// @Produce html
// @Success 200 {string} string "HTML representation of the signature"
// @Failure 404 {object} map[string]interface{} "Signature not found"
// @Failure 500 {object} map[string]interface{} "Failed to generate HTML"
// @Security BearerAuth
// @Router /api/signature/{id}/export [get]
func ExportSignature(c *fiber.Ctx) error {
	// Get user_id from context
	userID := c.Locals("user_id").(string)

	// Get signature_id from URL params
	signatureID := c.Params("id")

	// Optional: Get template type from query params (default to "basic")
	templateType := c.Query("template", "basic")

	// Fetch the signature data from the database
	var templateData map[string]interface{}
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT template_data FROM signatures WHERE id = $1 AND user_id = $2",
		signatureID,
		userID,
	).Scan(&templateData)
	if err != nil {
		log.Printf("Failed to fetch signature: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Signature not found",
		})
	}

	// Generate HTML based on template type
	var html string
	switch templateType {
	case "modern":
		html = renderModernHTML(templateData)
	default:
		html = renderBasicHTML(templateData)
	}

	// Return the HTML as a response
	return c.Status(fiber.StatusOK).SendString(html)
}

// PreviewSignature godoc
// @Summary Preview an email signature
// @Description Renders the signature in an HTML page for browser preview
// @Tags Signatures
// @Param id path string true "Signature ID"
// @Param template query string false "Template type (basic or modern)"
// @Produce html
// @Success 200 {string} string "HTML preview of the signature"
// @Failure 404 {object} map[string]interface{} "Signature not found"
// @Failure 500 {object} map[string]interface{} "Failed to generate preview"
// @Security BearerAuth
// @Router /api/signature/{id}/preview [get]
func PreviewSignature(c *fiber.Ctx) error {
	// Get user_id from context
	userID := c.Locals("user_id").(string)

	// Get signature_id from URL params
	signatureID := c.Params("id")

	// Optional: Get template type from query params (default to "basic")
	templateType := c.Query("template", "basic")

	// Fetch the signature data from the database
	var templateData map[string]interface{}
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT template_data FROM signatures WHERE id = $1 AND user_id = $2",
		signatureID,
		userID,
	).Scan(&templateData)
	if err != nil {
		log.Printf("Failed to fetch signature: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Signature not found",
		})
	}

	// Generate HTML based on the template type
	var signatureHTML string
	switch templateType {
	case "modern":
		signatureHTML = renderModernHTML(templateData)
	default:
		signatureHTML = renderBasicHTML(templateData)
	}

	// Wrap the signature in a full HTML document
	html := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Signature Preview</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    margin: 20px;
                    padding: 20px;
                    background-color: #f9f9f9;
                }
                .signature-container {
                    padding: 20px;
                    background: #fff;
                    border: 1px solid #ddd;
                    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
                    max-width: 600px;
                    margin: auto;
                }
            </style>
        </head>
        <body>
            <div class="signature-container">
                %s
            </div>
        </body>
        </html>
    `, signatureHTML)

	// Return the HTML response
	return c.Status(fiber.StatusOK).SendString(html)
}

// renderBasicHTML generates a basic HTML signature
func renderBasicHTML(data map[string]interface{}) string {
	// Same implementation as renderHTML above
	return renderHTML(data)
}

// renderModernHTML generates a modern HTML signature
func renderModernHTML(data map[string]interface{}) string {
	name := data["name"].(string)
	jobTitle := data["job_title"].(string)
	company := data["company"].(string)
	phone := data["phone"].(string)
	website := data["website"].(string)

	socialLinks := data["social_links"].(map[string]interface{})
	linkedin := socialLinks["linkedin"].(string)
	twitter := socialLinks["twitter"].(string)

	return fmt.Sprintf(`
        <div style="font-family: Verdana, sans-serif; color: #222; font-size: 16px; line-height: 1.8;">
            <table style="width: 100%%; border-spacing: 10px; background-color: #f9f9f9; padding: 10px;">
                <tr>
                    <td style="padding: 5px;">
                        <div style="font-size: 20px; font-weight: bold;">%s</div>
                        <div style="color: #555;">%s</div>
                        <div style="font-size: 12px; color: #777;">%s</div>
                    </td>
                </tr>
                <tr>
                    <td style="padding: 5px;">
                        <a href="tel:%s" style="color: #0a66c2; text-decoration: none; font-size: 14px;">Call: %s</a><br>
                        <a href="%s" style="color: #0a66c2; text-decoration: none; font-size: 14px;">Website: %s</a>
                    </td>
                </tr>
                <tr>
                    <td style="padding: 5px;">
                        <a href="%s" style="color: #0a66c2; text-decoration: none; margin-right: 15px;">LinkedIn</a>
                        <a href="%s" style="color: #0a66c2; text-decoration: none;">Twitter</a>
                    </td>
                </tr>
            </table>
        </div>
    `, name, jobTitle, company, phone, phone, website, website, linkedin, twitter)
}

// renderHTML generates an HTML string from the template data
func renderHTML(data map[string]interface{}) string {
	name := data["name"].(string)
	jobTitle := data["job_title"].(string)
	company := data["company"].(string)
	phone := data["phone"].(string)
	website := data["website"].(string)

	socialLinks := data["social_links"].(map[string]interface{})
	linkedin := socialLinks["linkedin"].(string)
	twitter := socialLinks["twitter"].(string)

	return fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; color: #444; font-size: 14px; line-height: 1.5;">
            <table>
                <tr>
                    <td>
                        <div style="font-size: 18px; font-weight: bold; color: #222;">%s</div>
                        <div style="color: #666;">%s</div>
                        <div style="color: #999; font-size: 12px;">%s</div>
                    </td>
                </tr>
                <tr>
                    <td>
                        <div style="margin-top: 10px;">
                            <a href="tel:%s" style="color: #0a66c2; text-decoration: none;">%s</a> | 
                            <a href="%s" style="color: #0a66c2; text-decoration: none;">%s</a>
                        </div>
                        <div style="margin-top: 10px;">
                            <a href="%s" style="color: #0a66c2; text-decoration: none; margin-right: 10px;">LinkedIn</a>
                            <a href="%s" style="color: #0a66c2; text-decoration: none;">Twitter</a>
                        </div>
                    </td>
                </tr>
            </table>
        </div>
    `, name, jobTitle, company, phone, phone, website, website, linkedin, twitter)
}
