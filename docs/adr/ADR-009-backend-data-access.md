# ADR-009 — Use raw SQL for data access

**Status:** Accepted  
**Date:** 2026-03-18

## Context

The backend needs a data access strategy for interacting with Postgres. This decision affects how queries are written, how much abstraction sits between the application and the database, and how easy it is to write complex queries as the domain grows.

## Decision

Use raw SQL with the `pgx` driver for all database interactions. No ORM, no query builder at this stage.

## Options Considered

**Raw SQL (pgx)**

- Pros: Full control over queries, no abstraction overhead, easy to optimize, SQL is explicit and readable, no hidden N+1 queries, straightforward to write complex joins and filters, no migration from one abstraction to another as queries grow complex
- Cons: More verbose, no compile-time query safety, SQL strings can drift from schema without being caught until runtime, more boilerplate for common CRUD operations

**Query builder (squirrel)**

- Pros: Programmatic query construction which useful for dynamic queries with optional filters, type-safe query building, reduces string concatenation for complex dynamic queries, sits close to raw SQL so migration cost is low
- Cons: Additional dependency, adds abstraction over SQL that can obscure what query is actually being run, overkill for simple static queries

**ORM (GORM, ent)**

- Pros: Reduces boilerplate for CRUD operations, schema defined in Go code, auto-migration support, relationships handled automatically
- Cons: black boxed behavior that can produce unexpected queries, N+1 query problems are easy to introduce accidentally, complex queries require dropping down to raw SQL anyway, higher migration cost if you need to move away later, abstracts away SQL knowledge that is valuable to develop, starting with ORM on a greenfield project is a higher-cost migration path than starting with raw SQL and adding an ORM later

## Rationale

Raw SQL gives full control and keeps the data access layer explicit. The domain will eventually require complex queries, dependency graph traversal, aggregations for crew assignment. These queries are clearer and more optimizable in raw SQL than through an ORM abstraction. The verbosity cost of raw SQL is acceptable. The repository pattern already abstracts persistence from the domain, so the SQL complexity is contained. A query builder like squirrel remains a future option for dynamic query construction when optional filters make string concatenation unwieldy. Starting with raw SQL and adding abstraction later is a lower-cost migration path than starting with an ORM and removing it.

## Consequences

- Full control over all queries
- No hidden behavior or unexpected query plans
- SQL knowledge is built and maintained directly
- More verbose for simple CRUD operations
- Schema drift between SQL strings and actual schema will not be caught at compile time
- squirrel remains a future option for dynamic query construction as filtering requirements grow
- No ORM migration cost
