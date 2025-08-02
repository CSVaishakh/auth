# Authentication Service Documentation

## Overview

The Authentication Service is a JWT-based authentication microservice built with Go and Fiber framework. It provides secure user registration, login, and logout functionality with role-based access control for the Office Management Platform.

## Architecture

### Technology Stack
- **Language**: Go 1.24.5
- **Web Framework**: Fiber v2.52.9
- **Database**: Supabase (PostgreSQL)
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Environment Management**: godotenv

### Design Patterns
- **Handler-Based Architecture**: Separate handlers for different authentication operations
- **Middleware Pattern**: Token verification middleware for protected routes
- **Helper Functions**: Utility functions for common operations
- **Type Safety**: Structured types for data models

## File Structure

```
auth-OMP/
├── main.go                 # Application entry point and route definitions
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── README.md               # Project overview
├── LICENSE                 # Project license
├── handlers/               # HTTP request handlers
│   ├── signUp.go          # User registration handler
│   ├── signIn.go          # User authentication handler
│   └── signOut.go         # User logout handler
├── utils/                # Utility functions
│   ├── initClient.go      # Supabase client initialization
│   ├── genJWT.go          # JWT token generation
│   ├── genUUID.go         # UUID generation
│   └── hashPass.go        # Password hashing and validation
├── middleware/             # HTTP middleware
│   └── middleware.go      # JWT token verification
└── types/                  # Data structure definitions
    └── types.go           # All type definitions
```

### File Descriptions

#### Core Files
- **main.go**: Application entry point that sets up Fiber server and defines API routes
- **go.mod**: Go module file with dependency management

#### Handlers Directory
- **signUp.go**: Handles user registration with role validation and password hashing
- **signIn.go**: Manages user authentication and JWT token generation
- **signOut.go**: Handles user logout by revoking JWT tokens

#### utils Directory
- **initClient.go**: Initializes and configures Supabase database client
- **genJWT.go**: Generates JWT tokens with custom claims and expiration
- **genUUID.go**: Generates unique identifiers for users and tokens
- **hashPass.go**: Provides password hashing and validation using bcrypt

#### Middleware Directory
- **middleware.go**: JWT token verification middleware for protected routes

#### Types Directory
- **types.go**: Contains all data structure definitions used across the application

## Database Schema

### Tables

#### users
Stores user account information:
```sql
- userid (string, primary key) - Unique user identifier
- email (string, unique) - User email address
- name (string) - Username/display name
- role (string) - User role in the system
- created_at (timestamp) - Account creation timestamp
```

#### secrets
Stores encrypted password data:
```sql
- userid (string, foreign key) - References users.userid
- password (string) - bcrypt hashed password
```

#### rolecodes
Defines available roles and their access codes:
```sql
- role (string) - Role name (e.g., "admin", "user", "manager")
- code (string) - Role access code for registration
```

#### jwt_tokens
Manages active JWT tokens:
```sql
- token_id (string, primary key) - Unique token identifier
- userid (string, foreign key) - References users.userid
- role (string) - User role at token creation
- token_type (string) - Token type (refresh/access)
- expires_at (string) - Token expiration timestamp
- issued_at (string) - Token creation timestamp
- status (boolean) - Token validity status
```

## API Documentation

### Base URL
```
http://localhost:5000
```

### Endpoints

#### POST /signup
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "username": "johndoe",
  "role_code": "USR001"
}
```

**Success Response (200):**
```json
{
  "message": "SignUp successful, Please Login"
}
```

**Error Responses:**
- `400`: Invalid request body
- `500`: Database error or role code validation failure

#### POST /signin
Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Success Response (200):**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "lifetime": 14400
}
```

**Error Responses:**
- `400`: Invalid request body
- `401`: Invalid credentials or user not found
- `500`: Database or token generation error

