# sqlk - SQL Database Library

A powerful SQL abstraction library with support for MySQL and PostgreSQL, featuring schema management, query builder, soft deletes, audit fields, bulk operations, and transaction helpers. Part of the konsultin backend boilerplate.

## Features

- **Multi-Database**: Support for MySQL and PostgreSQL
- **Schema Management**: Define schemas from structs or manually
- **Query Builder**: Fluent API for SELECT, INSERT, UPDATE, DELETE
- **Soft Delete**: Built-in soft delete support with `deleted_at`
- **Audit Fields**: Automatic tracking of created/updated timestamps and users
- **Bulk Operations**: Efficient bulk insert, update, and delete
- **Transactions**: Helper functions for safe transaction management
- **Connection Pooling**: Configurable connection pool settings

## Quick Start

```go
import "github.com/konsultin/project-goes-here/libs/sqlk"

// 1. Configure database
config := sqlk.Config{
    Driver:   "postgres",
    Host:     "localhost",
    Port:     5432,
    Username: "user",
    Password: "pass",
    DBName:   "mydb",
}

// 2. Create and initialize database
db, err := sqlk.NewDatabase(config)
if err != nil {
    log.Fatal(err)
}

if err := db.Init(); err != nil {
    log.Fatal(err)
}
defer db.Close()

// 3. Use with context
ctx := context.Background()
dbCtx := db.WithContext(ctx)
```

## Configuration

```go
type Config struct {
    Driver          string // "mysql" or "postgres"
    Host            string
    Port            int
    Username        string
    Password        string
    DBName          string
    MaxIdleConn     *int          // Default: 10
    MaxOpenConn     *int          // Default: 100
    MaxConnLifetime *int          // Seconds, Default: 3600
    Timeout         *int          // Seconds, Default: 30
}
```

## Schema Definition

### From Struct (Recommended)

```go
import "github.com/konsultin/project-goes-here/libs/sqlk/schema"

type User struct {
    ID        int       `db:"id"`
    Email     string    `db:"email"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// Create schema from struct
userSchema := schema.New(
    schema.WithModel(&User{}),
    schema.WithPrimaryKey("id"),
    schema.WithAutoIncrement(true),
)
```

### Manual Schema

```go
productSchema := schema.New(
    schema.WithTable("products"),
    schema.WithColumns("id", "name", "price", "stock"),
    schema.WithPrimaryKey("id"),
    schema.WithAutoIncrement(true),
)
```

### With Soft Delete

```go
userSchema := schema.New(
    schema.WithModel(&User{}),
    schema.WithPrimaryKey("id"),
    schema.WithSoftDelete(true, "deleted_at"),
)
```

### With Audit Fields

```go
userSchema := schema.New(
    schema.WithModel(&User{}),
    schema.WithPrimaryKey("id"),
    schema.WithAuditFields(true, "created_at", "updated_at", "created_by", "updated_by"),
)
```

## Query Builder

### SELECT Queries

```go
import "github.com/konsultin/project-goes-here/libs/sqlk/pq/query"

// Basic select
q := query.Select(userSchema, query.AllColumns).Build()
// SELECT "id", "email", "name" FROM "User"

// Select specific columns
q := query.Select(userSchema, "id", "email").Build()
// SELECT "id", "email" FROM "User"

// With WHERE
q := query.Select(userSchema, query.AllColumns).
    Where("email", op.Equal("john@example.com")).
    Build()

// With multiple conditions
q := query.Select(userSchema, query.AllColumns).
    Where("status", op.Equal("active")).
    Where("age", op.GreaterThan(18)).
    Build()

// With ORDER BY
q := query.Select(userSchema, query.AllColumns).
    OrderBy("created_at", query.DESC).
    Build()

// With LIMIT and OFFSET
q := query.Select(userSchema, query.AllColumns).
    Limit(10).
    Offset(20).
    Build()
```

### INSERT Queries

```go
// Insert all columns
q := query.Insert(userSchema, query.AllColumns).Build()
// INSERT INTO "User"("email", "name") VALUES (:email, :name) RETURNING "id"

// Insert specific columns
q := query.Insert(userSchema, "email", "name").Build()
```

### UPDATE Queries

```go
// Update all columns
q := query.Update(userSchema, query.AllColumns).
    Where("id", op.Equal(1)).
    Build()

// Update specific columns
q := query.Update(userSchema, "name", "email").
    Where("id", op.Equal(1)).
    Build()
```

### DELETE Queries

```go
// Hard delete
q := query.Delete(userSchema).
    Where("id", op.Equal(1)).
    Build()

// Soft delete (if enabled)
q := query.SoftDelete(userSchema).
    Where("id", op.Equal(1)).
    Build()
// UPDATE "User" SET "deleted_at" = NOW() WHERE "id" = :id_0
```

## Operators

```go
import "github.com/konsultin/project-goes-here/libs/sqlk/op"

// Comparison
op.Equal(value)
op.NotEqual(value)
op.GreaterThan(value)
op.GreaterThanOrEqual(value)
op.LessThan(value)
op.LessThanOrEqual(value)

// Pattern matching
op.Like(pattern)
op.NotLike(pattern)

// NULL checks
op.IsNull()
op.IsNotNull()

// IN clause
op.In(values...)
op.NotIn(values...)

// BETWEEN
op.Between(start, end)
```

## Bulk Operations

### Bulk Insert

```go
rows := []map[string]interface{}{
    {"email": "user1@example.com", "name": "User 1"},
    {"email": "user2@example.com", "name": "User 2"},
    {"email": "user3@example.com", "name": "User 3"},
}

q := query.BulkInsert(userSchema, query.AllColumns).
    Values(rows).
    Build()

