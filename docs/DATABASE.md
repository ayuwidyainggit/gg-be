# Database Documentation

Dokumentasi struktur database untuk Scylla Backend.

## 📋 Table of Contents

- [Database Overview](#database-overview)
- [Connection Configuration](#connection-configuration)
- [Schema Structure](#schema-structure)
- [Migration](#migration)
- [Best Practices](#best-practices)

## 🗄️ Database Overview

### Database System

- **Type**: PostgreSQL
- **Version**: 12+
- **ORM**: GORM (Go Object-Relational Mapping)

### Database Connection

Setiap service terhubung ke database PostgreSQL yang sama atau terpisah tergantung konfigurasi.

## ⚙️ Connection Configuration

### Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=password
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600
```

### Connection Pool Settings

- **Max Connections**: Maximum open connections (default: 100)
- **Max Idle Connections**: Maximum idle connections (default: 10)
- **Max Lifetime**: Connection max lifetime in seconds (default: 3600)

### Connection String Format

```
host={DB_HOST} port={DB_PORT} user={DB_USER} dbname={DB_NAME} password={DB_PASS}
```

## 📐 Schema Structure

### Schema Naming

- **Schema**: `public` (default PostgreSQL schema)
- **Tables**: `snake_case` (contoh: `jobs`, `users`, `products`)

### Table Naming Convention

- Plural form: `jobs`, `users`, `products`
- Descriptive names: `account_payables`, `warehouse_stocks`
- Junction tables: `user_roles`, `product_categories`

### Common Table Structure

#### Standard Fields

Kebanyakan table memiliki field berikut:

```sql
id              UUID PRIMARY KEY DEFAULT gen_random_uuid()
created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
deleted_at      TIMESTAMP NULL (soft delete)
```

#### Active Status

Beberapa table memiliki field `active`:

```sql
active          BOOLEAN DEFAULT true
```

#### Audit Fields

Beberapa table memiliki audit fields:

```sql
created_by      VARCHAR(100)
updated_by      VARCHAR(100)
```

## 🔄 Migration

### GORM Auto Migration

GORM dapat melakukan auto migration untuk model:

```go
db.AutoMigrate(&model.Job{})
```

### Manual Migration

Untuk migration manual, gunakan SQL files di folder `migration/`:

```sql
-- migration/001_create_jobs_table.sql
CREATE TABLE IF NOT EXISTS public.jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_name VARCHAR(100) NOT NULL,
    job_desc VARCHAR(255) NOT NULL,
    job_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Migration Best Practices

1. **Version Control**: Semua migration harus di-version control
2. **Reversible**: Migration harus bisa di-rollback jika diperlukan
3. **Testing**: Test migration di development environment dulu
4. **Backup**: Backup database sebelum migration di production

## 📊 Model Examples

### Job Model (Cronjob Service)

```go
type Job struct {
    JobID            uuid.UUID `gorm:"column:job_id;primaryKey;type:uuid;default:gen_random_uuid()"`
    JobName          string    `gorm:"column:job_name"`
    JobDesc          string    `gorm:"column:job_desc"`
    JobType          string    `gorm:"column:job_type;not null"`
    CronExpression   string    `gorm:"column:cron_expression"`
    DayOfWeekOrMonth int       `gorm:"column:day_of_week_or_month;type:integer"`
    TimeOfDay        *string   `gorm:"column:time_of_day;type:text"`
    RunAt            *string  `gorm:"column:run_at;type:text"`
    Task             string   `gorm:"column:task;not null"`
    Url              string   `gorm:"column:url"`
    Payload          string   `gorm:"column:payload"`
    Active           bool     `gorm:"column:active;default:true"`
    CreatedBy        string   `gorm:"column:created_by"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (Job) TableName() string {
    return "public.jobs"
}
```

### Common Model Patterns

#### UUID Primary Key
```go
ID uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
```

#### Timestamps
```go
CreatedAt time.Time
UpdatedAt time.Time
DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete
```

#### Foreign Keys
```go
UserID uuid.UUID `gorm:"column:user_id;type:uuid"`
User   User      `gorm:"foreignKey:UserID"`
```

#### Indexes
```go
Email string `gorm:"column:email;uniqueIndex"`
Name  string `gorm:"column:name;index"`
```

## 🔒 Transaction Management

### Transaction Pattern

```go
err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
    // All database operations here use txCtx
    if err := service.Repo1.Store(txCtx, data1); err != nil {
        return err
    }
    if err := service.Repo2.Store(txCtx, data2); err != nil {
        return err
    }
    return nil
})
```

### Transaction Context

Repository harus check context untuk transaction:

```go
func (repo *Repository) model(ctx context.Context) *gorm.DB {
    tx := extractTx(ctx)
    if tx != nil {
        return tx.WithContext(ctx)
    }
    return repo.WithContext(ctx)
}
```

## 💡 Best Practices

### 1. Use Transactions for Multiple Operations

```go
// ✅ Good
service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
    repo1.Store(txCtx, data1)
    repo2.Store(txCtx, data2)
    return nil
})

// ❌ Bad
repo1.Store(ctx, data1)
repo2.Store(ctx, data2)
```

### 2. Always Handle Errors

```go
// ✅ Good
if err := repo.Store(ctx, data); err != nil {
    log.Errorf("Error: %+v", err)
    return err
}

// ❌ Bad
repo.Store(ctx, data)
```

### 3. Use Prepared Statements

GORM automatically uses prepared statements, but be careful with raw queries:

```go
// ✅ Good
db.Where("id = ?", id).First(&model)

// ❌ Bad (SQL injection risk)
db.Raw(fmt.Sprintf("SELECT * FROM users WHERE id = %s", id))
```

### 4. Use Indexes for Frequently Queried Fields

```go
Email string `gorm:"column:email;uniqueIndex"`
Status string `gorm:"column:status;index"`
```

### 5. Soft Delete for Important Data

```go
DeletedAt gorm.DeletedAt `gorm:"index"`
```

### 6. Use Connection Pooling

Configure connection pool in environment variables:
- `DB_MAX_CONNECTIONS`: Set based on expected load
- `DB_MAX_IDLE_CONNECTIONS`: Keep some connections ready
- `DB_MAX_LIFETIME_CONNECTIONS`: Refresh connections periodically

### 7. Avoid N+1 Queries

```go
// ❌ Bad (N+1 query)
var users []User
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Orders").Find(&user.Orders)
}

