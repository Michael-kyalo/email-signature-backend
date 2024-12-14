package handlers

import (
	"context"
	"crypto/sha256"
	"email-signature-backend/database"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Creates a new user account with a hashed password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration payload"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 500 {object} map[string]interface{} "Failed to register user"
// @Router /api/register [post]
func RegisterUser(c *fiber.Ctx) error {
	req := new(RegisterRequest)

	// Parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Hash the password
	hashedPassword := sha256.Sum256([]byte(req.Password))
	hashedPasswordHex := hex.EncodeToString(hashedPassword[:]) // Convert to hex string

	// Insert the user into the database
	userID := uuid.New()
	_, err := database.DB.Exec(
		context.Background(),
		"INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
		userID,
		req.Email,
		hashedPasswordHex,
	)
	if err != nil {
		log.Printf("Failed to register user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

// LoginUser godoc
// @Summary Authenticate a user
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login payload"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Failed to generate token"
// @Router /api/login [post]
func LoginUser(c *fiber.Ctx) error {
	req := new(LoginRequest)

	// Parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Hash the password and convert to hex
	hashedPassword := sha256.Sum256([]byte(req.Password))
	hashedPasswordHex := hex.EncodeToString(hashedPassword[:]) // Convert to hex string

	// Check if user exists
	var userID string
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT id FROM users WHERE email = $1 AND password = $2",
		req.Email,
		hashedPasswordHex,
	).Scan(&userID)
	if err != nil {
		log.Printf("Invalid credentials: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign token with secret
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Failed to sign token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": signedToken,
	})
}
