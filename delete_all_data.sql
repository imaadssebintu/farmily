-- ============================================
-- Family Tree Database - Delete All Data
-- ============================================
-- This script will DELETE ALL DATA from all tables
-- but keep the table structures intact.
-- 
-- WARNING: This will DELETE ALL DATA in the database!
-- Make sure to backup your data before running this script.
-- ============================================

-- Disable triggers temporarily to avoid constraint issues
SET session_replication_role = 'replica';

-- Delete data from all tables in reverse dependency order
-- (child tables first, then parent tables)
TRUNCATE TABLE notes CASCADE;
TRUNCATE TABLE media CASCADE;
TRUNCATE TABLE events CASCADE;
TRUNCATE TABLE relationships CASCADE;
TRUNCATE TABLE people CASCADE;
TRUNCATE TABLE users CASCADE;

-- Re-enable triggers
SET session_replication_role = 'origin';

-- ============================================
-- Verification
-- ============================================

-- Display row counts for all tables
SELECT 'users' as table_name, COUNT(*) as row_count FROM users
UNION ALL
SELECT 'people', COUNT(*) FROM people
UNION ALL
SELECT 'relationships', COUNT(*) FROM relationships
UNION ALL
SELECT 'events', COUNT(*) FROM events
UNION ALL
SELECT 'media', COUNT(*) FROM media
UNION ALL
SELECT 'notes', COUNT(*) FROM notes
ORDER BY table_name;
