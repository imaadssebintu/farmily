package people

import (
	"database/sql"
	"farmily/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupPeopleRoutes(app *fiber.App, db *sql.DB) {
	// People pages
	app.Get("/people", auth.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.Render("people/index", fiber.Map{
			"Title":       "People - Farmily Tree",
			"CurrentPage": "people",
		})
	})

	app.Get("/people/:id", auth.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.Render("people/view", fiber.Map{
			"Title":       "Person Details - Farmily Tree",
			"CurrentPage": "people",
			"PersonID":    c.Params("id"),
		})
	})

	// People API
	api := app.Group("/api/people")
	api.Use(auth.AuthMiddleware)

	api.Get("/", func(c *fiber.Ctx) error {
		return GetAllPeopleAPI(c, db)
	})

	api.Get("/search", func(c *fiber.Ctx) error {
		return SearchPeopleAPI(c, db)
	})

	api.Get("/:id", func(c *fiber.Ctx) error {
		return GetPersonAPI(c, db)
	})

	api.Post("/", func(c *fiber.Ctx) error {
		return CreatePersonAPI(c, db)
	})

	api.Put("/:id", func(c *fiber.Ctx) error {
		return UpdatePersonAPI(c, db)
	})

	api.Delete("/:id", func(c *fiber.Ctx) error {
		return DeletePersonAPI(c, db)
	})
}
