## Directory Structure

```
vokrub-api/
├── main.go                  # entry: env, db, fiber, middleware, routes, listen
├── go.mod / go.sum
├── .env
├── database/
│   └── db.go                # Connect() + DB package var
├── models/
│   ├── user.go              # User struct
│   └── customer.go          # Customer struct
├── handlers/
│   ├── auth.go              # Login, Register, Me
│   ├── user.go              # GetUsers, GetUser, CreateUser, DeleteUser
│   └── customer.go          # GetCustomers, GetCustomer, CreateCustomer, DeleteCustomer
├── middleware/
│   └── auth.go              # Protected middleware
└── routes/
    └── routes.go            # Setup(app) — all route registration
```

### Adding a new domain (e.g. `word`)

1. `models/word.go` — define `Word` struct, add to AutoMigrate in `database/db.go`
2. `handlers/word.go` — handler functions
3. `routes/routes.go` — add group + routes inside `Setup(app)`

One file per domain per layer. No handler logic in `routes.go`, no route registration in `main.go`.
