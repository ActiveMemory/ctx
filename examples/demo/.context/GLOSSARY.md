# Glossary

Domain terms, abbreviations, and project-specific vocabulary.

---

## Terms

- **Claim (JWT)**: A key-value pair embedded in a JWT token that
  carries user identity or permission data without a database lookup.

- **Connection pool**: A cache of database connections maintained
  so they can be reused for future requests instead of opening new
  connections each time.

- **Handler**: A function that processes an incoming HTTP request
  and returns a response. Lives in `internal/api/handlers/`.

- **Migration**: A versioned SQL script that modifies the database
  schema. Applied in order to bring the database to a target state.

- **Middleware**: A function that wraps handlers to add cross-cutting
  concerns (authentication, logging, rate limiting) without modifying
  handler logic.

- **Refresh token**: A long-lived token used to obtain new access
  tokens without re-authenticating. Stored server-side in the database.

- **Repository**: A data access layer that abstracts database
  operations behind a Go interface. One repository per domain entity.

- **Service**: A business logic layer between handlers and
  repositories. Enforces domain rules and orchestrates operations
  across multiple repositories.

## Abbreviations

- **API**: Application Programming Interface
- **JWT**: JSON Web Token (RFC 7519)
- **CRUD**: Create, Read, Update, Delete
- **DTO**: Data Transfer Object (request/response structs)
- **ORM**: Object-Relational Mapping (not used; raw SQL preferred)
