# Quick Start Guide - Farmily Tree 🌳

## Prerequisites
- PostgreSQL installed and running
- Database named `family` created

## Setup Steps

### 1. Create the Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE family;

# Exit
\q
```

### 2. Configure Environment (Optional)
To use a local database instead of the remote one:
```powershell
$env:LOCAL_DB="true"
```

### 3. Run the Application
```bash
cd c:\Users\ssebi\Desktop\projects\Newfolder\farmily
.\farmily.exe
```

Or for development with hot reload:
```bash
air
```

### 4. Access the Application
Open your browser and navigate to:
```
http://localhost:8080
```

## First Time Setup

1. **Register an Account**
   - Go to `http://localhost:8080/auth/register`
   - Fill in your details
   - Click "Register"

2. **Add Your First Family Member**
   - Click "People" in the sidebar
   - Click "+ Add Person"
   - Fill in the form (First Name, Last Name, and Gender are required)
   - Click "Add Person"

3. **Create Relationships**
   - Click on a person to view their details
   - Click "+ Add" in the Relationships section
   - Select relationship type and the other person
   - Save

## Features Available

✅ **Dashboard** - View statistics and quick actions
✅ **People Management** - Add, edit, delete, and search family members
✅ **Relationships** - Create and manage family connections
✅ **Person Profiles** - View detailed information about each family member
✅ **Search** - Find family members by name

## Troubleshooting

### Database Connection Failed
If you see a database connection error:
1. Make sure PostgreSQL is running
2. Verify the database `family` exists
3. Check the connection details in `app/config/config.go`
4. Try using local database: `$env:LOCAL_DB="true"`

### Port Already in Use
If port 8080 is already in use:
1. Stop the other application using port 8080
2. Or modify the port in `main.go` (line 238)

### Build Errors
If you encounter build errors:
```bash
go mod tidy
go build -o farmily.exe .
```

## Next Steps

- Add more family members
- Create relationships between them
- Explore the person detail pages
- Use the search functionality
- Check out the walkthrough document for more features

Enjoy building your family tree! 🌳
