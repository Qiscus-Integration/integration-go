# Integration Boilerplate (Golang)

### Overview

The Qiscus Integration team's Go project boilerplate is designed to speed up application development and ensure consistency across all integration projects. By providing a standardized structure and set of guidelines, this boilerplate makes it easier for developers to understand, review, and maintain projects. By using this boilerplate, developers can ensure that all integration projects follow a consistent structure and coding style, making it easier to onboard new team members and maintain projects in the long term.

This boilerplate is used to standardize the directory structure for projects of medium to large complexity or with the potential for it. However, **for other cases that only handle one or a few processes, it is not necessary to implement this boilerplate** in order to avoid over-abstraction. For example, you can use a single file like main.go or a flat architecture where all components such as handlers, services, and repositories are placed in a single directory for simplicity and ease of navigation. Here’s a repository that implements this kind of flat structure: [Repo](https://bitbucket.org/qiscus/panin/src/main/).

One significant change from the [v1](https://bitbucket.org/qiscus/integration-go/src/v1/) is moving away from grouping code by function and, instead, organizing it by module. This approach offers several advantages:

- **Single Responsibility Principle**: One of the SOLID principles of object-oriented design, states that a class or module should have only one reason to change.
- **Reusability**: Modules become more reusable across different clients, promoting code sharing and reducing duplication.
- **Loose Coupling, High Cohesion**: Two words that describe how easy or difficult it is to change a piece of software. Grouping by module enforces loose coupling between different parts of the code while promoting high cohesion within each module.
- **Faster Contribution**: Developers can contribute to specific modules without causing **collateral damage** in unrelated areas, speeding up the development process.
- **Ease of Understanding**: The codebase becomes more accessible and understandable as it's organized around modules and use cases. A use case repository clarifies what each module does.

### Create New Module/API

This section guides you through creating a new API module following the established patterns in this codebase.

#### Architecture Overview

This project follows "Clean" Architecture principles with these layers:

- **Handler** (`internal/{module}/handler.go`) - HTTP layer that handles requests/responses
- **Service** (`internal/{module}/service.go`) - Business logic layer
- **Repository** (`internal/{module}/repo.go`) - Data access layer
- **Entity** (`internal/entity/{module}.go`) - Domain models

#### Step-by-Step Guide

**1. Create the Entity**

Create your domain model in `internal/entity/{module}.go`:

```go
package entity

import "time"

type YourModule struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name" gorm:"index"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**2. Create the Repository**

Create `internal/{module}/repo.go`:

```go
package yourmodule

import (
    "context"
    "integration-go/internal/entity"
    "gorm.io/gorm"
)

type repo struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *repo {
    return &repo{db: db}
}

func (r *repo) Save(ctx context.Context, item *entity.YourModule) error {
    return r.db.WithContext(ctx).Save(item).Error
}

func (r *repo) FindByID(ctx context.Context, id int64) (*entity.YourModule, error) {
    var item entity.YourModule
    err := r.db.WithContext(ctx).First(&item, id).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}
```

**3. Create the Service**

Create `internal/{module}/service.go`:

```go
package yourmodule

import (
    "context"
    "errors"
    "fmt"
    "integration-go/internal/entity"
    "gorm.io/gorm"
)

//go:generate mockery --with-expecter --case snake --name Repository
type Repository interface {
    Save(ctx context.Context, item *entity.YourModule) error
    FindByID(ctx context.Context, id int64) (*entity.YourModule, error)
}

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*entity.YourModule, error) {
    item, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("item not found")
        }
        return nil, fmt.Errorf("failed to find item: %w", err)
    }
    return item, nil
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) error {
    item := &entity.YourModule{
        Name: req.Name,
    }

    if err := s.repo.Save(ctx, item); err != nil {
        return fmt.Errorf("failed to save item: %w", err)
    }

    return nil
}

type CreateRequest struct {
    Name string `json:"name" validate:"required"`
}
```

**4. Create the Handler**

Create `internal/{module}/handler.go`:

```go
package yourmodule

import (
    "encoding/json"
    "integration-go/internal/api/resp"
    "net/http"
    "strconv"

    "github.com/rs/zerolog/log"
)

type httpHandler struct {
    svc *Service
}

func NewHttpHandler(svc *Service) *httpHandler {
    return &httpHandler{svc: svc}
}

