# Learnings

<!-- INDEX:START -->
| Date | Learning |
|------|----------|
| 2026-01-15 | Database connections need explicit timeouts |
| 2026-01-10 | Environment variables override config files |
| 2026-01-05 | Rate limiter must be per-user, not global |
<!-- INDEX:END -->

---

## [2026-01-15-143022] Database connections need explicit timeouts

**Context**: Production outage caused by database connection pool exhaustion.
Connections were hanging indefinitely waiting for slow queries.

**Lesson**: Always set explicit timeouts on database connections: connect timeout,
read timeout, and write timeout. Default "no timeout" is never acceptable in production.

**Application**: Add to connection config:
```go
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
```

---

## [2026-01-10-091500] Environment variables override config files

**Context**: Debugging why staging had different behavior than local. Config file
was correct, but an old environment variable was overriding it.

**Lesson**: Document the precedence order clearly: ENV > config file > defaults.
When debugging config issues, always check environment variables first.

**Application**: Add config source logging at startup:
```
Config loaded: database.host=localhost (source: ENV)
Config loaded: database.port=5432 (source: config.yaml)
```

---

## [2026-01-05-160030] Rate limiter must be per-user, not global

**Context**: Implemented global rate limiter (100 req/sec total). One heavy user
could starve all other users.

**Lesson**: Rate limiting should be per-user (or per-API-key) to ensure fair
resource allocation. Global limits are only useful as a last-resort circuit breaker.

**Application**: Use user ID or API key as the rate limiter bucket key, not a
single global bucket.

---