// Execute
result, err := dbCtx.NamedExec(q, flattenBulkInsertParams(rows))
```

### Bulk Update

```go
rows := []map[string]interface{}{
    {"id": 1, "name": "Updated 1", "status": "active"},
    {"id": 2, "name": "Updated 2", "status": "inactive"},
}

q := query.BulkUpdate(userSchema, "name", "status").
    Values(rows).
    Build()
```

### Bulk Delete

```go
// Hard delete multiple IDs
q := query.BulkDelete(userSchema).
    IDs(1, 2, 3, 4, 5).
    Build()

// Soft delete multiple IDs
q := query.BulkSoftDelete(userSchema).
    IDs(1, 2, 3, 4, 5).
    Build()
```

## Transactions

### Using WithTransaction (Recommended)

```go
import "github.com/konsultin/project-goes-here/libs/sqlk"

err := sqlk.WithTransaction(ctx, db.conn, nil, func(tx *sqlk.TxContext) error {
    // Insert user
    _, err := tx.Tx().NamedExec(insertUserQuery, user)
    if err != nil {
        return err // Auto rollback
    }
    
    // Insert profile
    _, err = tx.Tx().NamedExec(insertProfileQuery, profile)
    if err != nil {
        return err // Auto rollback
    }
    
    return nil // Auto commit
})
```

### Manual Transaction Control

```go
tx, err := sqlk.BeginTx(ctx, db.conn, nil)
if err != nil {
    return err
}

_, err = tx.Tx().Exec(query)
if err != nil {
    return sqlk.CommitOrRollback(tx, err)
}

return sqlk.CommitOrRollback(tx, nil)
```

## Production Examples

### Repository Pattern

```go
type UserRepository struct {
    db     *sqlk.Database
    schema *schema.Schema
}

func NewUserRepository(db *sqlk.Database) *UserRepository {
    userSchema := schema.New(
        schema.WithModel(&User{}),
        schema.WithPrimaryKey("id"),
        schema.WithAutoIncrement(true),
        schema.WithSoftDelete(true, "deleted_at"),
    )
    
    return &UserRepository{
        db:     db,
        schema: userSchema,
    }
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    q := query.Select(r.schema, query.AllColumns).
        Where("id", op.Equal(id)).
        Build()
    
    var user User
    err := r.db.WithContext(ctx).Get(&user, q, map[string]interface{}{"id_0": id})
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    q := query.Insert(r.schema, query.AllColumns).Build()
    
    result, err := r.db.WithContext(ctx).NamedExec(q, user)
    if err != nil {
        return err
    }
    
    id, _ := result.LastInsertId()
    user.ID = int(id)
    return nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
    q := query.Update(r.schema, query.AllColumns).
        Where("id", op.Equal(user.ID)).
        Build()
    
    _, err := r.db.WithContext(ctx).NamedExec(q, user)
    return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
    q := query.SoftDelete(r.schema).
        Where("id", op.Equal(id)).
        Build()
    
    _, err := r.db.WithContext(ctx).Exec(q, map[string]interface{}{"id_0": id})
    return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
    q := query.Select(r.schema, query.AllColumns).
        OrderBy("created_at", query.DESC).
        Limit(limit).
        Offset(offset).
        Build()
    
    var users []User
    err := r.db.WithContext(ctx).Select(&users, q, nil)
    return users, err
}
```

### Complex Queries

```go
// Search with multiple filters
func (r *UserRepository) Search(ctx context.Context, filters SearchFilters) ([]User, error) {
    builder := query.Select(r.schema, query.AllColumns)
    
    if filters.Email != "" {
        builder.Where("email", op.Like("%"+filters.Email+"%"))
    }
    
    if filters.Status != "" {
        builder.Where("status", op.Equal(filters.Status))
    }
    
    if filters.AgeMin > 0 {
        builder.Where("age", op.GreaterThanOrEqual(filters.AgeMin))
    }
    
    q := builder.
        OrderBy("created_at", query.DESC).
        Limit(filters.Limit).
        Build()
    
    var users []User
    err := r.db.WithContext(ctx).Select(&users, q, builder.Params())
    return users, err
}
```

### Using Transactions

```go
func (s *UserService) CreateUserWithProfile(ctx context.Context, user *User, profile *Profile) error {
    return sqlk.WithTransaction(ctx, s.db.conn, nil, func(tx *sqlk.TxContext) error {
        // Insert user
        userQ := query.Insert(s.userSchema, query.AllColumns).Build()
        result, err := tx.Tx().NamedExec(userQ, user)
        if err != nil {
            return err
        }
        
        userID, _ := result.LastInsertId()
        user.ID = int(userID)
        
        // Insert profile
        profile.UserID = user.ID
        profileQ := query.Insert(s.profileSchema, query.AllColumns).Build()
        _, err = tx.Tx().NamedExec(profileQ, profile)
        return err
    })
}
```

## Best Practices

- **Use Schemas**: Define schemas once and reuse throughout the application
- **Context Support**: Always pass context for timeout/cancellation
- **Connection Pooling**: Configure pool size based on your load
- **Soft Delete**: Enable for entities that shouldn't be permanently deleted
- **Transactions**: Use `WithTransaction` for automatic rollback on errors
- **Query Builder**: Use builder for complex dynamic queries
- **Repository Pattern**: Encapsulate data access logic in repositories
- **Close Connections**: Always `defer db.Close()` after initialization
- **Error Handling**: Check and handle all database errors properly 
- **Prepared Statements**: Use named parameters (`:param`) to prevent SQL injection

## Database Drivers

Supported drivers:
- `mysql` - MySQL database
- `postgres` - PostgreSQL database

Both drivers are automatically imported.
