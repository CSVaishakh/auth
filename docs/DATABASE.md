# Database Connections and Configuration

## Overview

The Authentication Service uses Supabase (PostgreSQL) as its primary database. This document outlines the database architecture, connection management, and configuration requirements.

## Database Architecture

### Supabase Integration

The application connects to Supabase using the official Go client library. All database operations are performed through the Supabase REST API, which provides:

- Automatic connection pooling
- Built-in security features
- Real-time capabilities (for future features)
- Row-level security policies
- Auto-generated API endpoints

### Connection Management

#### Client Initialization
```go
// helpers/initClient.go
func InItClient() (*supabase.Client, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, err
    }

    url := os.Getenv("SUPABASE_URL")
    key := os.Getenv("SUPABASE_ANON_KEY")

    client := supabase.CreateClient(url, key)
    return client, nil
}
```

#### Connection Pattern
- **Per-Request Initialization**: New client instance created for each request
- **Environment-Based Configuration**: Credentials loaded from environment variables
- **Error Handling**: Proper error propagation for connection failures

## Database Schema

### Table Structure

#### users
Primary user account information:
```sql
CREATE TABLE users (
    userid VARCHAR PRIMARY KEY,           -- 6-character UUID
    email VARCHAR UNIQUE NOT NULL,        -- User email (login identifier)
    name VARCHAR,                        -- Display name/username
    role VARCHAR,                        -- User role (admin, user, manager)
    created_at TIMESTAMP DEFAULT NOW()   -- Account creation timestamp
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

#### secrets
Encrypted password storage:
```sql
CREATE TABLE secrets (
    userid VARCHAR REFERENCES users(userid) ON DELETE CASCADE,
    password VARCHAR NOT NULL            -- bcrypt hashed password
);

-- Indexes
CREATE INDEX idx_secrets_userid ON secrets(userid);
```

#### rolecodes
Role validation codes:
```sql
CREATE TABLE rolecodes (
    role VARCHAR NOT NULL,               -- Role name (admin, user, manager)
    code VARCHAR UNIQUE NOT NULL         -- Registration code (ADM001, USR001, MGR001)
);

-- Sample data
INSERT INTO rolecodes (role, code) VALUES 
('admin', 'ADM001'),
('user', 'USR001'),
('manager', 'MGR001');
```

#### jwt_tokens
Token management and tracking:
```sql
CREATE TABLE jwt_tokens (
    token_id VARCHAR PRIMARY KEY,        -- Unique token identifier
    userid VARCHAR REFERENCES users(userid) ON DELETE CASCADE,
    role VARCHAR,                        -- User role at token creation
    token_type VARCHAR,                  -- Token type (refresh, access)
    expires_at VARCHAR,                  -- Expiration timestamp (Unix)
    issued_at VARCHAR,                   -- Creation timestamp (Unix)
    status BOOLEAN DEFAULT true          -- Token validity status
);

-- Indexes
CREATE INDEX idx_jwt_tokens_userid ON jwt_tokens(userid);
CREATE INDEX idx_jwt_tokens_status ON jwt_tokens(status);
CREATE INDEX idx_jwt_tokens_expires ON jwt_tokens(expires_at);
```

### Entity Relationships

```
users (1) ←→ (1) secrets
users (1) ←→ (N) jwt_tokens
rolecodes (1) ←→ (N) users (via role validation)
```

## Configuration

### Environment Variables

#### Required Variables
```env
# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key

# JWT Configuration  
JWT_SECRET=your-secure-jwt-secret-key
```

#### Supabase Setup
1. Create a new Supabase project
2. Navigate to Settings → API
3. Copy the Project URL and anon/public key
4. Configure Row Level Security (RLS) policies if needed

### Database Security

#### Row Level Security (RLS)
Recommended RLS policies for enhanced security:

```sql
-- Enable RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE secrets ENABLE ROW LEVEL SECURITY;
ALTER TABLE jwt_tokens ENABLE ROW LEVEL SECURITY;

-- Users can only access their own data
CREATE POLICY "Users can view own profile" ON users 
FOR SELECT USING (userid = current_user_id());

CREATE POLICY "Users can update own profile" ON users 
FOR UPDATE USING (userid = current_user_id());

-- Secrets are only accessible during authentication
CREATE POLICY "Secrets access during auth" ON secrets 
FOR SELECT USING (true); -- Controlled by application logic

-- JWT tokens are user-specific
CREATE POLICY "Users can view own tokens" ON jwt_tokens 
FOR ALL USING (userid = current_user_id());
```

#### API Key Security
- **anon key**: Used for client-side operations
- **service_role key**: For server-side operations (more privileged)
- Store keys securely in environment variables
- Rotate keys periodically

## Database Operations

### Connection Patterns

#### Standard Operation Pattern
```go
func DatabaseOperation(c *fiber.Ctx) error {
    // Initialize client
    client, err := helpers.InItClient()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // Perform database operation
    var result []TypeStruct
    queryErr := client.DB.From("table_name").
        Select("*").
        Execute(&result)
    
    if queryErr != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": queryErr.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(result)
}
```

#### Query Examples

**User Registration**:
```go
// Insert user
queryErr := client.DB.From("users").Insert(user).Execute(nil)

