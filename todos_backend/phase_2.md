# Phase 2: Detailed Plan - Event-Driven Synchronization

> **Architectural Note:** This document provides a detailed, step-by-step implementation plan for Phase 2. The goal is to implement an event-driven architecture using **Watermill** and **asynq**. See `ARCHITECTURE_EVENTS.md` for the high-level concept.

---

### Step 1: Core Infrastructure Setup
*   **Goal:** Establish the foundational infrastructure for handling events and background jobs.
*   **Status:** ✅ COMPLETED

| Task | Implementation Details | Status |
| :--- | :--- | :--- |
| **Create `eventbus` and `taskqueue` packages** | **`src/internal/pkg/eventbus/watermill.go`**: <br> - `NewPublisher()` -> `message.Publisher` <br> - `NewRouter(logger watermill.LoggerAdapter)` -> `*message.Router` <br> - Start with `GoChannel` implementation for local dev. <br> **`src/internal/pkg/taskqueue/asynq.go`**: <br> - `NewClient()` -> `*asynq.Client` <br> - `NewServer(redisOpt asynq.RedisConnOpt)` -> `*asynq.Server` | `[x]` |
| **Create worker entry points** | **`src/cmd/event_worker/main.go`**: <br> - Initializes config, logger. <br> - Creates a Watermill router using `eventbus.NewRouter()`. <br> - Ready to register event subscriber handlers. <br> - Runs the router: `router.Run(ctx)`. <br> **`src/cmd/job_worker/main.go`**: <br> - Initializes config, logger, redis. <br> - Creates an `asynq.ServeMux` to act as the handler. <br> - Ready to register job handlers to the mux. <br> - Creates and runs the Asynq server using `taskqueue.NewServer()`. | `[x]` |
| **Update `src/docker-compose.yml`** | Add two new services by duplicating the `api` service configuration. The key change will be the `command` property: <br> ```yaml <br> event_worker: <br>   build: <br>     context: . <br>     dockerfile: Dockerfile.dev <br>   working_dir: /app <br>   command: sh -c "go run ./cmd/event_worker" <br>   env_file: [ .env ] <br>   volumes: [ ".:/app", ... ] <br>   depends_on: <br>     psql_bp: { condition: service_healthy } <br>     redis: { condition: service_healthy } <br> <br> job_worker: <br>   # ... similar config ... <br>   command: sh -c "go run ./cmd/job_worker" <br> ``` <br> *Note: Using `go run` directly for development simplicity.* | `[x]` |
| **Define core domain events** | Create `src/internal/modules/shared/domain/events/events.go`. <br> Add the first event struct and topic constant: <br> ```go <br> package events <br> <br> const NotionWebhookReceivedTopic = "notion.webhook.received" <br> <br> type NotionWebhookReceived struct { <br>   Payload []byte <br> } <br> ``` | `[x]` |

---

### Step 1.5: Infrastructure Testing & Validation
*   **Goal:** Test and validate that all infrastructure components work together properly.
*   **Status:** ✅ COMPLETED

| Task | Implementation Details | Status |
| :--- | :--- | :--- |
| **Test docker-compose setup** | Run `docker-compose up` and verify that all services start correctly: <br> - `api` service <br> - `event_worker` service <br> - `job_worker` service <br> - PostgreSQL and Redis health checks | `[x]` |
| **Validate worker connectivity** | Ensure workers can connect to: <br> - PostgreSQL database <br> - Redis instance <br> - Watermill event bus (GoChannel) | `[x]` |
| **Test basic event flow** | Create a simple test to verify event publishing: <br> - Publish `NotionWebhookReceived` event from a test script <br> - Verify it's received by event_worker (add temporary logging) | `[x]` |
| **Test basic job queuing** | Create a simple test to verify job queuing: <br> - Enqueue a basic job using Asynq client <br> - Verify it's processed by job_worker (add temporary logging) | `[x]` |

---

### Step 2: Project Module Implementation
*   **Goal:** Create the business logic and data structures for managing projects.
*   **Status:** ✅ COMPLETED

| Task | Implementation Details | Status |
| :--- | :--- | :--- |
| **Create `projects` module structure** | Create `src/internal/modules/projects/` with subdirectories: <br> - `application` (use cases) <br> - `domain` (entities, value objects, repo interface) <br> - `infrastructure/postgres` (GORM implementation) <br> - `infrastructure/http` (handlers) <br> - `infrastructure/events` (subscribers) | `[x]` |
| **Database Migration** | Create `src/migrations/20250921123742_create_projects.go`. <br> Inside, define a GORM model: `type ProjectRecord struct { gorm.Model; UserID uuid.UUID; NotionDatabaseID string; NotionWebhookSecret string; }`. <br> Use `m.AutoMigrate(&ProjectRecord{})` to create the table. | `[x]` |
| **Define Entity and Repository** | **`src/internal/modules/projects/domain/project.go`**: <br> - Define the `Project` entity with validation and Clock interface. <br> **`src/internal/modules/projects/domain/repository.go`**: <br> - `type Repository interface { Save(ctx context.Context, project *Project) error; FindByID(ctx context.Context, id uuid.UUID) (*Project, error); ... }` | `[x]` |
| **Implement PostgreSQL Repository** | **`src/internal/modules/projects/infrastructure/postgres/repository.go`**: <br> - `type ProjectRepository struct { db *gorm.DB }` <br> - Implement the `domain.Repository` interface with full CRUD operations. | `[x]` |

---

### Step 2.5: Authentication Middleware
*   **Goal:** Implement JWT-based authentication middleware for secure API access.
*   **Status:** 📝 To Do

| Task | Implementation Details | Status |
| :--- | :--- | :--- |
| **Create JWT middleware** | **`src/internal/pkg/middleware/auth.go`**: <br> - JWT token validation middleware <br> - Extract user ID from token claims <br> - Handle authentication errors | `[ ]` |
| **Update projects router** | Modify `src/internal/modules/projects/interfaces/http/router.go` to: <br> - Remove placeholder user ID <br> - Use authenticated user ID from middleware context | `[ ]` |
| **Add user context helper** | **`src/internal/pkg/middleware/context.go`**: <br> - Helper functions to get/set user ID in request context <br> - Type-safe context keys | `[ ]` |
| **Test authentication flow** | Verify JWT tokens are properly validated and user context is available | `[ ]` |

---

### Step 3: First End-to-End Flow (API & Webhook)
*   **Goal:** Connect all components into a single, working, simplified flow.
*   **Status:** 📝 To Do

| Task | Implementation Details | Status |
| :--- | :--- | :--- |
| **Create use cases** | **`src/internal/modules/projects/application/create_project.go`**: <br> - `CreateProjectUseCase` with validation and repository interaction | `[x]` |
| **Create HTTP handlers** | **`src/internal/modules/projects/interfaces/http/router.go`**: <br> - POST `/` for creating projects <br> - GET `/` for listing user projects | `[x]` |
| **Create DTOs** | **`src/internal/modules/projects/interfaces/http/dto.go`**: <br> - `CreateProjectRequestDTO`, `ProjectResponseDTO`, `ProjectsListResponseDTO` <br> - Functions for mapping between domain and DTO | `[x]` |
| **Integrate with main API** | Update `src/internal/server/routes.go` to mount projects router at `/api/v1/projects` | `[ ]` |
| **Test end-to-end flow** | Create a project via API and verify it's stored in database | `[ ]` |
