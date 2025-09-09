# Backend To-Do List (Golang)

This file outlines the technical tasks required to build the backend for the Project Overlay Management for Notion application using Golang.

## Phase 1: Core MVP Setup (4-6 weeks)

### 1. Project Setup & API Foundation
- [ ] Initialize Go module (`go mod init`).
- [ ] Set up project structure (e.g., cmd, internal, pkg).
- [ ] Choose and set up a web framework (e.g., Gin, Echo, Fiber) for the REST API.
- [ ] Implement a basic health check endpoint (e.g., `/ping`).

### 2. Notion Integration
- [ ] Implement OAuth 2.0 flow to authorize users with their Notion accounts.
- [ ] Create a client/service to interact with the Notion API (`@notionhq/client` equivalent for Go).
- [ ] Implement functions to query databases (`POST /v1/databases/{database_id}/query`).
- [ ] Implement functions to update pages (`PATCH /v1/pages/{page_id}`).
- [ ] Implement functions to create new pages (`POST /v1/pages`).

### 3. Database & Caching
- [ ] Set up PostgreSQL database connection (e.g., using `gorm` or `sqlx`).
- [ ] Design and create database schemas for extended data:
    - User accounts and Notion tokens.
    - Task dependencies.
    - Task hierarchy.
- [ ] Set up Redis connection for caching API responses and session management.

## Phase 2: Advanced Features (6-8 weeks)

### 4. Real-time Synchronization
- [ ] Implement WebSocket server to push updates to the frontend.
- [ ] Create a background service (using goroutines) to poll the Notion API for changes periodically (e.g., every 10 seconds).
- [ ] Implement logic to detect changes and broadcast them to relevant clients via WebSockets.

### 5. Core Feature Logic
- [ ] Develop API endpoints for the frontend to fetch project and task data.
- [ ] Implement business logic for task dependencies:
    - Store dependency relationships in PostgreSQL.
    - Calculate cascade effects when a task's date changes.
    - Implement logic to identify the critical path.
- [ ] Implement business logic for hierarchical tasks, managing parent-child relationships.

## Phase 3: PM Features (4-6 weeks)

### 6. Additional Features & Scalability
- [ ] Implement API endpoints for time tracking.
- [ ] Implement API endpoints for multi-project dashboard data.
- [ ] Develop logic for team workload visualization.
- [ ] Create a cron job system (or scheduled goroutines) for generating periodic progress reports.
- [ ] Implement robust error handling, request queuing, and exponential backoff to respect Notion's API rate limits.
