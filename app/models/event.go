package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Event types
const (
	EventBirth      = "birth"
	EventDeath      = "death"
	EventMarriage   = "marriage"
	EventDivorce    = "divorce"
	EventGraduation = "graduation"
	EventEmployment = "employment"
	EventRetirement = "retirement"
	EventOther      = "other"
)

type Event struct {
	ID          uuid.UUID      `json:"id"`
	PersonID    uuid.UUID      `json:"person_id"`
	EventType   string         `json:"event_type"`
	EventDate   sql.NullTime   `json:"event_date"`
	EventPlace  sql.NullString `json:"event_place"`
	Description sql.NullString `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type EventResponse struct {
	ID          uuid.UUID  `json:"id"`
	PersonID    uuid.UUID  `json:"person_id"`
	PersonName  string     `json:"person_name"`
	EventType   string     `json:"event_type"`
	EventDate   *time.Time `json:"event_date"`
	EventPlace  string     `json:"event_place"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
