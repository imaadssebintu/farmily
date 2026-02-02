package main

import (
	"log"

	"farmily/app/config"
	"farmily/app/database"
	"farmily/app/routes/auth"
	"farmily/app/routes/dashboard"
	"farmily/app/routes/people"
	"farmily/app/routes/relationships"
	"farmily/app/routes/tree"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

// customErrorHandler handles HTTP errors with custom templates
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Check if this is an API request
	if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
		return c.Status(code).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"code":    code,
		})
	}

	// Handle different error codes for web requests
	switch code {
	case 404:
		return c.Status(404).Render("404", fiber.Map{
			"Title":       "Page Not Found - Farmily Tree",
			"CurrentPage": "",
		})
	case 403:
		return c.Status(403).Render("error", fiber.Map{
			"Title":        "Access Forbidden - Farmily Tree",
			"CurrentPage":  "",
			"ErrorCode":    "403",
			"ErrorTitle":   "Access Forbidden",
			"ErrorMessage": "You don't have permission to access this resource.",
		})
	case 401:
		return c.Status(401).Render("error", fiber.Map{
			"Title":        "Unauthorized - Farmily Tree",
			"CurrentPage":  "",
			"ErrorCode":    "401",
			"ErrorTitle":   "Unauthorized",
			"ErrorMessage": "Please log in to access this resource.",
		})
	case 500:
		return c.Status(500).Render("500", fiber.Map{
			"Title":        "Server Error - Farmily Tree",
			"CurrentPage":  "",
			"ErrorCode":    "500",
			"ErrorTitle":   "Internal Server Error",
			"ErrorMessage": "We're experiencing technical difficulties. Please try again later.",
			"ShowRetry":    true,
		})
	default:
		return c.Status(code).Render("error", fiber.Map{
			"Title":        "Error - Farmily Tree",
			"CurrentPage":  "",
			"ErrorCode":    code,
			"ErrorTitle":   "An Error Occurred",
			"ErrorMessage": err.Error(),
		})
	}
}

func main() {
	// Initialize database
	config.InitDB()
	defer config.GetDB().Close()

	// Run database migrations
	if err := database.RunMigrations(config.GetDB()); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize template engine
	engine := html.New("./app/templates", ".html")
	engine.Reload(true) // Enable template reloading for development
	engine.Debug(false) // Disable debug mode to reduce verbose logs

	// Create Fiber app
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
		ErrorHandler:      customErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Static files
	app.Static("/static", "./static")

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/auth/login")
	})

	// Setup auth routes
	auth.SetupAuthRoutes(app, config.GetDB())

	// Setup dashboard routes
	dashboard.SetupDashboardRoutes(app, config.GetDB())

	// Setup people routes
	people.SetupPeopleRoutes(app, config.GetDB())

	// Setup relationships routes
	relationships.SetupRelationshipsRoutes(app, config.GetDB())

	// Setup tree routes
	tree.SetupTreeRoutes(app, config.GetDB())

	// Catch-all route for 404 errors (must be last)
	app.Use("*", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "Page not found")
	})

	// Start server
	log.Println("Server starting on :8000")
	log.Fatal(app.Listen(":8000"))
}