#### POST /signout
Revoke user's JWT token (logout).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "message": "SignOut successful"
}
```

**Error Responses:**
- `401`: Invalid or missing token
- `500`: Database error during token revocation

## Authentication Flow

### Registration Process
1. Client sends registration data with role code
2. System validates role code against `rolecodes` table
3. Password is hashed using bcrypt
4. User record created in `users` table
5. Hashed password stored in `secrets` table
6. Success response sent to client

### Login Process
1. Client sends email and password
2. System retrieves user data from `users` table
3. Password validation against stored hash
4. JWT token generated with user claims
5. Token metadata stored in `jwt_tokens` table
6. Token returned to client

### Logout Process
1. Client sends request with JWT token
2. Middleware validates token
3. Token status updated to false in database
4. Success response sent to client

### Token Verification
1. Extract token from Authorization header
2. Verify token signature using JWT_SECRET
3. Check token claims and expiration
4. Store user context in request locals
5. Allow request to proceed

## Environment Configuration

Create a `.env` file in the project root with the following variables:

```env
# Supabase Configuration
SUPABASE_URL=your_supabase_project_url
SUPABASE_ANON_KEY=your_supabase_anon_key

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Optional: Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=office_management
DB_USER=your_db_user
DB_PASSWORD=your_db_password
```

### Environment Variables Description

- **SUPABASE_URL**: Your Supabase project URL
- **SUPABASE_ANON_KEY**: Supabase anonymous key for database access
- **JWT_SECRET**: Secret key for JWT token signing and verification

## Security Features

### Password Security
- **bcrypt Hashing**: All passwords hashed with bcrypt using default cost
- **Salt Integration**: bcrypt automatically generates unique salts
- **Secure Comparison**: Constant-time password comparison to prevent timing attacks

### JWT Security
- **HS256 Signing**: Tokens signed using HMAC SHA-256 algorithm
- **Custom Claims**: Includes user ID, role, token type, and metadata
- **Expiration Control**: 4-hour token lifetime with automatic expiration
- **Token Revocation**: Ability to invalidate tokens on logout

### Role-Based Access
- **Role Validation**: Registration requires valid role codes
- **Role Persistence**: User roles stored and included in JWT claims
- **Access Control**: Middleware can check user roles for authorization

## Error Handling

### Common Error Responses

#### 400 Bad Request
- Invalid JSON in request body
- Missing required fields
- Malformed data

#### 401 Unauthorized
- Invalid credentials
- Expired or invalid JWT token
- Missing authorization header

#### 500 Internal Server Error
- Database connection failures
- JWT generation errors
- Unexpected system errors

### Error Response Format
```json
{
  "error": "Descriptive error message"
}
```

## Dependencies

### Production Dependencies
```go
github.com/gofiber/fiber/v2 v2.52.9      // Web framework
github.com/golang-jwt/jwt/v5 v5.2.3       // JWT implementation
github.com/joho/godotenv v1.5.1           // Environment variable loading
github.com/nedpals/supabase-go v0.5.0     // Supabase client
github.com/google/uuid v1.6.0             // UUID generation
golang.org/x/crypto v0.40.0               // Cryptographic functions
```

### Development Dependencies
```go
github.com/andybalholm/brotli v1.1.0      // Compression
github.com/go-viper/mapstructure/v2 v2.2.1 // Data mapping
github.com/google/go-querystring v1.1.0   // Query string encoding
```

## Performance Considerations

### Database Optimization
- Indexed email field for fast user lookups
- Separate secrets table to optimize user queries
- Token status indexing for quick validation

### Memory Management
- Efficient JWT claim structure
- Minimal data structures in types
- Proper error handling to prevent memory leaks

### Scalability
- Stateless JWT design enables horizontal scaling
- Database connection pooling through Supabase
- Middleware pattern allows easy feature extension

## Monitoring and Logging

### Built-in Logging
- Database operation logging
- User registration confirmations
- Error tracking with detailed messages

### Recommended Monitoring
- Response time tracking
- Error rate monitoring
- Database query performance
- JWT token usage patterns

## Integration Guidelines

### Client Integration
1. Implement proper error handling for all endpoints
2. Store JWT tokens securely (HTTP-only cookies recommended)
3. Include Authorization header for protected routes
4. Handle token expiration gracefully

### Microservice Integration
1. Use JWT tokens for service-to-service authentication
2. Validate tokens at service boundaries
3. Include user role information in service calls
4. Implement proper timeout and retry mechanisms

## Troubleshooting

### Common Issues

#### "Authheader not found"
- Ensure Authorization header is included in requests
- Verify header format: `Authorization: Bearer <token>`

#### "no valid user"
- Check email format and spelling
- Verify user exists in database
- Confirm email case sensitivity

#### "invalid token claims"
- Token may be corrupted or expired
- Verify JWT_SECRET environment variable
- Check token format and structure

#### Database Connection Errors
- Verify Supabase credentials in .env file
- Check network connectivity
- Confirm Supabase project status

### Debug Mode
Enable detailed logging by adding debug statements in handlers and middleware for development environments.
