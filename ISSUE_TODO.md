## Architecture and Operational TODOs

- [x] Swap Watermill transport from in-memory gochannel to a shared backend (Redis/Kafka/Postgres) so `api` publishers and `event_worker` subscribers share the same bus; wire via config and inject publishers/subscribers from `cmd` instead of constructing inside routers. (Done: Redis transport implemented and default; publisher/subscriber wired in `cmd`/server, routers no longer construct them.)
- [x] Provide a real transaction manager (GORM-backed) and inject it into use cases; stop using `NoopTransactionManager` in HTTP routers so “transactional” use cases actually run in transactions. (Done: GORM-backed manager, repos use transactional session, routers receive txMgr from `cmd/api`.)
- [x] Move dependency wiring (repos, clocks, ID generators, event publishers) out of interface routers and into `cmd/api`; keep routers thin to improve testability and layering. (Done for users/projects/auth routers; wiring now in `cmd/api`/server deps.)
- [x] Separate migration execution from API startup (or guard with a flag) to avoid blocking boot and competing migrations; rely on the `migrate` command. (Done: `RUN_MIGRATIONS=true` gates migrations in `cmd/api`.)
- [x] Add graceful shutdown and resource cleanup: close DB/Redis clients; wire HTTP server shutdown hooks; fix `job_worker` to register handlers, start server in a goroutine, and handle signals before `Start`. (Partial: `job_worker` has handler, goroutine/signal flow, deferred shutdown, and closes Redis; `event_worker` validates/closes Redis and closes DB; API defers pubsub and closes Redis/DB; remaining cleanup/refinement may be needed.)
- [x] Replace gochannel health checks that `log.Fatalf` on DB ping failures with non-fatal health reporting so `/health` doesn’t crash the process when the DB is down. (Done in `internal/database/database.go`.)

## Domain and API Correctness

- [x] Introduce explicit “already exists” errors: projects (duplicate Notion database) and users (duplicate email) currently misuse “not found/invalid” errors and return the wrong HTTP status; propagate to HTTP 409s.
- [ ] In webhook flow, avoid double signature validation: validate in middleware, then pass a verified payload to the service; keep config access in one place.
- [x] Harden JWT middleware: verify issuer/audience and reject tokens with unexpected alg/claims; ensure Notion OAuth validates `state` to prevent CSRF. (Done: JWT enforces HS256 + issuer/audience; OAuth generates state, stores in HttpOnly cookie, validates and clears it on callback.)
- [ ] Complete task queue wiring: register Asynq handlers for domain tasks, add tests around webhook triage → task creation/upsert → queue dispatch.
- [ ] Revisit webhook triage handling for non-database parents (skip vs. resolve); ensure rollup/project/task hierarchies can be fully reconstructed.

## Testing and Observability

- [ ] Add integration tests for webhook triage covering page created/updated/deleted and database updated paths, including task upsert and event publishing.
- [x] Add smoke tests for `/health` that assert non-fatal behavior when DB/Redis are down. (Done in `internal/server/health_test.go`.)
- [ ] Add startup/logging checks for missing critical config (Notion client/secret, webhook secret, JWT secret) and fail fast with clear errors.
