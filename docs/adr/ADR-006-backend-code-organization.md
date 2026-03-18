# ADR-006 — Use flat-by-layer code organization for the backend

**Status:** Accepted  
**Date:** 2026-03-18

## Context

The backend needs a code organization strategy that determines how files and packages are grouped. This decision affects how easy it is to navigate the codebase, reason about features, and enforce layer boundaries as the system grows.

## Decision

Use flat-by-layer organization for the backend e.g.`internal/handler`, `internal/domain`, `internal/repository`. All handlers live together, all domain types live together, all repositories live together.

## Options Considered

**Flat by layer**

- Pros: Simple to start, easy to find all handlers or all repositories at a glance, low overhead for small codebases, natural fit for layered architecture
- Cons: As features grow, files within each layer become hard to navigate, a change to one feature touches multiple directories, no natural encapsulation of feature concerns

**Package by feature** (`internal/modules`, `internal/alerts`, `internal/crew`)

- Pros: Everything for a feature lives together, high cohesion, easy to reason about one feature at a time, natural encapsulation, easier to extract a feature later if needed
- Cons: Harder to enforce layer boundaries across features, can lead to duplication of patterns across feature packages, more upfront thinking required about feature boundaries

## Rationale

Flat by layer is the right starting point for a codebase with one entity. The friction of flat-by-layer becomes real around 3-4 features when navigating across directories for a single change becomes painful. That friction is intentional and it will makes the value of package-by-feature concrete when that pain is felt. Starting with flat-by-layer keeps cognitive overhead low while the domain and patterns are still being established.

## Consequences

- Simple to navigate at this stage
- A change to one feature touches multiple directories
- Layer boundaries are visually enforced by the folder structure
