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
