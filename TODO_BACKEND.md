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

## Phase 2: Advanced Features (6-8 weeks) - CURRENT FOCUS 🎯

### 4. Real-time Synchronization
- [ ] Create WebSocket hub in `pkg/websocket` for managing connections and rooms
- [ ] Implement WebSocket endpoints (`/api/v1/ws`) with user authentication
- [ ] Create projects domain module following DDD pattern
- [ ] Create tasks domain module with dependency and hierarchy support
- [ ] Implement background service (using goroutines) to poll Notion API for changes
- [ ] Add change detection logic and broadcast via WebSockets to subscribed clients

### 5. Core Feature Logic - Tasks & Projects
- [ ] **Projects Module**: 
    - [ ] Domain entities (Project, ProjectRepository)
    - [ ] Application services (CreateProject, GetProject, SyncWithNotion) 
    - [ ] Infrastructure (PostgreSQL repository, HTTP handlers)
    - [ ] Database migrations for projects table
- [ ] **Tasks Module**:
    - [ ] Domain entities (Task, TaskDependency, TaskHierarchy)
    - [ ] Business logic for task dependencies and cascade effects
    - [ ] Critical path calculation algorithms
    - [ ] Parent-child relationship management
    - [ ] Database migrations for tasks, dependencies, and hierarchy tables
- [ ] **API Endpoints**:
    - [ ] `GET /api/v1/projects` - List user projects
    - [ ] `POST /api/v1/projects` - Create/sync project from Notion
    - [ ] `GET /api/v1/projects/{id}/tasks` - Get project tasks with dependencies
    - [ ] `PUT /api/v1/tasks/{id}/dependencies` - Update task dependencies
    - [ ] `GET /api/v1/tasks/{id}/critical-path` - Calculate critical path

## Phase 3: PM Features (4-6 weeks)

### 6. Additional Features & Scalability
- [ ] Implement API endpoints for time tracking.
- [ ] Implement API endpoints for multi-project dashboard data.
- [ ] Develop logic for team workload visualization.
- [ ] Create a cron job system (or scheduled goroutines) for generating periodic progress reports.
- [ ] Implement robust error handling, request queuing, and exponential backoff to respect Notion's API rate limits.
