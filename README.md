# Farmily Tree 🌳

A modern web-based family tree management system built with Go Fiber and PostgreSQL. Track your family history, relationships, and life events with an intuitive interface.

## Features

- 👥 **People Management** - Add and manage family members with detailed biographical information
- 🔗 **Relationships** - Define parent-child, spouse, and sibling relationships
- 📅 **Events Timeline** - Record important life events (births, marriages, deaths, etc.)
- 🌲 **Family Tree Visualization** - Interactive family tree diagrams
- 🔍 **Search & Filter** - Quickly find family members
- 🔐 **Authentication** - Secure user accounts with JWT
- 📱 **Responsive Design** - Works on desktop, tablet, and mobile

## Tech Stack

- **Backend**: Go 1.21 with Fiber v2
- **Database**: PostgreSQL
- **Frontend**: HTML, CSS, JavaScript
- **Authentication**: JWT tokens
- **Template Engine**: Go HTML templates

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

## Installation

1. **Clone the repository**
   ```bash
   cd c:\Users\ssebi\Desktop\projects\Newfolder\farmily
   ```

2. **Create the database**
   ```bash
   # Connect to PostgreSQL
   psql -U postgres

   # Create database
   CREATE DATABASE family;
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Configure database connection**
   
   The application uses a remote PostgreSQL database by default. To use a local database:
   
   ```powershell
   $env:LOCAL_DB="true"
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

## Development

For hot-reload during development, install Air:

```bash
go install github.com/cosmtrek/air@latest
air
```

## Project Structure

```
farmily/
├── app/
│   ├── config/          # Database configuration
│   ├── database/        # Migrations and queries
│   ├── models/          # Data models
│   ├── routes/          # Route handlers
│   │   ├── auth/        # Authentication
│   │   ├── dashboard/   # Dashboard
│   │   ├── people/      # People management
│   │   └── relationships/ # Relationship management
│   └── templates/       # HTML templates
│       ├── layouts/     # Layout templates
│       ├── auth/        # Auth pages
│       ├── dashboard/   # Dashboard pages
│       └── people/      # People pages
├── static/
│   └── css/            # Stylesheets
├── main.go             # Application entry point
└── go.mod              # Go dependencies
```

## Database Schema

### Tables

- **users** - User accounts
- **people** - Family members
- **relationships** - Family relationships
- **events** - Life events
- **media** - Photos and documents
- **notes** - Personal notes

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout

### People
- `GET /api/people` - Get all people
- `GET /api/people/:id` - Get person by ID
- `POST /api/people` - Create person
- `PUT /api/people/:id` - Update person
- `DELETE /api/people/:id` - Delete person
- `GET /api/people/search?q=query` - Search people

### Relationships
- `POST /api/relationships` - Create relationship
- `DELETE /api/relationships/:id` - Delete relationship
- `GET /api/people/:id/relationships` - Get person's relationships

## Usage

1. **Register an account** at `/auth/register`
2. **Login** at `/auth/login`
3. **Add family members** from the People page
4. **Create relationships** between family members
5. **View your family tree** on the Tree page

## Contributing

This is a personal project, but suggestions and improvements are welcome!

## License

MIT License

## Author

Built with ❤️ for preserving family legacies
