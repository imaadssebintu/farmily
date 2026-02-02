package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID              uuid.UUID      `json:"id"`
	FirstName       string         `json:"first_name"`
	MiddleName      sql.NullString `json:"middle_name"`
	LastName        string         `json:"last_name"`
	MaidenName      sql.NullString `json:"maiden_name"`
	Gender          string         `json:"gender"`
	BirthDate       sql.NullTime   `json:"birth_date"`
	BirthPlace      sql.NullString `json:"birth_place"`
	DeathDate       sql.NullTime   `json:"death_date"`
	DeathPlace      sql.NullString `json:"death_place"`
	IsLiving        bool           `json:"is_living"`
	Occupation      sql.NullString `json:"occupation"`
	Biography       sql.NullString `json:"biography"`
	ProfilePhotoURL sql.NullString `json:"profile_photo_url"`
	CreatedBy       uuid.UUID      `json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// PersonResponse is used for API responses with formatted data
type PersonResponse struct {
	ID              uuid.UUID  `json:"id"`
	FirstName       string     `json:"first_name"`
	MiddleName      string     `json:"middle_name"`
	LastName        string     `json:"last_name"`
	MaidenName      string     `json:"maiden_name"`
	Gender          string     `json:"gender"`
	BirthDate       *time.Time `json:"birth_date"`
	BirthPlace      string     `json:"birth_place"`
	DeathDate       *time.Time `json:"death_date"`
	DeathPlace      string     `json:"death_place"`
	IsLiving        bool       `json:"is_living"`
	Occupation      string     `json:"occupation"`
	Biography       string     `json:"biography"`
	ProfilePhotoURL string     `json:"profile_photo_url"`
	DisplayName     string     `json:"display_name"`
	Age             *int       `json:"age"`
	Lifespan        string     `json:"lifespan"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// GetDisplayName returns the full name of the person
func (p *Person) GetDisplayName() string {
	name := p.FirstName
	if p.MiddleName.Valid && p.MiddleName.String != "" {
		name += " " + p.MiddleName.String
	}
	name += " " + p.LastName
	if p.MaidenName.Valid && p.MaidenName.String != "" {
		name += " (" + p.MaidenName.String + ")"
	}
	return name
}

// GetAge calculates the age of the person
func (p *Person) GetAge() *int {
	if !p.BirthDate.Valid {
		return nil
	}

	endDate := time.Now()
	if p.DeathDate.Valid {
		endDate = p.DeathDate.Time
	}

	age := endDate.Year() - p.BirthDate.Time.Year()
	if endDate.YearDay() < p.BirthDate.Time.YearDay() {
		age--
	}

	return &age
}

// GetLifespan returns a formatted lifespan string
func (p *Person) GetLifespan() string {
	if !p.BirthDate.Valid {
		return "Unknown"
	}

	birthYear := p.BirthDate.Time.Year()
	if p.DeathDate.Valid {
		deathYear := p.DeathDate.Time.Year()
		return fmt.Sprintf("%d - %d", birthYear, deathYear)
	}

	if p.IsLiving {
		return fmt.Sprintf("%d - Present", birthYear)
	}

	return fmt.Sprintf("%d - ?", birthYear)
}

// ToResponse converts Person to PersonResponse
func (p *Person) ToResponse() PersonResponse {
	resp := PersonResponse{
		ID:          p.ID,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Gender:      p.Gender,
		IsLiving:    p.IsLiving,
		DisplayName: p.GetDisplayName(),
		Age:         p.GetAge(),
		Lifespan:    p.GetLifespan(),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if p.MiddleName.Valid {
		resp.MiddleName = p.MiddleName.String
	}
	if p.MaidenName.Valid {
		resp.MaidenName = p.MaidenName.String
	}
	if p.BirthDate.Valid {
		resp.BirthDate = &p.BirthDate.Time
	}
	if p.BirthPlace.Valid {
		resp.BirthPlace = p.BirthPlace.String
	}
	if p.DeathDate.Valid {
		resp.DeathDate = &p.DeathDate.Time
	}
	if p.DeathPlace.Valid {
		resp.DeathPlace = p.DeathPlace.String
	}
	if p.Occupation.Valid {
		resp.Occupation = p.Occupation.String
	}
	if p.Biography.Valid {
		resp.Biography = p.Biography.String
	}
	if p.ProfilePhotoURL.Valid {
		resp.ProfilePhotoURL = p.ProfilePhotoURL.String
	}

	return resp
}
