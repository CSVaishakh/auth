# Development Conventions

## Overview

This document outlines the coding conventions, naming standards, and commit message guidelines for the Authentication Service project. Following these conventions ensures code consistency, readability, and maintainability across the codebase.

## Table of Contents

- [Naming Conventions](#naming-conventions)
- [File and Directory Conventions](#file-and-directory-conventions)
- [Code Structure Conventions](#code-structure-conventions)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Branch Naming Conventions](#branch-naming-conventions)
- [Documentation Conventions](#documentation-conventions)
- [API Conventions](#api-conventions)

## Naming Conventions

### Go Language Conventions

#### Variables
```go
// Use camelCase for variables
var userID string
var tokenString string
var maxRetryCount int
var isAuthenticated bool

// Use descriptive names
var emailAddress string        // Good
var e string                   // Avoid

// Boolean variables should be questions
var isValid bool
var hasPermission bool
var canAccess bool
```

#### Functions and Methods
```go
// Exported functions: PascalCase
func SignUp(c *fiber.Ctx) error { }
func GenJWT(userID string) string { }
func ValidatePassword(password string) error { }

// Unexported functions: camelCase
func validateInput(data map[string]string) error { }
func hashPassword(password string) (string, error) { }
func initializeClient() (*supabase.Client, error) { }
```

#### Constants
```go
// Use ALL_CAPS with underscores
const JWT_SECRET_KEY = "jwt_secret"
const DEFAULT_TOKEN_EXPIRY = 4 * time.Hour
const MAX_RETRY_ATTEMPTS = 3
const DATABASE_CONNECTION_TIMEOUT = 30 * time.Second

// Group related constants
const (
    ROLE_ADMIN   = "admin"
    ROLE_USER    = "user"
    ROLE_MANAGER = "manager"
)

const (
    STATUS_ACTIVE   = "active"
    STATUS_INACTIVE = "inactive"
    STATUS_PENDING  = "pending"
)
```

#### Types and Structs
```go
// Use PascalCase for types
type User struct {
    UserID    string `json:"userid" db:"userid"`
    Email     string `json:"email" db:"email"`
    Username  string `json:"name" db:"name,omitempty"`
    Role      string `json:"role" db:"role,omitempty"`
    CreatedAt string `json:"created_at" db:"created_at,omitempty"`
}

type AuthRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int64  `json:"expires_in"`
}

// Interface names should describe behavior
type TokenGenerator interface {
    Generate(userID, role string) (string, error)
    Validate(token string) (*Claims, error)
}

type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(password, hash string) error
}
```

#### Packages
```go
// Use lowercase, single word when possible
package handlers
package utils  
package middleware
package types

// For multi-word packages, use lowercase without separators
package authhandlers  // Instead of auth_handlers or auth-handlers
package dbutils       // Instead of db_utils or db-utils
```

#### Interface Methods
```go
// Use verb-based names for interface methods
type UserService interface {
    CreateUser(user User) error
    GetUser(userID string) (*User, error)
    UpdateUser(userID string, updates User) error
    DeleteUser(userID string) error
    ListUsers(filters UserFilters) ([]User, error)
}

type TokenManager interface {
    IssueToken(userID, role string) (*Token, error)
    ValidateToken(tokenString string) (*Claims, error)
    RevokeToken(tokenID string) error
    RefreshToken(refreshToken string) (*Token, error)
}
```

### Database Conventions

#### Table Names
```sql
-- Use lowercase with underscores
users
secrets
role_codes  -- Instead of rolecodes
jwt_tokens
user_sessions
audit_logs
```

#### Column Names
```sql
-- Use lowercase with underscores
user_id
email_address
created_at
updated_at
is_active
token_expires_at

-- Be descriptive and consistent
first_name    -- Good
fname         -- Avoid
f_name        -- Avoid
```

#### Index Names
```sql
-- Format: idx_tablename_columnname(s)
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_jwt_tokens_user_id ON jwt_tokens(user_id);
CREATE INDEX idx_jwt_tokens_status_expires ON jwt_tokens(status, expires_at);
```

### JSON and API Conventions

#### JSON Field Names
```json
{
    "user_id": "abc123",
    "email_address": "user@example.com",
    "created_at": "2025-01-01T00:00:00Z",
    "is_active": true,
    "role_permissions": ["read", "write"]
}
```

#### HTTP Status Codes
```go
// Use appropriate status codes consistently
fiber.StatusOK          // 200 - Success
fiber.StatusCreated     // 201 - Resource created
fiber.StatusBadRequest  // 400 - Client error
fiber.StatusUnauthorized // 401 - Authentication required
fiber.StatusForbidden   // 403 - Permission denied
fiber.StatusNotFound    // 404 - Resource not found
fiber.StatusConflict    // 409 - Resource conflict
fiber.StatusInternalServerError // 500 - Server error
```

## File and Directory Conventions

### File Names
```
// Go files: lowercase with underscores
sign_up.go
sign_in.go
password_hasher.go
jwt_generator.go

// Test files: same name with _test suffix
sign_up_test.go
password_hasher_test.go

// Documentation: UPPERCASE
README.md
CONTRIBUTING.md
CONVENTIONS.md
LICENSE
```

### Directory Structure
```
auth-service/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go
├── internal/              # Private application code
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Business logic
│   ├── repositories/      # Data access layer
│   └── middleware/        # HTTP middleware
├── pkg/                   # Public library code
│   ├── auth/             # Authentication utilities
│   ├── database/         # Database utilities
│   └── logger/           # Logging utilities
├── configs/              # Configuration files
├── scripts/              # Build and deployment scripts
├── docs/                 # Documentation
└── tests/                # Integration tests
```

## Code Structure Conventions

### Function Organization
```go
// Order functions logically within files
// 1. Constants and variables
// 2. Type definitions
// 3. Constructor functions
// 4. Public methods
// 5. Private methods
// 6. Helper functions

// Example structure for handlers/sign_up.go
package handlers

import (...)

// Constants
const (
    MAX_SIGNUP_ATTEMPTS = 5
    SIGNUP_RATE_LIMIT   = time.Minute
)

// Variables
var signupAttempts = make(map[string]int)

// Public handler function
func SignUp(c *fiber.Ctx) error {
    // Implementation
}

// Private helper functions
func validateSignupRequest(data map[string]string) error {
    // Implementation
}

func createUserAccount(userData User) error {
    // Implementation
}
```

### Error Handling Patterns
```go
// Consistent error handling
func ProcessUser(userData User) error {
    if err := validateUser(userData); err != nil {
        return fmt.Errorf("user validation failed: %w", err)
    }
    
    if err := saveUser(userData); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}

// Error variable naming
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidPassword  = errors.New("invalid password")
    ErrTokenExpired     = errors.New("token has expired")
    ErrInsufficientRole = errors.New("insufficient role permissions")
)
```

### Import Organization
```go
package handlers

import (
    // Standard library imports first
    "context"
    "fmt"
    "time"
    
    // Third-party imports
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    
    // Local imports last
    "go-auth-app/utils"
    "go-auth-app/types"
)
```

## Commit Message Guidelines

### Conventional Commit Format
```
type(scope): description

[optional body]

[optional footer(s)]
```

### Commit Types
- **feat**: New feature for the user
- **fix**: Bug fix for the user
- **docs**: Documentation changes
- **style**: Code style changes (formatting, semicolons, etc.)
- **refactor**: Code refactoring without changing functionality
- **test**: Adding or updating tests
- **chore**: Maintenance tasks, dependency updates
- **perf**: Performance improvements
- **ci**: Continuous integration changes
- **build**: Build system or external dependency changes

### Scope Examples
- **auth**: Authentication-related changes
- **handlers**: HTTP handler changes
- **middleware**: Middleware changes
- **database**: Database-related changes
- **config**: Configuration changes
- **api**: API changes

### Commit Message Examples

#### Feature Commits
```
feat(auth): add two-factor authentication support

- Implement TOTP-based 2FA
- Add QR code generation for authenticator apps
- Update user registration flow

Closes #123
```

```
feat(handlers): add user profile management endpoints

- GET /profile - retrieve user profile
- PUT /profile - update user profile
- DELETE /profile - deactivate user account
```

#### Bug Fix Commits
```
fix(middleware): resolve JWT token validation race condition

The token validation was failing intermittently due to
concurrent access to the validation cache.

Fixes #456
```

```
fix(database): handle connection timeout errors gracefully

- Add retry logic for database operations
- Improve error messages for connection failures
- Update connection pool configuration
```

#### Documentation Commits
```
docs(api): update authentication endpoint documentation

- Add request/response examples
- Document error codes and messages
- Include rate limiting information
```

#### Refactoring Commits
```
refactor(utils): extract common validation logic

- Create shared validation utilities
- Reduce code duplication across handlers
- Improve test coverage for validation functions
```

#### Test Commits
```
test(handlers): add comprehensive signup handler tests

- Test valid signup scenarios
- Test validation error cases
- Test database error handling
- Add integration tests for full signup flow
```

### Breaking Changes
```
feat(api)!: redesign authentication response format

BREAKING CHANGE: The authentication response format has changed.
Previous format returned 'token', new format returns 'access_token'
and 'refresh_token' separately.

Before:
{
  "token": "...",
  "expires": 3600
}

After:
{
  "access_token": "...",
  "refresh_token": "...",
  "expires_in": 3600
}
```

### Multi-line Commit Messages
```
feat(auth): implement role-based access control

Add comprehensive RBAC system with the following features:
- Role hierarchy (admin > manager > user)
- Permission-based access control
- Dynamic role assignment
- Role inheritance

Implementation details:
- New roles table in database
- Updated JWT claims to include permissions
- Middleware for permission checking
- Admin endpoints for role management

This change maintains backward compatibility by defaulting
existing users to 'user' role.

Closes #789
Co-authored-by: Jane Smith <jane@example.com>
```

## Branch Naming Conventions

### Branch Types and Naming
```bash
# Feature branches
feature/user-authentication
feature/jwt-token-refresh
feature/role-based-permissions

# Bug fix branches
fix/password-validation-error
fix/database-connection-timeout
fix/memory-leak-in-middleware

# Hotfix branches (critical production fixes)
hotfix/security-vulnerability-patch
hotfix/database-connection-failure

# Release branches
release/v1.2.0
release/v2.0.0-beta

# Documentation branches
docs/api-documentation-update
docs/contributing-guidelines

# Refactoring branches
refactor/extract-validation-logic
refactor/improve-error-handling
```

### Branch Naming Rules
- Use lowercase letters
- Use hyphens to separate words
- Be descriptive but concise
- Include issue number when applicable: `feature/123-user-authentication`
- Use prefixes to categorize: `feature/`, `fix/`, `docs/`, etc.

## Documentation Conventions

### Markdown Standards
```markdown
# Main Title (H1)

## Section Title (H2)

### Subsection Title (H3)

#### Sub-subsection Title (H4)

- Use bullet points for lists
- Use **bold** for emphasis
- Use `code` for inline code
- Use ```language for code blocks

> Use blockquotes for important notes

| Table | Headers |
|-------|---------|
| Data  | Values  |
```

### Code Documentation
```go
// Package documentation
// Package handlers provides HTTP request handlers for the authentication service.
// It includes handlers for user registration, authentication, and session management.
package handlers

// Function documentation
// SignUp handles user registration requests.
// It validates the input data, creates a new user account, and returns
// a success response or appropriate error message.
//
// Parameters:
//   - c: Fiber context containing the HTTP request
//
// Returns:
//   - error: nil on success, error on failure
//
// Example request body:
//   {
//     "email": "user@example.com",
//     "password": "securepassword",
//     "username": "johndoe",
//     "role_code": "USR001"
//   }
func SignUp(c *fiber.Ctx) error {
    // Implementation
}

// Type documentation
// User represents a user account in the system.
// It contains basic user information and authentication details.
type User struct {
    // UserID is the unique identifier for the user
    UserID string `json:"userid" db:"userid"`
    
    // Email is the user's email address (used for login)
    Email string `json:"email" db:"email"`
    
    // Username is the display name for the user
    Username string `json:"name" db:"name,omitempty"`
    
    // Role defines the user's permissions level
    Role string `json:"role" db:"role,omitempty"`
    
    // CreatedAt is the timestamp when the account was created
    CreatedAt string `json:"created_at" db:"created_at,omitempty"`
}
```

## API Conventions

### Endpoint Naming
```
# RESTful endpoint patterns
POST   /auth/signup          # User registration
POST   /auth/signin          # User authentication
POST   /auth/signout         # User logout
POST   /auth/refresh         # Token refresh
GET    /auth/profile         # Get user profile
PUT    /auth/profile         # Update user profile
DELETE /auth/profile         # Delete user account

# Use kebab-case for multi-word resources
GET    /user-sessions        # List user sessions
POST   /password-reset       # Initiate password reset
PUT    /password-change      # Change password
```

### HTTP Headers
```
# Request headers
Content-Type: application/json
Authorization: Bearer <token>
X-Request-ID: <unique-request-id>

# Response headers
Content-Type: application/json
X-Response-Time: <response-time-ms>
X-Request-ID: <same-as-request>
```

### Response Format Standards
```json
// Success response
{
  "success": true,
  "data": {
    "user_id": "abc123",
    "email": "user@example.com"
  },
  "message": "User created successfully"
}

// Error response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": {
      "field": "email",
      "value": "invalid-email"
    }
  }
}

// List response
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "pages": 5
  }
}
```

## Validation and Enforcement

### Pre-commit Hooks
```bash
#!/bin/bash
# .git/hooks/pre-commit

# Format Go code
go fmt ./...

# Run go vet
go vet ./...

# Run tests
go test ./...

# Check commit message format
if ! grep -qE "^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .{1,50}" "$1"; then
    echo "Invalid commit message format"
    exit 1
fi
```

### Linting Configuration
```yaml
# .golangci.yml
linters:
  enable:
    - gofmt
    - goimports
    - govet
    - ineffassign
    - misspell
    - deadcode
    - varcheck
    - typecheck

linters-settings:
  gofmt:
    simplify: true
  goimports:
    local-prefixes: go-auth-app
```

By following these conventions, we ensure:
- Consistent code style across the project
- Better readability and maintainability
- Easier onboarding for new contributors
- Professional development practices
- Clear project history through commit messages
