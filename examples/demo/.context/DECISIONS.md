# Decisions

Record of significant technical decisions with rationale.

## 2024-01-15 Use PostgreSQL for primary database

**Context**: We need a reliable, scalable database for our application.

**Decision**: Use PostgreSQL instead of MySQL or MongoDB.

**Rationale**:
- Strong ACID compliance for financial transactions
- Excellent JSON support for flexible schema needs
- Rich ecosystem of tools and extensions
- Team has existing expertise

**Consequences**:
- Need to learn PostgreSQL-specific features
- Deployment requires PostgreSQL setup

---

## 2024-01-10 Use Go for API server

**Context**: Choosing a backend language for our API.

**Decision**: Use Go instead of Node.js or Python.

**Rationale**:
- Excellent performance characteristics
- Strong typing catches bugs at compile time
- Simple deployment with single binary
- Great concurrency primitives

**Consequences**:
- Smaller talent pool than JavaScript
- Some team members need Go training