func (h *httpHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        resp.WriteJSONFromError(w, err)
        return
    }

    item, err := h.svc.GetByID(ctx, int64(id))
    if err != nil {
        log.Ctx(ctx).Error().Msgf("failed to get item: %s", err.Error())
        resp.WriteJSONFromError(w, err)
        return
    }

    resp.WriteJSON(w, http.StatusOK, item)
}

func (h *httpHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var req CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        resp.WriteJSONFromError(w, err)
        return
    }

    if err := h.svc.Create(ctx, &req); err != nil {
        log.Ctx(ctx).Error().Msgf("failed to create item: %s", err.Error())
        resp.WriteJSONFromError(w, err)
        return
    }

    resp.WriteJSON(w, http.StatusCreated, "created")
}
```

**5. Register Routes in Server**

Add your module to `internal/api/server.go` in the `NewServer()` function:

```go
// YourModule
yourModuleRepo := yourmodule.NewRepository(db)
yourModuleSvc := yourmodule.NewService(yourModuleRepo)
yourModuleHandler := yourmodule.NewHttpHandler(yourModuleSvc)

// Add routes
r.Handle("GET /api/v1/yourmodule/{id}", authMidd.StaticToken(http.HandlerFunc(yourModuleHandler.GetByID)))
r.Handle("POST /api/v1/yourmodule", authMidd.StaticToken(http.HandlerFunc(yourModuleHandler.Create)))
```

**6. Add Database Migration (if needed)**

If your entity needs database tables, add migration to `internal/postgres/migrate.go`:

```go
func Migrate(db *gorm.DB) {
    db.AutoMigrate(
        &entity.Room{},
        &entity.YourModule{}, // Add your entity here
    )
}
```

#### Key Patterns to Follow

- **Error Handling**: Use `resp.WriteJSONFromError(w, err)` for consistent error responses
- **Logging**: Use `log.Ctx(ctx).Error().Msgf()` for contextual logging
- **Validation**: Use struct tags with `validate` for request validation
- **Database**: Always use `WithContext(ctx)` for database operations
- **Mocking**: Add `//go:generate mockery` comments for interfaces that need mocks

#### Available Response Utilities

- `resp.WriteJSON(w, statusCode, data)` - Standard JSON response
- `resp.WriteJSONFromError(w, err)` - Error response with proper status codes
- `resp.WriteJSONWithPaginate(w, statusCode, data, total, page, limit)` - Paginated response

### Sample Use Case

You can find the documentation for this application [here](/docs/README.md)

Create tagging and storing room data in the database when a user initiates a chat on Qiscus Omnichannel by utilizing a new session webhook, besides that there is a cronjob that does an auto resolved room when the on going room has reached 10 minutes. And also create an API to get all rooms that are stored in the database with an API key authentication.

### Environment Variables

Set up the environment file by copying .env.example:

Mac/Linux:
`cp .env.example .env`

Windows:
`copy .env.example .env`

Alternatively, you can create a copy of `.env.example` and rename it to `.env` in your project's root directory

### Run Locally

To run the project locally, follow these steps:

- Clone this repository
- Navigate to the directory
- Format code and tidy modfile: `make tidy`
- Run test: `make test`, make sure that all tests are passing
- Run the server: `make run bin=api`, or run the application with reloading on file changes with: `make run/live bin=api`. You can also apply this to the cron application by changing the parameter to `bin=cron`
- The backend server will be accessible at `http://localhost:8080`
- You can find another usefull commands in `Makefile`

### Generate Mock from Interface

- Install [Mockery](https://github.com/vektra/mockery):

  ```bash
  go install github.com/vektra/mockery/v2@v2.53.4
  ```

  > **Note:** In this repo, we use Mockery v2.

- Add the following code in the interface code file: `//go:generate mockery --case snake --name XXXX`
- Run go generate using `make generate`

### Handle HTTP Client Exceptions

`client.Error` complies with Go standard error. which support Error, Unwrap, Is, As

```go
// Sample using errors.As
err := s.omni.AssignAgent(ctx, agentID, roomID)
if err != nil {
    var cerr *client.Error
    if errors.As(err, &cerr) {
        fmt.Println(cerr.Message)     		// General error message
        fmt.Println(cerr.StatusCode)  		// HTTP status code e.g: 400, 401 etc.
        fmt.Println(cerr.RawError)    		// Raw Go error object
        fmt.Println(cerr.RawAPIResponse)  // Raw API response body in byte
    }
}

```