// Insert password hash
queryErr = client.DB.From("secrets").Insert(secret).Execute(nil)
```

**User Authentication**:
```go
// Find user by email
queryErr := client.DB.From("users").Select("*").Execute(&users)

// Get password hash
queryErr = client.DB.From("secrets").Select("*").
    Eq("userid", user.UserId).Execute(&storedHashs)
```

**Token Management**:
```go
// Store token
queryErr := client.DB.From("jwt_tokens").Insert(token).Execute(nil)

// Revoke token
queryErr := client.DB.From("jwt_tokens").
    Update(map[string]interface{}{"status": false}).
    Eq("token_id", tokenId).Execute(nil)
```

### Error Handling

#### Database Error Types
```go
// Connection errors
if strings.Contains(err.Error(), "connection") {
    // Handle connection issues
}

// Query errors
if strings.Contains(err.Error(), "syntax") {
    // Handle SQL syntax errors
}

// Constraint violations
if strings.Contains(err.Error(), "unique") {
    // Handle duplicate key errors
}
```

#### Retry Logic
```go
func retryDatabaseOperation(operation func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        if !isRetryableError(err) {
            return err
        }
        
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return errors.New("max retries exceeded")
}
```

## Performance Optimization

### Query Optimization

#### Efficient Queries
```go
// Good: Specific field selection
client.DB.From("users").Select("userid,email,role").
    Eq("email", email).Execute(&users)

// Avoid: Select all when not needed
client.DB.From("users").Select("*").Execute(&users)
```

#### Indexing Strategy
- Index frequently queried columns (email, userid, token_id)
- Composite indexes for multi-column queries
- Partial indexes for status-based queries

### Connection Pooling
Supabase automatically handles connection pooling:
- Default pool size: 15 connections
- Configurable through Supabase dashboard
- Monitor connection usage in Supabase metrics

### Caching Considerations
```go
// Role codes caching (rarely change)
var roleCodesCache map[string]string
var cacheExpiry time.Time

func getCachedRoleCodes() (map[string]string, error) {
    if time.Now().After(cacheExpiry) {
        // Refresh cache
        roleCodesCache = refreshRoleCodesFromDB()
        cacheExpiry = time.Now().Add(1 * time.Hour)
    }
    return roleCodesCache, nil
}
```

## Monitoring and Maintenance

### Database Monitoring

#### Supabase Dashboard Metrics
- Query performance
- Connection usage
- Storage utilization
- Error rates

#### Custom Monitoring
```go
// Log database operation times
start := time.Now()
err := client.DB.From("users").Insert(user).Execute(nil)
duration := time.Since(start)

log.Printf("Database operation took: %v", duration)
if duration > 5*time.Second {
    log.Printf("Slow query detected: %v", duration)
}
```

### Maintenance Tasks

#### Regular Cleanup
```sql
-- Clean up expired tokens (run periodically)
DELETE FROM jwt_tokens 
WHERE status = false 
AND CAST(expires_at AS BIGINT) < EXTRACT(epoch FROM NOW());

-- Archive old user sessions
DELETE FROM jwt_tokens 
WHERE CAST(expires_at AS BIGINT) < EXTRACT(epoch FROM NOW() - INTERVAL '30 days');
```

#### Backup Strategy
- Supabase provides automatic backups
- Point-in-time recovery available
- Export critical data periodically
- Test restore procedures

## Troubleshooting

### Common Connection Issues

#### Environment Variable Problems
```bash
# Check environment variables
echo $SUPABASE_URL
echo $SUPABASE_ANON_KEY

# Verify .env file loading
grep SUPABASE .env
```

#### Network Connectivity
```bash
# Test Supabase endpoint
curl -H "apikey: YOUR_ANON_KEY" \
     "https://your-project.supabase.co/rest/v1/users?select=*"
```

#### Authentication Failures
- Verify API keys are correct
- Check RLS policies
- Confirm database permissions

### Performance Issues

#### Slow Queries
- Check query execution plans
- Review index usage
- Monitor connection pool utilization
- Consider query optimization

#### Connection Timeouts
- Increase timeout values
- Check network latency
- Monitor connection pool status
- Consider connection retry logic

### Error Diagnostics

#### Database Error Logging
```go
func logDatabaseError(operation string, err error) {
    log.Printf("Database Error - Operation: %s, Error: %v", operation, err)
    
    // Additional context for debugging
    if supabaseErr, ok := err.(*supabase.Error); ok {
        log.Printf("Supabase Error Code: %d, Message: %s", 
            supabaseErr.Code, supabaseErr.Message)
    }
}
```

## Security Best Practices

### Credential Management
- Store sensitive data in environment variables
- Use separate keys for different environments
- Rotate keys regularly
- Monitor key usage

### Query Security
- Use parameterized queries through Supabase client
- Validate input data before database operations
- Implement proper error handling
- Log security-relevant events

### Access Control
- Implement RLS policies
- Use principle of least privilege
- Regular security audits
- Monitor database access patterns
