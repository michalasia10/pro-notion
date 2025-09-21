# Phase 2: Detailed Plan - Event-Driven Synchronization

> **Architectural Note:** This document provides a detailed, step-by-step implementation plan for Phase 2. The goal is to implement an event-driven architecture using **Watermill** and **asynq**. See `ARCHITECTURE_EVENTS.md` for the high-level concept.

---

### Step 1: Core Infrastructure Setup
*   **Goal:** Establish the foundational infrastructure for handling events and background jobs.
*   **Status:** 📝 To Do

| Task | Implementation Details | Status |
| :--- | :--- | :--- |
| **Create `eventbus` and `taskqueue` packages** | **`src/internal/pkg/eventbus/watermill.go`**: <br> - `NewPublisher()` -> `message.Publisher` <br> - `NewRouter(logger watermill.LoggerAdapter)` -> `*message.Router` <br> - Start with `GoChannel` implementation for local dev. <br> **`src/internal/pkg/taskqueue/asynq.go`**: <br> - `NewClient()` -> `*asynq.Client` <br> - `NewServer(redisOpt asynq.RedisConnOpt, handler asynq.Handler)` -> `*asynq.Server` | `[ ]` |
| **Create worker entry points** | **`src/cmd/event_worker/main.go`**: <br> - Initializes config, logger, db connection. <br> - Creates a Watermill router using `eventbus.NewRouter()`. <br> - Registers all event subscriber handlers to the router. <br> - Runs the router: `router.Run(ctx)`. <br> **`src/cmd/job_worker/main.go`**: <br> - Initializes config, logger, db, redis. <br> - Creates an `asynq.ServeMux` to act as the handler. <br> - Registers all job handlers to the mux. <br> - Creates and runs the Asynq server using `taskqueue.NewServer()`. | `[ ]` |
| **Update `src/docker-compose.yml`** | Add two new services by duplicating the `api` service configuration. The key change will be the `command` property: <br> ```yaml <br> event_worker: <br>   build: <br>     context: . <br>     dockerfile: Dockerfile.dev <br>   working_dir: /app <br>   command: sh -c "air -c .air.toml event_worker" # Or similar for live reload <br>   env_file: [ .env ] <br>   volumes: [ ".:/app", ... ] <br>   depends_on: <br>     psql_bp: { condition: service_healthy } <br>     redis: { condition: service_healthy } <br> <br> job_worker: <br>   # ... similar config ... <br>   command: sh -c "air -c .air.toml job_worker" <br> ``` <br> *Note: This assumes we will adapt `.air.toml` to run different main packages.* | `[ ]` |
| **Define core domain events** | Create `src/internal/modules/shared/domain/events/events.go`. <br> Add the first event struct and topic constant: <br> ```go <br> package events <br> <br> const NotionWebhookReceivedTopic = "notion.webhook.received" <br> <br> type NotionWebhookReceived struct { <br>   Payload []byte <br> } <br> ``` | `[ ]` |

---

### Step 2: Project Module Implementation
*   **Goal:** Create the business logic and data structures for managing projects.
*   **Status:** 📝 To Do

| Task | Implementation Details | Status |
| :--- | :--- | :--- |
| **Create `projects` module structure** | Create `src/internal/modules/projects/` with subdirectories: <br> - `application` (use cases) <br> - `domain` (entities, value objects, repo interface) <br> - `infrastructure/postgres` (GORM implementation) <br> - `infrastructure/http` (handlers) <br> - `infrastructure/events` (subscribers) | `[ ]` |
| **Database Migration** | Create `src/migrations/00002_create_projects.go`. <br> Inside, define a GORM model: `type Project struct { gorm.Model; UserID uuid.UUID; NotionDatabaseID string; NotionWebhookSecret string; }`. <br> Use `db.AutoMigrate(&Project{})` to create the table. | `[ ]` |
| **Define Entity and Repository** | **`src/internal/modules/projects/domain/project.go`**: <br> - Define the `Project` entity, which will be the "pure" domain model without GORM tags. <br> **`src/internal/modules/projects/domain/repository.go`**: <br> - `type Repository interface { Save(ctx context.Context, project *Project) error; FindByID(ctx context.Context, id uuid.UUID) (*Project, error); ... }` | `[ ]` |
| **Implement PostgreSQL Repository** | **`src/internal/modules/projects/infrastructure/postgres/repository.go`**: <br> - `type GormRepository struct { db *gorm.DB }` <br> - Implement the `domain.Repository` interface. This layer will handle mapping between the domain entity and the GORM model. | `[ ]` |
