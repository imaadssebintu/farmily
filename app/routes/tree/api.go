package tree

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Node struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	PhotoURL  string `json:"photo_url,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	DeathDate string `json:"death_date,omitempty"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

func GetTreeDataAPI(c *fiber.Ctx, db *sql.DB) error {
	// check authentication or if user is logged in
	// ignoring strict auth check for now to allow easier testing if needed,
	// but typically we should use auth.GetUserID(c) if per-user trees were a thing.
	// Since it's a shared family tree, we just check if they are authenticated via middleware on the route group.

	// 1. Fetch All People (Nodes)
	rows, err := db.Query(`
		SELECT id, first_name, last_name, gender, profile_photo_url, birth_date, death_date
		FROM people
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch people",
		})
	}
	defer rows.Close()

	// Initialize as empty slices to ensure JSON [] instead of null
	nodes := []Node{}
	for rows.Next() {
		var id, firstName, lastName, gender string
		var photoURL, birthDate, deathDate sql.NullString

		if err := rows.Scan(&id, &firstName, &lastName, &gender, &photoURL, &birthDate, &deathDate); err != nil {
			continue
		}

		node := Node{
			ID:     id,
			Name:   firstName + " " + lastName,
			Gender: gender,
		}
		if photoURL.Valid {
			node.PhotoURL = photoURL.String
		}
		if birthDate.Valid {
			// Format date if needed, or send as is
			node.BirthDate = birthDate.String
		}
		if deathDate.Valid {
			node.DeathDate = deathDate.String
		}

		nodes = append(nodes, node)
	}

	// 2. Fetch All Relationships (Links)
	relRows, err := db.Query(`
		SELECT person1_id, person2_id, relationship_type
		FROM relationships
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch relationships",
		})
	}
	defer relRows.Close()

	links := []Link{}
	for relRows.Next() {
		var p1, p2, relType string
		if err := relRows.Scan(&p1, &p2, &relType); err != nil {
			continue
		}

		// D3 expects source and target to match node IDs
		links = append(links, Link{
			Source: p1,
			Target: p2,
			Type:   relType,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"nodes": nodes,
			"links": links,
		},
	})
}
