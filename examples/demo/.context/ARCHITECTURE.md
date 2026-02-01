# Architecture

System overview and component organization.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Clients                              │
│            (Web App, Mobile App, CLI)                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      Load Balancer                           │
│                    (nginx / AWS ALB)                         │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       API Server                             │
│                    (Go / net/http)                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Handlers  │  │  Services   │  │   Repositories      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────┬───────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
    ┌───────────┐  ┌───────────┐  ┌───────────────┐
    │ PostgreSQL│  │   Redis   │  │  Object Store │
    │ (primary) │  │  (cache)  │  │    (S3)       │
    └───────────┘  └───────────┘  └───────────────┘
```

## Directory Structure

```
.
├── cmd/
│   ├── api/           # API server entrypoint
│   └── worker/        # Background worker entrypoint
├── internal/
│   ├── handler/       # HTTP handlers
│   ├── service/       # Business logic
│   ├── repository/    # Data access
│   └── model/         # Domain types
├── pkg/               # Shared libraries (importable)
├── migrations/        # Database migrations
├── docs/              # Documentation
└── .context/          # AI context files
```

## Key Components

### API Server (`cmd/api`)
- Handles HTTP requests
- Validates input, calls services, returns responses
- Stateless — all state in database or cache

### Services (`internal/service`)
- Contains business logic
- Orchestrates multiple repositories
- Enforces business rules

### Repositories (`internal/repository`)
- Data access layer
- One repository per aggregate root
- Handles database queries and caching

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

## Data Flow

1. Request arrives at handler
2. Handler validates input, extracts user context
3. Handler calls service with validated data
4. Service applies business logic, calls repositories
5. Repository reads/writes to database
6. Response flows back up the stack

## Scaling Strategy

- **Horizontal**: Add more API server instances behind load balancer
- **Database**: Read replicas for read-heavy workloads
- **Cache**: Redis for session data and frequently accessed records
- **Background work**: Separate worker processes for async jobs
