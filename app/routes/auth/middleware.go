package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Get token from cookie or Authorization header
	token := c.Cookies("token")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized - No token provided",
		})
	}

	// Parse and validate token
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized - Invalid token",
		})
	}

	// Store user info in context
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)

	return c.Next()
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID := c.Locals("userID")
	if userID == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}
	return userID.(uuid.UUID), nil
}
