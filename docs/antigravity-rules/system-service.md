---
trigger: always_on
---

# System Service Rules

Rules untuk **System Service** (Port 9001) - Fiber v2 + sqlx.

## Service Overview

System Service menangani:
- User management & authentication
- Configuration management
- File uploads (Huawei OBS)
- Push notifications
- Day/calendar management

---

## Project Structure

```
system/
├── adapter/          # External adapters (OBS, HTTP client)
├── controller/       # HTTP handlers
├── entity/           # Request/Response DTOs
├── model/            # Database models
├── pkg/              # Shared utilities
│   ├── config/       # Configuration & database
│   ├── middleware/   # Fiber middleware
│   ├── validation/   # Request validation
│   └── server/       # Server setup
├── repository/       # Data access layer
├── service/          # Business logic
└── main.go           # Entry point
```

---

## Coding Conventions

### Controller Pattern

```go
type UserController struct {
    UserService service.UserService
    validator   *validation.Validate
}

func NewUserController(userService service.UserService, validator *validation.Validate) *UserController {
    return &UserController{
        UserService: userService,
        validator:   validator,
    }
}

func (controller *UserController) Route(app *fiber.App) {
    route := app.Group("/v1/users", middleware.JWTProtected())
    route.Get("", controller.List)
    route.Get("/:id", controller.Detail)
    route.Post("", controller.Create)
    route.Patch("/:id", controller.Update)
    route.Delete("/:id", controller.Delete)
}
```

### Handler Method Pattern

```go
func (controller *UserController) Create(c *fiber.Ctx) error {
    var request entity.CreateUserBody
    responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), "")

    // 1. Parse request body
    if err := c.BodyParser(&request); err != nil {
        responsePayload.Setmsg(err.Error())
        return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
    }

    // 2. Inject context data
    request.CustId = c.Locals("cust_id").(string)
    request.CreatedBy = c.Locals("user_id").(int64)

    // 3. Validate request
    if errs := controller.validator.ValidateStruct(request, ""); errs != nil {
        responsePayload.Setmsg(fiber.ErrBadRequest.Message)
        responsePayload.Seterrors(errs)
        return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
    }

    // 4. Call service
    _, err := controller.UserService.Store(request)
    if err != nil {
        responsePayload.Setmsg(err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
    }

    responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
    return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
```

### Service Pattern

```go
type UserService interface {
    Detail(params entity.DetailUserParams) (entity.UserResponse, error)
    List(filter entity.UserQueryFilter) ([]entity.UserResponse, int, int, error)
    Store(request entity.CreateUserBody) (entity.UserResponse, error)
    Update(id int64, request entity.UpdateUserRequest) error
    Delete(id int64) error
}

func NewUserService(repo repository.UserRepository) *userServiceImpl {
    return &userServiceImpl{UserRepository: repo}
}

type userServiceImpl struct {
    UserRepository repository.UserRepository
}
```

### Repository Pattern

```go
type UserRepository interface {
    FindAll(filter entity.UserQueryFilter) ([]model.User, int, int, error)
    FindById(id int64) (model.User, error)
    Store(user model.User) (int64, error)
    Update(id int64, user model.User) error
    Delete(id int64) error
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepositoryImpl{db}
}

type userRepositoryImpl struct {
    *sqlx.DB
}
```

---

## Database Access (sqlx)

### Query Patterns

```go
// Single row query
func (r *userRepositoryImpl) FindById(id int64) (model.User, error) {
    var user model.User
    query := `SELECT * FROM sys.m_user WHERE user_id = $1 AND is_del = false`
    err := r.Get(&user, query, id)
    return user, err
}

// Multiple rows query  
func (r *userRepositoryImpl) FindAll() ([]model.User, error) {
    var users []model.User
    query := `SELECT * FROM sys.m_user WHERE is_del = false`
    err := r.Select(&users, query)
    return users, err
}

// Insert with returning ID
func (r *userRepositoryImpl) Store(user model.User) (int64, error) {
    query := `INSERT INTO sys.m_user (user_name, email) VALUES ($1, $2) RETURNING user_id`
    var id int64
    err := r.QueryRow(query, user.UserName, user.Email).Scan(&id)
    return id, err
}
```

---

## Validation

```go
type CreateUserBody struct {
    UserName  string `json:"user_name" validate:"required,max=100"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    RoleId    int64  `json:"role_id" validate:"required"`
    CustId    string `json:"-"`
    CreatedBy int64  `json:"-"`
}
```

---

## Error Handling

- Return descriptive error messages
- Use `log.Error()` for logging  
- Handle `sql: no rows in result set` as 404 Not Found

---

## Testing Requirements

- Unit tests for service layer with mocked repositories
- Target coverage: >75%
