# TaskID Implementation Plan

## 1. Goal and Scope
- Establish a consistent internal identifier (`taskID`) for all task-related flows.
- Ensure every event and background job works with the same canonical task record stored in Postgres.
- Maintain compatibility with Notion identifiers (`page_id`, `database_id`) while decoupling domain logic from external IDs.

## 2. Deliverables Overview
- `tasks` module (domain + application + infrastructure) with persistence support **✅ Implemented (Task aggregate, repository, migration)**
- Database migration creating `tasks` table, including mapping to Notion IDs **✅ Added migration `20251001120000_create_tasks.go`**
- Updated `WebhookTriage` that finds or creates tasks based on Notion payloads **✅ Implemented (repo lookup/create, logging)**
- Watermill subscribers (`TaskSynchronizer`, others) that react using `taskID` **⏳ Pending**
- Asynq task DTOs/handlers operating on `taskID` **⏳ Pending**
- Integration tests covering the end-to-end flow from webhook to job execution **🟡 Partially done (domain + repo tests green)**

## 3. Data Model & Persistence
- Created `internal/modules/tasks/domain/task.go` with entity fields (IDs, Notion linkage, sync metadata) **✅**
- Defined repository interface (`FindByNotionIDs`, `FindByID`, `Save`, `Update`, `Upsert`) **✅**
- Implemented GORM model in `infrastructure/postgres` with unique index on `(notion_database_id, notion_page_id)` **✅**
- Added migration `20251001120000_create_tasks.go` **✅**

## 4. Initial Synchronisation Workflow
- Expand or create `ProjectSyncService` to fetch tasks from Notion and populate the `tasks` table **⏳ Pending**
- Upon creation, assign `task.ID` and optional `task.PublicID` **✅ Ready via domain constructor**
- Store mapping between `task.ID` and Notion IDs for subsequent webhook usage **✅ Supported by schema**
- Mark sync status/metadata for incremental updates **⏳ Requires sync implementation**

## 5. WebhookTriage Adjustments
- Inject `TaskRepository` (and possibly `ProjectRepository` for cross-checks) **✅ TasksRepo + ProjectsRepo injected**
- Replace `uuid.New()` with lookup/create logic **✅ resolveTaskID uses repositories**
- Publish events using canonical `taskID` **✅ Events enriched with persisted IDs**
- Emit structured logs indicating whether task was found or created **✅ Added logging**

## 6. Watermill Subscribers
- Implement `TaskSynchronizer` handlers under `tasks/infrastructure/events/` **⏳ Pending**
- Register subscribers inside `cmd/event_worker/main.go` **⏳ Pending**

## 7. Asynq Integration
- Define DTOs for queued jobs in `tasks/application/jobs/` **⏳ Pending**
- Update event subscribers to enqueue jobs using canonical IDs **⏳ Pending**
- Register handlers in `cmd/job_worker/main.go` **⏳ Pending**

## 8. Testing Strategy
- Unit tests for new repository methods and entity invariants **✅ Implemented (`task_test.go`)**
- Integration tests (Ginkgo) covering GORM persistence via testcontainers **✅ (`repository_test.go`)**
- E2E scenario: webhook → triage → job → result **⏳ Pending**

## 9. Configuration & Observability
- Extend `config.Config` if additional settings required **⏳ To evaluate after sync implementation**
- Add structured logging for triage/subscribers/job handlers including `taskID` **🟡 Logging added for triage; subscribers pending**
- Prepare metrics hooks to monitor sync health **⏳ Pending**

## 10. Migration & Rollout Notes
- Migration order: run DB migration before deploying new binaries **🟡 Migration ready; rollout planning pending**
- Backfill existing projects/tasks via one-off sync **⏳ Pending**
- Gradually enable Asynq handlers post-backfill **⏳ Pending**
- Monitor logs for missing `taskID` or duplicate mapping errors **⏳ Pending**

## 11. Open Questions
- Q: Should `PublicID` be introduced immediately or postponed? **Closed**
  A: Implemented

- Q: Do we require soft-delete semantics (`deleted_at`) for tasks? **Answered**
  A: Yes 
- How to reconcile concurrent webhook vs initial sync creating the same task? **Needs design**
- Need for idempotency keys when publishing follow-up events referencing `taskID`? **Still open**
