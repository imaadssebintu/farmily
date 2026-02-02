-- ============================================
-- Family Tree Database Reset Script
-- ============================================
-- This script will:
-- 1. Drop all existing tables and data
-- 2. Recreate all tables with fresh schema
-- 
-- WARNING: This will DELETE ALL DATA in the database!
-- Make sure to backup your data before running this script.
-- ============================================

-- Drop all tables in reverse dependency order
-- (child tables first, then parent tables)
DROP TABLE IF EXISTS notes CASCADE;
DROP TABLE IF EXISTS media CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS relationships CASCADE;
DROP TABLE IF EXISTS people CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ============================================
-- Create Tables
-- ============================================

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create people table
CREATE TABLE people (
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
);

-- Create relationships table
CREATE TABLE relationships (
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
);

-- Create events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('birth', 'death', 'marriage', 'divorce', 'graduation', 'employment', 'retirement', 'other')),
    event_date DATE,
    event_place VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create media table
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id UUID REFERENCES people(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    file_path VARCHAR(500) NOT NULL,
    file_type VARCHAR(50) NOT NULL CHECK (file_type IN ('image', 'document', 'video')),
    title VARCHAR(255),
    description TEXT,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create notes table
CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- Create Indexes for Performance
-- ============================================

CREATE INDEX idx_people_name ON people(last_name, first_name);
CREATE INDEX idx_people_birth_date ON people(birth_date);
CREATE INDEX idx_relationships_person1 ON relationships(person1_id);
CREATE INDEX idx_relationships_person2 ON relationships(person2_id);
CREATE INDEX idx_events_person ON events(person_id);
CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_media_person ON media(person_id);
CREATE INDEX idx_notes_person ON notes(person_id);

-- ============================================
-- Verification
-- ============================================

-- Display all tables
SELECT 
    table_name,
    (SELECT COUNT(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
FROM information_schema.tables t
WHERE table_schema = 'public' 
    AND table_type = 'BASE TABLE'
ORDER BY table_name;
