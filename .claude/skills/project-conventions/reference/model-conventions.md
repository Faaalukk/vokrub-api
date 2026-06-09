## Model Conventions

### Struct layout

```go
type Word struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    // domain fields below
    Word    string `json:"word"`
    Meaning string `json:"meaning"`
    UserID  uint   `json:"user_id"`
}
```

- Always include the four base fields (ID, CreatedAt, UpdatedAt, DeletedAt) manually — gives full control over JSON tags.
- `DeletedAt` always uses `json:"-"` — never expose soft-delete timestamp.
- Unique columns: `gorm:"uniqueIndex"`.
- Foreign keys: store the ID field (`UserID uint`) + add `gorm:"index"` tag.
- Enums stored as strings with a comment on valid values.

### AutoMigrate

Add every new model to `database.Connect()` in `database/db.go`:

```go
DB.AutoMigrate(&models.Word{}, &models.Customer{}, &models.User{})
```

AutoMigrate only adds columns and indexes — it never drops. Safe to run on every startup.

### Input structs

Never parse request body directly into a model struct (exposes all fields). Use a dedicated input type:

```go
type CreateWordInput struct {
    Word    string `json:"word"`
    Meaning string `json:"meaning"`
}
```
