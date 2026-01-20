# Architecture

System overview and key design decisions.

## System Overview

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│   Frontend   │───▶│  API Server  │───▶│   Database   │
│   (React)    │    │    (Go)      │    │  (Postgres)  │
└──────────────┘    └──────────────┘    └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  Redis Cache │
                    └──────────────┘
```

## Directory Structure

- `/cmd/` - Application entry points
- `/internal/` - Private application code
- `/pkg/` - Public library code
- `/api/` - API definitions and OpenAPI specs
- `/web/` - Frontend React application

## Key Patterns

### Repository Pattern
Data access is abstracted through repositories. Business logic never
directly queries the database.

### Dependency Injection
All dependencies are injected through constructors, making testing
easier and components more modular.

### Event-Driven Updates
The system uses an event bus for decoupled component communication.
Events are published when state changes, and interested components
subscribe to relevant events.
