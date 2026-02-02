package dashboard

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		return DashboardPage(c, db)
	})
}

func DashboardPage(c *fiber.Ctx, db *sql.DB) error {
	// Get statistics
	var totalPeople, livingPeople, deceasedPeople int
	db.QueryRow("SELECT COUNT(*) FROM people").Scan(&totalPeople)
	db.QueryRow("SELECT COUNT(*) FROM people WHERE is_living = true").Scan(&livingPeople)
	db.QueryRow("SELECT COUNT(*) FROM people WHERE is_living = false").Scan(&deceasedPeople)

	var totalRelationships int
	db.QueryRow("SELECT COUNT(*) FROM relationships").Scan(&totalRelationships)

	var totalEvents int
	db.QueryRow("SELECT COUNT(*) FROM events").Scan(&totalEvents)

	return c.Render("dashboard/index", fiber.Map{
		"Title":              "Dashboard - Farmily Tree",
		"CurrentPage":        "dashboard",
		"TotalPeople":        totalPeople,
		"LivingPeople":       livingPeople,
		"DeceasedPeople":     deceasedPeople,
		"TotalRelationships": totalRelationships,
		"TotalEvents":        totalEvents,
	})
}
