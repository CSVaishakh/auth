# Contributing Guide

## Welcome Contributors

Thank you for your interest in contributing to the Authentication Service! This guide will help you get started with contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Submission Process](#submission-process)
- [Issue Guidelines](#issue-guidelines)
- [Pull Request Process](#pull-request-process)
- [Code Review Process](#code-review-process)
- [Release Process](#release-process)

## Code of Conduct

### Our Pledge
We are committed to making participation in this project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Expected Behavior
- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior
- The use of sexualized language or imagery
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate in a professional setting

## Getting Started

### Prerequisites
- Go 1.24.5 or later
- Git
- Supabase account (for database)
- Basic understanding of JWT and authentication concepts

### Fork and Clone
1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/auth-OMP.git
   cd auth-OMP
   ```
3. Add the original repository as upstream:
   ```bash
   git remote add upstream https://github.com/CSVaishakh/auth-OMP.git
   ```

## Development Setup

### Environment Configuration
1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Fill in your environment variables:
   ```env
   SUPABASE_URL=your_supabase_url
   SUPABASE_ANON_KEY=your_supabase_key
   JWT_SECRET=your_jwt_secret
   ```

### Install Dependencies
```bash
go mod download
go mod tidy
```

### Database Setup
1. Create the required tables in your Supabase database:
   ```sql
   -- Users table
   CREATE TABLE users (
       userid VARCHAR PRIMARY KEY,
       email VARCHAR UNIQUE NOT NULL,
       name VARCHAR,
       role VARCHAR,
       created_at TIMESTAMP DEFAULT NOW()
   );

   -- Secrets table
   CREATE TABLE secrets (
       userid VARCHAR REFERENCES users(userid),
       password VARCHAR NOT NULL
   );

   -- Role codes table
   CREATE TABLE rolecodes (
       role VARCHAR NOT NULL,
       code VARCHAR UNIQUE NOT NULL
   );

   -- JWT tokens table
   CREATE TABLE jwt_tokens (
       token_id VARCHAR PRIMARY KEY,
       userid VARCHAR REFERENCES users(userid),
       role VARCHAR,
       token_type VARCHAR,
       expires_at VARCHAR,
       issued_at VARCHAR,
       status BOOLEAN DEFAULT true
   );
   ```

2. Insert sample role codes:
   ```sql
   INSERT INTO rolecodes (role, code) VALUES 
   ('admin', 'ADM001'),
   ('user', 'USR001'),
   ('manager', 'MGR001');
   ```

### Run the Application
```bash
go run main.go
```

The server will start on port 5000.

## Project Structure

### Directory Organization
```
auth-OMP/
├── handlers/          # HTTP request handlers
├── helpers/           # Utility functions
├── middleware/        # HTTP middleware
├── types/            # Data structures
├── tests/            # Test files (when added)
└── docs/             # Documentation
```

### Package Guidelines
- **handlers**: Contains HTTP route handlers
- **helpers**: Utility functions that can be reused
- **middleware**: HTTP middleware for request processing
- **types**: Data structure definitions

## Coding Standards

This project follows strict coding standards to ensure consistency and maintainability. For detailed conventions including naming standards, commit message guidelines, and code structure patterns, please refer to **[CONVENTIONS.md](CONVENTIONS.md)**.

### Quick Reference

#### Go Style Guide
Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

#### Key Conventions
- **Variables**: camelCase (e.g., `userID`, `tokenString`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Constants**: ALL_CAPS with underscores
- **Types**: PascalCase
- **Packages**: lowercase, single word when possible

#### Commit Message Format
```
type(scope): description

Examples:
feat(auth): add role-based access control
fix(handlers): resolve password validation issue
docs(readme): update installation instructions
```

For complete guidelines, see **[CONVENTIONS.md](CONVENTIONS.md)**.

### Code Formatting
```bash
# Format code
go fmt ./...

# Check for common issues
go vet ./...

# Run linter (if available)
golint ./...
```

### Error Handling
```go
// Good: Explicit error handling
result, err := someFunction()
if err != nil {
    return fiber.Map{"error": err.Error()}
}

// Avoid: Ignoring errors
result, _ := someFunction() // Don't do this
```

### Database Operations
```go
// Good: Proper error handling for database operations
query_err := client.DB.From("users").Insert(user).Execute(nil)
if query_err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": query_err.Error(),
    })
}
```

### JSON Response Format
```go
// Success responses
return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "message": "Operation successful",
    "data": responseData,
})

