package auth

import (
	"database/sql"
	"farmily/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production")

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func RegisterAPI(c *fiber.Ctx, db *sql.DB) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.AuthResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(400).JSON(models.AuthResponse{
			Success: false,
			Message: "All fields are required",
		})
	}

	// Check if user already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return c.Status(500).JSON(models.AuthResponse{
			Success: false,
			Message: "Database error",
		})
	}

	if exists {
		return c.Status(400).JSON(models.AuthResponse{
			Success: false,
			Message: "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(models.AuthResponse{
			Success: false,
			Message: "Failed to hash password",
		})
	}

	// Create user
	userID := uuid.New()
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, req.Email, string(hashedPassword), req.FirstName, req.LastName)

	if err != nil {
		return c.Status(500).JSON(models.AuthResponse{
			Success: false,
			Message: "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := generateToken(userID, req.Email)
	if err != nil {
		return c.Status(500).JSON(models.AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	user := &models.User{
		ID:        userID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	return c.JSON(models.AuthResponse{
		Success: true,
		Message: "Registration successful",
		Token:   token,
		User:    user,
	})
}

func LoginAPI(c *fiber.Ctx, db *sql.DB) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.AuthResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Get user from database
	var user models.User
	err := db.QueryRow(`
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at
		FROM users WHERE email = $1
	`, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(401).JSON(models.AuthResponse{
			Success: false,
			Message: "Invalid email or password",
		})
	} else if err != nil {
		return c.Status(500).JSON(models.AuthResponse{
			Success: false,
			Message: "Database error",
		})
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.Status(401).JSON(models.AuthResponse{
			Success: false,
			Message: "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(models.AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(models.AuthResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    &user,
	})
}

func LogoutAPI(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

func generateToken(userID uuid.UUID, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
