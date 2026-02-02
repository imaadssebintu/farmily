package tree

import (
	"database/sql"
	"farmily/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupTreeRoutes(app *fiber.App, db *sql.DB) {
	// Tree Page
	app.Get("/tree", auth.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.Render("tree/index", fiber.Map{
			"Title":       "Family Tree - Farmily Tree",
			"CurrentPage": "tree",
		})
	})

	// Tree API
	app.Get("/api/tree/data", auth.AuthMiddleware, func(c *fiber.Ctx) error {
		return GetTreeDataAPI(c, db)
	})
}
