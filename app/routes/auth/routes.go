package auth

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, db *sql.DB) {
	// Auth pages
	app.Get("/auth/login", func(c *fiber.Ctx) error {
		return c.Render("auth/login", fiber.Map{
			"Title": "Login - Farmily Tree",
		}, "layouts/auth")
	})

	app.Get("/auth/register", func(c *fiber.Ctx) error {
		return c.Render("auth/register", fiber.Map{
			"Title": "Register - Farmily Tree",
		}, "layouts/auth")
	})

	// Auth API
	app.Post("/api/auth/register", func(c *fiber.Ctx) error {
		return RegisterAPI(c, db)
	})

	app.Post("/api/auth/login", func(c *fiber.Ctx) error {
		return LoginAPI(c, db)
	})

	app.Post("/api/auth/logout", LogoutAPI)
}
