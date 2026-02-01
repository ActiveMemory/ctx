# Session: Database Timeout Investigation

**Date**: 2026-01-15
**start_time**: 2026-01-15-140000
**end_time**: 2026-01-15-160000
**Topic**: Investigating production database connection issues
**Type**: bugfix

---

## Summary

Investigated production outage caused by database connection pool exhaustion.
Found that connections were hanging indefinitely on slow queries. Implemented
explicit timeouts and connection lifecycle management.

## Problem

- Production API started returning 503 errors
- Database connection pool was exhausted (all 100 connections in use)
- Connections were stuck waiting for queries that never returned
- No timeout configured on database connections

## Root Cause

Default Go database driver has no timeout. When the database is slow or
unresponsive, connections wait forever, eventually exhausting the pool.

## Fix Applied

```go
// Before: no timeouts
db, err := sql.Open("postgres", connStr)

// After: explicit lifecycle management
db, err := sql.Open("postgres", connStr)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)

// Query-level timeouts
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
rows, err := db.QueryContext(ctx, query)
```

## Key Decisions

- Set connection max lifetime to 5 minutes (prevents stale connections)
- Set query timeout to 30 seconds (fail fast on slow queries)
- Added circuit breaker for database calls

## Tasks for Next Session

- Add monitoring for connection pool metrics
- Set up alerting for connection pool utilization > 80%
- Review other services for similar timeout issues

## Files Changed

- `internal/repository/db.go`
- `internal/config/database.go`
- `docs/operations.md`
