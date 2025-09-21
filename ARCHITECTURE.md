# Backend Architecture (Golang)

This document outlines the proposed backend architecture for the Project Overlay Management for Notion application, built with Golang. The architecture is heavily inspired by Domain-Driven Design (DDD) principles, aiming for a modular, scalable, and maintainable codebase.

## Guiding Principles

*   **Domain-Driven Design (DDD)**: The core logic is organized into distinct domains (or modules) that represent specific business capabilities. This keeps concerns separated and the codebase easier to understand.
*   **Clean Architecture**: The architecture follows the principles of Clean Architecture, ensuring that dependencies flow inwards. The core domain logic does not depend on external concerns like databases, APIs, or frameworks.
*   **Modularity**: The system is composed of loosely coupled modules that can be developed, tested, and deployed independently.

## Project Structure

The proposed project structure will follow the standard Go layout, similar to the `pet4u-go` project.

```
/
├── cmd/
│   ├── api/                # Main application entry point (HTTP server)
│   │   └── main.go
│   └── migrate/            # Database migration tool
│       └── main.go
├── internal/
│   ├── cache/              # Redis client and caching logic
│   ├── config/             # Configuration loading and management
│   ├── database/           # PostgreSQL connection and helpers
│   ├── modules/            # Core business domains (DDD modules)
│   │   ├── projects/       # Domain for managing projects
│   │   ├── tasks/          # Domain for managing tasks and dependencies
│   │   ├── users/          # Domain for users and authentication
│   │   └── shared/         # Shared kernel logic (e.g., value objects)
│   ├── pkg/                # Shared utility packages
│   │   ├── notion/         # Client for Notion API
│   │   └── sse/            # Server-Sent Events (SSE) hub
│   └── server/             # HTTP server setup (routing, middleware)
├── migrations/             # SQL migration files
├── go.mod
├── go.sum
└── Dockerfile
```

## Domain Module Structure

Each domain module (e.g., `internal/modules/tasks/`) will have a consistent internal structure:

```
tasks/
├── application/            # Application services (use cases)
│   ├── task_service.go
│   └── dto/                # Data Transfer Objects
├── domain/                 # Core domain logic
│   ├── task.go             # Aggregate root (entity)
│   ├── task_repository.go  # Repository interface
│   └── value_objects/      # Domain-specific value objects
└── infrastructure/         # Implementation of external concerns
    ├── task_postgres_repository.go # GORM/sqlx implementation of the repository
    └── task_handler.go     # HTTP handlers for this module
```

*   **Domain**: Contains the core business logic, entities, value objects, and repository interfaces. This layer is pure and has no external dependencies.
*   **Application**: Orchestrates the domain logic to perform specific use cases. It depends on the domain layer.
*   **Infrastructure**: Provides the implementation details for interfaces defined in the domain layer, such as database repositories and HTTP handlers. It depends on the application and domain layers.

## Testing Strategy

*   **Unit Tests**: Each component in the `domain` and `application` layers will be thoroughly unit-tested in isolation. We will use mocks for dependencies (like repositories).
*   **Integration Tests**: We will write integration tests to verify the interaction between different layers, especially the infrastructure layer's connection to the database and external APIs. These will run against a test database.
*   **End-to-End (E2E) Tests**: A small number of E2E tests will verify entire user flows from the API endpoint down to the database.

### Testing Libraries and Approach

To implement our testing strategy, we will adopt the following libraries and conventions, similar to the `pet4u-go` project:

*   **Ginkgo**: For structuring our tests in a Behavior-Driven Development (BDD) style. This makes tests more readable and expressive by using `Describe` and `It` blocks.
*   **Gomega**: As the primary assertion library, paired with Ginkgo. Its rich set of matchers allows for clear and fluent assertions (`Expect(err).ToNot(HaveOccurred())`).
*   **Testify**: We will use `testify/require` for essential assertions and `testify/mock` for creating mock objects for our interfaces (like repositories) in unit tests. This allows us to test components in isolation.
*   **Test Suites**: We will group related tests into suites for better organization and to share setup/teardown logic (e.g., `suite_test.go`).

## Reusable Components (`/pkg`)

The `/internal/pkg` directory will house packages that are shared across different modules but are not core domain logic.

*   `notion`: A dedicated client for interacting with the Notion API, handling authentication, requests, and error responses.
*   `sse`: A central hub for managing Server-Sent Events (SSE) connections and broadcasting events to clients. This will be used for real-time updates.

This architecture will provide a solid foundation for building a robust and scalable backend for the Project Overlay Management for Notion application.