// ✅ Good (Eager loading)
var users []User
db.Preload("Orders").Find(&users)
```

### 8. Use Select for Specific Fields

```go
// ✅ Good (only select needed fields)
db.Select("id", "name", "email").Find(&users)

// ❌ Bad (select all fields)
db.Find(&users)
```

## 🔍 Query Optimization

### Use Explain Analyze

```sql
EXPLAIN ANALYZE SELECT * FROM jobs WHERE active = true;
```

### Index Optimization

Create indexes for frequently queried columns:

```sql
CREATE INDEX idx_jobs_active ON jobs(active);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);
```

### Query Patterns

#### Pagination
```go
query.Limit(limit).Offset(offset).Find(&data)
```

#### Filtering
```go
if filter.Active == 1 {
    query = query.Where("active = ?", true)
}
```

#### Sorting
```go
query.Order("created_at DESC")
```

## 📈 Performance Tips

1. **Connection Pooling**: Configure appropriate pool size
2. **Query Optimization**: Use indexes, avoid N+1 queries
3. **Batch Operations**: Use batch insert/update when possible
4. **Selective Loading**: Only load needed fields
5. **Caching**: Use Redis for frequently accessed data

## 🔐 Security

### SQL Injection Prevention

GORM automatically prevents SQL injection with parameterized queries:

```go
// ✅ Safe
db.Where("email = ?", email).First(&user)

// ❌ Unsafe (don't do this)
db.Raw(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))
```

### Database Credentials

- Never commit `.env` files
- Use environment variables for credentials
- Rotate credentials regularly
- Use least privilege principle

## 📚 Related Documentation

- [Development Guidelines](./DEVELOPMENT.md)
- [API Structure](./API_STRUCTURE.md)

---

**Last Updated**: 2024

