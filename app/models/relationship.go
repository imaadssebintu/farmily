package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Relationship types
const (
	RelationshipParent  = "parent"
	RelationshipChild   = "child"
	RelationshipSpouse  = "spouse"
	RelationshipSibling = "sibling"
)

type Relationship struct {
	ID               uuid.UUID      `json:"id"`
	Person1ID        uuid.UUID      `json:"person1_id"`
	Person2ID        uuid.UUID      `json:"person2_id"`
	RelationshipType string         `json:"relationship_type"`
	StartDate        sql.NullTime   `json:"start_date"`
	EndDate          sql.NullTime   `json:"end_date"`
	Notes            sql.NullString `json:"notes"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type RelationshipResponse struct {
	ID               uuid.UUID  `json:"id"`
	Person1ID        uuid.UUID  `json:"person1_id"`
	Person1Name      string     `json:"person1_name"`
	Person1Gender    string     `json:"person1_gender"`
	Person2ID        uuid.UUID  `json:"person2_id"`
	Person2Name      string     `json:"person2_name"`
	Person2Gender    string     `json:"person2_gender"`
	RelationshipType string     `json:"relationship_type"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	Notes            string     `json:"notes"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
