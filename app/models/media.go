package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Media types
const (
	MediaImage    = "image"
	MediaDocument = "document"
	MediaVideo    = "video"
)

type Media struct {
	ID          uuid.UUID      `json:"id"`
	PersonID    sql.NullString `json:"person_id"`
	EventID     sql.NullString `json:"event_id"`
	FilePath    string         `json:"file_path"`
	FileType    string         `json:"file_type"`
	Title       sql.NullString `json:"title"`
	Description sql.NullString `json:"description"`
	UploadDate  time.Time      `json:"upload_date"`
}

type Note struct {
	ID        uuid.UUID `json:"id"`
	PersonID  uuid.UUID `json:"person_id"`
	Content   string    `json:"content"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