// Error responses
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error": "Descriptive error message",
})
```

## Testing Guidelines

### Test Structure
Create test files following Go conventions:
```
handlers/
├── signUp.go
├── signUp_test.go
├── signIn.go
└── signIn_test.go
```

### Test Categories
1. **Unit Tests**: Test individual functions
2. **Integration Tests**: Test handler endpoints
3. **Database Tests**: Test database operations

### Writing Tests
```go
func TestSignUp(t *testing.T) {
    // Setup
    app := fiber.New()
    app.Post("/signup", handlers.SignUp)

    // Test cases
    tests := []struct {
        name           string
        requestBody    string
        expectedStatus int
    }{
        {
            name:           "Valid signup",
            requestBody:    `{"email":"test@example.com","password":"password123","username":"testuser","role_code":"USR001"}`,
            expectedStatus: 200,
        },
        {
            name:           "Invalid email",
            requestBody:    `{"email":"invalid-email","password":"password123"}`,
            expectedStatus: 400,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/signup", strings.NewReader(tt.requestBody))
            req.Header.Set("Content-Type", "application/json")
            
            resp, _ := app.Test(req)
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
        })
    }
}
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestSignUp ./handlers
```

## Submission Process

### Before You Start
1. Check existing issues and pull requests
2. Create an issue to discuss major changes
3. Ensure your development environment is set up correctly

### Making Changes
1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. Make your changes following coding standards
3. Write or update tests for your changes
4. Update documentation if necessary
5. Commit your changes with clear messages

### Commit Messages
Follow conventional commit format as detailed in **[CONVENTIONS.md](CONVENTIONS.md)**:
```
type(scope): description
```

Available types: `feat`, `fix`, `docs`, `test`, `refactor`, `style`, `chore`

## Issue Guidelines

### Bug Reports
Include the following information:
- Go version
- Operating system
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages or logs
- Relevant code snippets

### Feature Requests
Include the following information:
- Clear description of the feature
- Use case and motivation
- Proposed implementation approach
- Potential impact on existing functionality

### Issue Templates
```markdown
**Bug Report**
- Go Version: 
- OS: 
- Description: 
- Steps to Reproduce: 
- Expected Behavior: 
- Actual Behavior: 
- Additional Context: 

**Feature Request**
- Description: 
- Motivation: 
- Proposed Solution: 
- Alternatives Considered: 
- Additional Context: 
```

## Pull Request Process

### PR Checklist
- [ ] Code follows project coding standards
- [ ] Tests pass locally
- [ ] Documentation updated if necessary
- [ ] Commit messages follow conventional format
- [ ] PR description explains the changes
- [ ] Related issues are linked

### PR Description Template
```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass
- [ ] New tests added
- [ ] Manual testing performed

## Related Issues
Fixes #(issue number)

## Additional Notes
Any additional information or context.
```

### Review Process
1. Automated checks must pass
2. At least one maintainer review required
3. Address review feedback
4. Maintainer approval required for merge

## Code Review Process

### For Reviewers
- Check code quality and standards compliance
- Verify test coverage
- Test functionality locally if needed
- Provide constructive feedback
- Approve when satisfied with changes

### For Contributors
- Respond to feedback promptly
- Make requested changes
- Ask questions if feedback is unclear
- Be open to suggestions and improvements

### Review Criteria
- **Functionality**: Does the code work as intended?
- **Code Quality**: Is the code clean and well-structured?
- **Performance**: Are there any performance concerns?
- **Security**: Are there any security vulnerabilities?
- **Testing**: Is there adequate test coverage?
- **Documentation**: Is documentation updated as needed?

## Release Process

### Version Numbering
Follow Semantic Versioning (SemVer):
- MAJOR.MINOR.PATCH
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

### Release Steps
1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release tag
4. Deploy to staging environment
5. Run integration tests
6. Deploy to production
7. Monitor for issues

## Getting Help

### Communication Channels
- GitHub Issues: For bugs and feature requests
- GitHub Discussions: For questions and general discussion
- Email: [Contact maintainers](mailto:maintainer@example.com)

### Documentation
- [README.md](README.md): Project overview
- [DOCUMENTATION.md](DOCUMENTATION.md): Detailed technical documentation
- [API Documentation](DOCUMENTATION.md#api-documentation): API reference

### Resources
- [Go Documentation](https://golang.org/doc/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [JWT.io](https://jwt.io/): JWT documentation
- [Supabase Documentation](https://supabase.com/docs)

## Recognition

### Contributors
We recognize and appreciate all contributors. Contributors will be:
- Listed in the repository contributors section
- Mentioned in release notes for significant contributions
- Invited to join the maintainer team for sustained contributions

### Types of Contributions
- Code contributions
- Documentation improvements
- Bug reports and testing
- Feature suggestions
- Code reviews
- Community support

Thank you for contributing to the Authentication Service!
