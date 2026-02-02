package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	log.Println("✓ Users table created/verified")

	// Create people table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS people (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			first_name VARCHAR(100) NOT NULL,
			middle_name VARCHAR(100),
			last_name VARCHAR(100) NOT NULL,
			maiden_name VARCHAR(100),
			gender VARCHAR(20) NOT NULL CHECK (gender IN ('Male', 'Female', 'Other')),
			birth_date DATE,
			birth_place VARCHAR(255),
			death_date DATE,
			death_place VARCHAR(255),
			is_living BOOLEAN DEFAULT true,
			occupation VARCHAR(255),
			biography TEXT,
			profile_photo_url VARCHAR(500),
			created_by UUID REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	log.Println("✓ People table created/verified")

	// Create relationships table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS relationships (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			person1_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
			person2_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
			relationship_type VARCHAR(50) NOT NULL CHECK (relationship_type IN ('parent', 'child', 'spouse', 'sibling')),
			start_date DATE,
			end_date DATE,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT different_people CHECK (person1_id != person2_id)
		)
	`)
	if err != nil {
		return err
	}
	log.Println("✓ Relationships table created/verified")

	// Create events table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
			event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('birth', 'death', 'marriage', 'divorce', 'graduation', 'employment', 'retirement', 'other')),
			event_date DATE,
			event_place VARCHAR(255),
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	log.Println("✓ Events table created/verified")

	// Create media table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS media (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			person_id UUID REFERENCES people(id) ON DELETE CASCADE,
			event_id UUID REFERENCES events(id) ON DELETE CASCADE,
			file_path VARCHAR(500) NOT NULL,
			file_type VARCHAR(50) NOT NULL CHECK (file_type IN ('image', 'document', 'video')),
			title VARCHAR(255),
			description TEXT,
			upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	log.Println("✓ Media table created/verified")

	// Create notes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_by UUID REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	log.Println("✓ Notes table created/verified")

	// Create indexes for better query performance
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_people_name ON people(last_name, first_name);
		CREATE INDEX IF NOT EXISTS idx_people_birth_date ON people(birth_date);
		CREATE INDEX IF NOT EXISTS idx_relationships_person1 ON relationships(person1_id);
		CREATE INDEX IF NOT EXISTS idx_relationships_person2 ON relationships(person2_id);
		CREATE INDEX IF NOT EXISTS idx_events_person ON events(person_id);
		CREATE INDEX IF NOT EXISTS idx_events_date ON events(event_date);
		CREATE INDEX IF NOT EXISTS idx_media_person ON media(person_id);
		CREATE INDEX IF NOT EXISTS idx_notes_person ON notes(person_id);
	`)
	if err != nil {
		return err
	}
	log.Println("✓ Database indexes created/verified")

	log.Println("All migrations completed successfully!")
	return nil
}
