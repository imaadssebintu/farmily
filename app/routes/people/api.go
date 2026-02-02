package people

import (
	"database/sql"
	"farmily/app/models"
	"farmily/app/routes/auth"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllPeopleAPI(c *fiber.Ctx, db *sql.DB) error {
	rows, err := db.Query(`
		SELECT id, first_name, middle_name, last_name, maiden_name, gender,
			birth_date, birth_place, death_date, death_place, is_living,
			occupation, biography, profile_photo_url, created_by, created_at, updated_at
		FROM people
		ORDER BY last_name, first_name
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch people",
		})
	}
	defer rows.Close()

	var people []models.PersonResponse
	for rows.Next() {
		var p models.Person
		err := rows.Scan(
			&p.ID, &p.FirstName, &p.MiddleName, &p.LastName, &p.MaidenName, &p.Gender,
			&p.BirthDate, &p.BirthPlace, &p.DeathDate, &p.DeathPlace, &p.IsLiving,
			&p.Occupation, &p.Biography, &p.ProfilePhotoURL, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			continue
		}
		people = append(people, p.ToResponse())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    people,
	})
}

func GetPersonAPI(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	personID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid person ID",
		})
	}

	var p models.Person
	err = db.QueryRow(`
		SELECT id, first_name, middle_name, last_name, maiden_name, gender,
			birth_date, birth_place, death_date, death_place, is_living,
			occupation, biography, profile_photo_url, created_by, created_at, updated_at
		FROM people WHERE id = $1
	`, personID).Scan(
		&p.ID, &p.FirstName, &p.MiddleName, &p.LastName, &p.MaidenName, &p.Gender,
		&p.BirthDate, &p.BirthPlace, &p.DeathDate, &p.DeathPlace, &p.IsLiving,
		&p.Occupation, &p.Biography, &p.ProfilePhotoURL, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Person not found",
		})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Database error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    p.ToResponse(),
	})
}

func CreatePersonAPI(c *fiber.Ctx, db *sql.DB) error {
	userID, err := auth.GetUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var req struct {
		FirstName       string  `json:"first_name"`
		MiddleName      *string `json:"middle_name"`
		LastName        string  `json:"last_name"`
		MaidenName      *string `json:"maiden_name"`
		Gender          string  `json:"gender"`
		BirthDate       *string `json:"birth_date"`
		BirthPlace      *string `json:"birth_place"`
		DeathDate       *string `json:"death_date"`
		DeathPlace      *string `json:"death_place"`
		IsLiving        bool    `json:"is_living"`
		Occupation      *string `json:"occupation"`
		Biography       *string `json:"biography"`
		ProfilePhotoURL *string `json:"profile_photo_url"`
		FatherID        *string `json:"father_id"`
		MotherID        *string `json:"mother_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Gender == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "First name, last name, and gender are required",
		})
	}

	personID := uuid.New()

	// Parse dates
	var birthDate, deathDate interface{}
	if req.BirthDate != nil && *req.BirthDate != "" {
		birthDate = *req.BirthDate
	}
	if req.DeathDate != nil && *req.DeathDate != "" {
		deathDate = *req.DeathDate
	}

	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to start transaction",
		})
	}

	_, err = tx.Exec(`
		INSERT INTO people (id, first_name, middle_name, last_name, maiden_name, gender,
			birth_date, birth_place, death_date, death_place, is_living,
			occupation, biography, profile_photo_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`, personID, req.FirstName, req.MiddleName, req.LastName, req.MaidenName, req.Gender,
		birthDate, req.BirthPlace, deathDate, req.DeathPlace, req.IsLiving,
		req.Occupation, req.Biography, req.ProfilePhotoURL, userID)

	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create person",
			"error":   err.Error(),
		})
	}

	// Handle parent relationships
	parents := []struct {
		ID   *string
		Type string
	}{
		{req.FatherID, "child"},
		{req.MotherID, "child"},
	}

	for _, p := range parents {
		if p.ID != nil && *p.ID != "" {
			parentID, err := uuid.Parse(*p.ID)
			if err == nil {
				relationshipID := uuid.New()
				_, err = tx.Exec(`
					INSERT INTO relationships (id, person1_id, person2_id, relationship_type)
					VALUES ($1, $2, $3, $4)
				`, relationshipID, personID, parentID, "child")

				if err != nil {
					tx.Rollback()
					return c.Status(500).JSON(fiber.Map{
						"success": false,
						"message": "Failed to create parent relationship",
					})
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to commit transaction",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Person created successfully",
		"id":      personID,
	})
}

func UpdatePersonAPI(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	personID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid person ID",
		})
	}

	var req struct {
		FirstName       string  `json:"first_name"`
		MiddleName      *string `json:"middle_name"`
		LastName        string  `json:"last_name"`
		MaidenName      *string `json:"maiden_name"`
		Gender          string  `json:"gender"`
		BirthDate       *string `json:"birth_date"`
		BirthPlace      *string `json:"birth_place"`
		DeathDate       *string `json:"death_date"`
		DeathPlace      *string `json:"death_place"`
		IsLiving        bool    `json:"is_living"`
		Occupation      *string `json:"occupation"`
		Biography       *string `json:"biography"`
		ProfilePhotoURL *string `json:"profile_photo_url"`
		FatherID        *string `json:"father_id"`
		MotherID        *string `json:"mother_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Parse dates
	var birthDate, deathDate interface{}
	if req.BirthDate != nil && *req.BirthDate != "" {
		birthDate = *req.BirthDate
	}
	if req.DeathDate != nil && *req.DeathDate != "" {
		deathDate = *req.DeathDate
	}

	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to start transaction",
		})
	}

	_, err = tx.Exec(`
		UPDATE people SET
			first_name = $1, middle_name = $2, last_name = $3, maiden_name = $4, gender = $5,
			birth_date = $6, birth_place = $7, death_date = $8, death_place = $9, is_living = $10,
			occupation = $11, biography = $12, profile_photo_url = $13, updated_at = $14
		WHERE id = $15
	`, req.FirstName, req.MiddleName, req.LastName, req.MaidenName, req.Gender,
		birthDate, req.BirthPlace, deathDate, req.DeathPlace, req.IsLiving,
		req.Occupation, req.Biography, req.ProfilePhotoURL, time.Now(), personID)

	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update person",
		})
	}

	// Sync parent relationships
	// For simplicity in this implementation:
	// 1. Delete all existing "child" relationships where this person is the child
	// 2. Add the new ones provided
	_, err = tx.Exec("DELETE FROM relationships WHERE person1_id = $1 AND relationship_type = 'child'", personID)
	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to sync relationships",
		})
	}

	parents := []struct {
		ID   *string
		Type string
	}{
		{req.FatherID, "child"},
		{req.MotherID, "child"},
	}

	for _, p := range parents {
		if p.ID != nil && *p.ID != "" {
			parentID, err := uuid.Parse(*p.ID)
			if err == nil {
				relationshipID := uuid.New()
				_, err = tx.Exec(`
					INSERT INTO relationships (id, person1_id, person2_id, relationship_type)
					VALUES ($1, $2, $3, $4)
				`, relationshipID, personID, parentID, "child")

				if err != nil {
					tx.Rollback()
					return c.Status(500).JSON(fiber.Map{
						"success": false,
						"message": "Failed to update parent relationship",
					})
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to commit transaction",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Person updated successfully",
	})
}

func DeletePersonAPI(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	personID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid person ID",
		})
	}

	_, err = db.Exec("DELETE FROM people WHERE id = $1", personID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete person",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Person deleted successfully",
	})
}

func SearchPeopleAPI(c *fiber.Ctx, db *sql.DB) error {
	query := c.Query("q")
	if query == "" {
		return GetAllPeopleAPI(c, db)
	}

	searchPattern := "%" + query + "%"
	rows, err := db.Query(`
		SELECT id, first_name, middle_name, last_name, maiden_name, gender,
			birth_date, birth_place, death_date, death_place, is_living,
			occupation, biography, profile_photo_url, created_by, created_at, updated_at
		FROM people
		WHERE first_name ILIKE $1 OR last_name ILIKE $1 OR middle_name ILIKE $1
		ORDER BY last_name, first_name
	`, searchPattern)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Search failed",
		})
	}
	defer rows.Close()

	var people []models.PersonResponse
	for rows.Next() {
		var p models.Person
		err := rows.Scan(
			&p.ID, &p.FirstName, &p.MiddleName, &p.LastName, &p.MaidenName, &p.Gender,
			&p.BirthDate, &p.BirthPlace, &p.DeathDate, &p.DeathPlace, &p.IsLiving,
			&p.Occupation, &p.Biography, &p.ProfilePhotoURL, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			continue
		}
		people = append(people, p.ToResponse())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    people,
	})
}
