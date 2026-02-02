package relationships

import (
	"database/sql"
	"farmily/app/models"
	"farmily/app/routes/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetPersonRelationshipsAPI(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	personID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid person ID",
		})
	}

	rows, err := db.Query(`
		SELECT r.id, r.person1_id, r.person2_id, r.relationship_type,
			r.start_date, r.end_date, r.notes, r.created_at, r.updated_at,
			p1.first_name || ' ' || p1.last_name as person1_name,
			p2.first_name || ' ' || p2.last_name as person2_name,
			p1.gender as person1_gender,
			p2.gender as person2_gender
		FROM relationships r
		JOIN people p1 ON r.person1_id = p1.id
		JOIN people p2 ON r.person2_id = p2.id
		WHERE r.person1_id = $1 OR r.person2_id = $1
		ORDER BY r.created_at DESC
	`, personID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch relationships",
		})
	}
	defer rows.Close()

	var relationships []models.RelationshipResponse
	for rows.Next() {
		var r models.RelationshipResponse
		var startDate, endDate sql.NullTime
		var notes sql.NullString

		err := rows.Scan(
			&r.ID, &r.Person1ID, &r.Person2ID, &r.RelationshipType,
			&startDate, &endDate, &notes, &r.CreatedAt, &r.UpdatedAt,
			&r.Person1Name, &r.Person2Name,
			&r.Person1Gender, &r.Person2Gender,
		)
		if err != nil {
			continue
		}

		if startDate.Valid {
			r.StartDate = &startDate.Time
		}
		if endDate.Valid {
			r.EndDate = &endDate.Time
		}
		if notes.Valid {
			r.Notes = notes.String
		}

		// Determine correct label based on direction
		// r.RelationshipType describes what person1 is to person2.
		// We want displayType to describe what the OTHER person is to the current personID.
		displayType := r.RelationshipType
		if r.Person1ID == personID {
			// I am person1. The other is person2.
			// What is person2 to me? The inverse of what I am to them.
			displayType = getInverseRelationshipType(r.RelationshipType)
		} else {
			// I am person2. The other is person1.
			// What is person1 to me? What person1 is to person2.
			displayType = r.RelationshipType
		}
		r.RelationshipType = displayType

		relationships = append(relationships, r)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    relationships,
	})
}

func getInverseRelationshipType(originalType string) string {
	switch originalType {
	case "parent":
		return "child"
	case "child":
		return "parent"
	// spouse and sibling are symmetric
	default:
		return originalType
	}
}

func CreateRelationshipAPI(c *fiber.Ctx, db *sql.DB) error {
	_, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var req struct {
		Person1ID        string  `json:"person1_id"`
		Person2ID        string  `json:"person2_id"`
		RelationshipType string  `json:"relationship_type"`
		StartDate        *string `json:"start_date"`
		EndDate          *string `json:"end_date"`
		Notes            *string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	person1ID, err := uuid.Parse(req.Person1ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid person1 ID",
		})
	}

	person2ID, err := uuid.Parse(req.Person2ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid person2 ID",
		})
	}

	relationshipID := uuid.New()

	var startDate, endDate interface{}
	if req.StartDate != nil && *req.StartDate != "" {
		startDate = *req.StartDate
	}
	if req.EndDate != nil && *req.EndDate != "" {
		endDate = *req.EndDate
	}

	_, err = db.Exec(`
		INSERT INTO relationships (id, person1_id, person2_id, relationship_type, start_date, end_date, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, relationshipID, person1ID, person2ID, req.RelationshipType, startDate, endDate, req.Notes)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create relationship",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Relationship created successfully",
		"id":      relationshipID,
	})
}

func DeleteRelationshipAPI(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	relationshipID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid relationship ID",
		})
	}

	_, err = db.Exec("DELETE FROM relationships WHERE id = $1", relationshipID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete relationship",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Relationship deleted successfully",
	})
}

func SetupRelationshipsRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api/relationships")
	api.Use(auth.AuthMiddleware)

	api.Post("/", func(c *fiber.Ctx) error {
		return CreateRelationshipAPI(c, db)
	})

	api.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteRelationshipAPI(c, db)
	})

	// Get relationships for a specific person
	app.Get("/api/people/:id/relationships", auth.AuthMiddleware, func(c *fiber.Ctx) error {
		return GetPersonRelationshipsAPI(c, db)
	})
}
