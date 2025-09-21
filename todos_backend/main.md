# Backend To-Do List (Golang)

This file outlines the technical tasks required to build the backend for the Project Overlay Management for Notion application using Golang.

## Phase 1: Core MVP Setup ✅ COMPLETED

### 1. Project Setup & API Foundation ✅
- [x] Initialize Go module (`go mod init`).
- [x] Set up project structure (cmd, internal, pkg) following DDD architecture.
- [x] Set up Chi web framework for the REST API with proper middleware.
- [x] Implement enhanced health check endpoint (`/health`) with database and Redis status.

### 2. Notion Integration ✅
- [x] Implement OAuth 2.0 flow to authorize users with their Notion accounts.
- [x] Create a comprehensive client/service to interact with the Notion API (`pkg/notion`).
- [x] Implement functions to query databases (`POST /v1/databases/{database_id}/query`).
- [x] Implement functions to update pages (`PATCH /v1/pages/{page_id}`).
- [x] Implement functions to create new pages (`POST /v1/pages`).
- [x] Complete type definitions for all Notion API objects (User, Database, Page, Properties, etc.).

### 3. Database & Caching ✅
- [x] Set up PostgreSQL database connection using GORM with singleton pattern.
- [x] Design and create database schemas for extended data:
    - [x] User accounts and Notion tokens (users table with OAuth integration).
    - [ ] Task dependencies (Phase 2).
    - [ ] Task hierarchy (Phase 2).
- [x] Set up Redis connection for caching API responses and session management.
- [x] Implement Go-based migrations using GORM AutoMigrate (following pet4u-go pattern).
- [x] Create comprehensive config system with singleton pattern and test support.

### 🚀 Available Endpoints:
```
GET  /health                              - Enhanced system health check
POST /api/v1/users                        - Create user
GET  /api/v1/users/{userID}              - Get user by ID  
GET  /api/v1/users/by-email/{email}      - Get user by email
GET  /api/v1/auth/notion/authorize        - Start Notion OAuth flow
GET  /api/v1/auth/notion/callback         - Handle OAuth callback
```

### 🏗️ Architecture Implemented:
- **Domain-Driven Design (DDD)**: Clean separation of domain/application/infrastructure layers
- **Users Module**: Complete with entities, repositories, use cases, and HTTP interfaces
- **Config System**: Centralized configuration with singleton pattern and test support
- **Testing**: Comprehensive test suite with testcontainers for integration tests
- **Error Handling**: Structured HTTP error responses and validation
- **Migrations**: Go-based database migrations using GORM AutoMigrate

## Phase 2: Event-Driven Synchronization (8-10 weeks) - CURRENT FOCUS 🎯

> **Architectural Note:** This phase implements an event-driven architecture using **Watermill** as an internal event bus and **asynq** for background job processing. This approach provides high scalability and loose coupling between modules. See `ARCHITECTURE_EVENTS.md` for a detailed explanation.

> **Progress Update:** Step 1 (Core Infrastructure Setup) - ✅ COMPLETED<br>
> **Step 1.5 (Infrastructure Testing & Validation) - ✅ COMPLETED**<br>
> **Step 2 (Project Module Implementation) - ✅ COMPLETED**<br>
> **Next:** Step 2.5 (Authentication Middleware) - 📝 Ready to start

### 4. Core Infrastructure Setup
- [x] Set up **Watermill** router and configure a publisher/subscriber model
- [x] Set up **asynq** client and server for background job processing
- [x] Define core domain events (e.g., `NotionWebhookReceived`, `TaskPropertiesUpdated`)

### 4.5. Infrastructure Testing & Validation
- [x] Test docker-compose setup with all new services
- [x] Validate worker connectivity to PostgreSQL and Redis
- [x] Test basic event publishing and consumption
- [x] Test basic job queuing and processing

### 5. Project Module & Initial Sync
- [x] Create `projects` domain module (DDD: entity, repository)
- [x] Add database migration for `projects` table (including `notion_webhook_secret`)
- [x] Implement PostgreSQL repository for projects
- [ ] Create `ProjectSyncService` for handling bulk data synchronization from Notion
- [ ] Implement `PerformInitialSync` logic to fetch and store all tasks when a project is first added

### 5.25. Authentication Middleware
- [x] Create JWT middleware for API authentication
- [x] Add user context helpers for request handling
- [x] Update projects router to use authenticated user ID
- [x] Test authentication flow with JWT tokens

### 5.5. First End-to-End Flow (API & Webhook)
- [x] Create use cases and HTTP handlers for projects
- [x] Create DTOs for project operations
- [x] Integrate projects router with main API
- [-] Test end-to-end project creation flow (Skipped for now)

### 6. Notion Webhook & Event-Driven Flow ✅ COMPLETED
- [x] Create `/api/v1/webhooks/notion` endpoint that validates and publishes a `NotionWebhookReceived` event to Watermill
- [ ] Create a `WebhookTriage` Watermill subscriber to process raw events and publish specific domain events (e.g., `TaskPropertiesUpdated`)
- [ ] Create a `TaskSynchronizer` Watermill subscriber to update the local database based on domain events
- [ ] Implement robust `X-Notion-Signature` validation for security

### 7. Core Feature Logic - Tasks & Background Jobs
- [ ] Create `tasks` domain module with dependency and hierarchy support
- [ ] Create a `CriticalPathService` Watermill subscriber that enqueues a job in **asynq** when a task's date changes
- [ ] Create an **asynq** worker for heavy-lifting tasks like critical path calculation
- [ ] Domain entities (Task, TaskDependency, TaskHierarchy)
- [ ] Business logic for task dependencies and cascade effects
- [ ] Database migrations for `tasks`, `dependencies`, and `hierarchy` tables

### 8. API Endpoints & Real-time Frontend Updates
- [ ] **Handle Eventual Consistency in API/UI:** Define a clear contract for notifying the frontend about ongoing background processes (e.g., a "syncing" status in API responses or via SSE) so it can display appropriate indicators until a final confirmation event is received.
- [ ] **Projects API**:
    - [ ] `GET /api/v1/projects` - List user projects
    - [ ] `POST /api/v1/projects` - Create/sync project from Notion (triggers initial sync)
    - [ ] `POST /api/v1/projects/{id}/resync` - Manually trigger a full re-synchronization
    - [ ] `DELETE /api/v1/projects/{id}` - Delete a project
- [ ] **Tasks API**:
    - [ ] `GET /api/v1/projects/{id}/tasks` - Get project tasks with dependencies
    - [ ] `PUT /api/v1/tasks/{id}/dependencies` - Update task dependencies
- [ ] **SSE Hub**:
    - [ ] Create SSE hub in `pkg/sse` for managing connections
    - [ ] Create a `SSENotifier` Watermill subscriber that listens for events (e.g., `CriticalPathCalculated`) and pushes updates to clients
    - [ ] Implement `/api/v1/events` endpoint with user authentication

## Phase 3: PM Features (4-6 weeks)

### 9. Additional Features & Scalability
- [ ] Implement API endpoints for time tracking.
- [ ] Implement API endpoints for multi-project dashboard data.
- [ ] Develop logic for team workload visualization.
- [ ] Create a cron job system (or scheduled goroutines) for generating periodic progress reports.
- [ ] Implement robust error handling, request queuing, and exponential backoff to respect Notion's API rate limits.
