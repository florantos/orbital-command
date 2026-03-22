# ADR-013 - Use Chi as the HTTP router

**Status:** Accepted  
**Date:** 2026-03-22

## Context

The project uses Go's standard `net/http` mux for routing. As the API grows, multiple HTTP methods need to be registered on the same route (e.g. GET and POST on `/modules`). The standard mux does not support method-specific registration. Workarounds require either a dispatcher function in `main.go` or a method-switch handler on the `Handler` struct. Both approaches push routing concerns into the wrong layer and become harder to maintain as routes grow.

## Decision

Use Chi as the HTTP router.

## Options Considered

**Standard `net/http` mux with inline dispatcher**

- Pros: No external dependency, idiomatic for simple cases
- Cons: No method-specific registration, requires a manual dispatcher per route, method validation logic either duplicated in each handler or centralized in an awkward switch block in `main.go`, scales poorly as routes grow

**Chi router**

- Pros: Method-specific registration (`r.Get`, `r.Post`), clean route definitions in one place, middleware support built in, idiomatic Go, no magic, widely used in production Go projects
- Cons: External dependency, small learning curve

**Gorilla Mux**

- Pros: Mature, feature-rich, method-specific registration
- Cons: Heavier than Chi, Gorilla org went unmaintained for a period before being revived, Chi is the more active and widely recommended choice for new projects

## Rationale

Chi solves a real limitation of the standard mux that we are already working around. The dispatcher workaround would become the pattern for every future route in the project. Introducing Chi now is cleaner than refactoring it out later. Chi is lightweight, has no magic, and is the professional standard for Go HTTP routing. The dependency cost is minimal relative to the routing clarity it provides across the full project lifecycle.

## Consequences

- Method-specific route registration across all handlers
- Method validation removed from individual handlers
- Middleware can be applied per-route or per-group as needed in later loops
- One external dependency added to the backend
