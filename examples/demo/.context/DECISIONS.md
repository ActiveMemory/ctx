# Decisions

Architectural decisions with rationale and consequences.

---

## [2026-01-05-110000] Use PostgreSQL for Primary Database

**Context**: Needed to choose a database for the application. Options were
PostgreSQL, MySQL, and MongoDB.

**Decision**: PostgreSQL

**Rationale**:
- Strong ACID compliance for financial transactions
- Excellent JSON support for flexible schema needs
- Team has existing PostgreSQL expertise
- Rich ecosystem of tools and extensions

**Consequences**:
- Need to manage schema migrations explicitly
- Requires more upfront schema design than document stores
- Horizontal scaling requires additional tooling (Citus, read replicas)

---

## [2026-01-08-140000] JWT for API Authentication

**Context**: Needed to choose authentication mechanism for the REST API.
Options were session cookies, JWT tokens, and API keys.

**Decision**: JWT tokens with short expiry + refresh tokens

**Rationale**:
- Stateless authentication scales horizontally without session storage
- Works well for both web and mobile clients
- Can embed user claims to reduce database lookups
- Industry standard with good library support

**Consequences**:
- Cannot immediately revoke tokens (must wait for expiry)
- Need secure storage for refresh tokens
- Must implement token refresh flow in all clients
- Larger request payload than session cookies

---

## [2026-01-10-090000] Use Go for API Server

**Context**: Choosing a backend language for the API. Options were Go,
Node.js, and Python.

**Decision**: Go

**Rationale**:
- Excellent performance characteristics
- Strong typing catches bugs at compile time
- Simple deployment with single binary
- Great concurrency primitives for handling many connections

**Consequences**:
- Smaller talent pool than JavaScript/Python
- Some team members need Go training
- Compile step required (vs interpreted languages)

---

## [2026-01-12-160000] Monorepo Structure

**Context**: Starting with multiple services (API, worker, CLI). Needed to
decide between monorepo and multi-repo structure.

**Decision**: Monorepo with shared packages

**Rationale**:
- Atomic commits across services
- Easier code sharing and refactoring
- Single CI/CD pipeline to maintain
- Better visibility into cross-service changes

**Consequences**:
- Need tooling to handle partial builds (only changed services)
- Repository will grow large over time
- All developers need access to entire codebase
- Must establish clear package boundaries

---
