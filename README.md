# Authentication Service

JWT-based authentication microservice for the Office Management Platform.

## Features

- User registration and login
- JWT token authentication
- Password hashing with bcrypt
- Role-based access control
- Supabase database integration

## Tech Stack

- Go 1.24.5
- Fiber web framework
- JWT tokens
- Supabase database
- bcrypt password hashing

## API Endpoints

- `POST /signup` - User registration
- `POST /signin` - User login
- `POST /signout` - User logout (requires authentication)

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables (create `.env` file)
4. Run the service:
   ```bash
   go run main.go
   ```

The service runs on port 5000.

## Environment Variables

Configure the following environment variables:
- Database connection details/Supabase credentials
- JWT secret key

## Dependencies

- github.com/gofiber/fiber/v2
- github.com/golang-jwt/jwt/v5
- github.com/nedpals/supabase-go
- golang.org/x/crypto
- github.com/google/uuid

## Documentation

- **[Complete Documentation](docs/DOCUMENTATION.md)** - Detailed technical documentation
- **[Database Guide](docs/DATABASE.md)** - Database schema, connections, and configuration
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute to the project
- **[Development Conventions](docs/CONVENTIONS.md)** - Naming conventions and commit guidelines

## Quick Links

- [API Documentation](DOCUMENTATION.md#api-documentation)
- [Database Schema](DATABASE.md#database-schema)
- [Environment Setup](DATABASE.md#configuration)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Coding Standards](CONVENTIONS.md)
