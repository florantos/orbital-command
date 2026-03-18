# ADR-008 — Use hook-based data fetching layer for frontend internal structure

**Status:** Accepted  
**Date:** 2026-03-18

## Context

The frontend needs a defined internal structure that governs how components depend on data fetching and business logic. This is the frontend equivalent of the backend's internal structure decision. It defines dependency direction and separation of concerns across layers. The wrong choice couples components to infrastructure concerns and makes the codebase hard to test and maintain as it grows.

## Decision

Use a hook-based data fetching layer as the primary abstraction boundary. Components never call APIs directly. Instead they consume custom hooks that encapsulate fetch logic, loading state, and error state. Hooks call the API directly for now using fetch. Components are responsible for rendering only.

## Options Considered

**No separation — direct fetch in component**

- Pros: Simple, no abstraction overhead, easy to understand for small components
- Cons: Components become bloated with data fetching logic, impossible to reuse fetch logic, hard to test components in isolation, mixes rendering and data concerns, changing the API layer requires touching every component

**Hook-based layer** (component → hook → fetch)

- Pros: Components stay focused on rendering, data fetching logic is reusable, loading and error state handled consistently, hooks are testable independently from components
- Cons: Additional abstraction layer, more files per feature

**Hook + API client layer** (component → hook → API client → fetch)

- Pros: Maximum separation, hooks know nothing about HTTP, API client is fully swappable, mirrors hexagonal ports and adapters most closely, API client can be tested independently
- Cons: Higher overhead upfront, three layers of indirection for simple fetches, premature at this stage

**Service layer** (component → hook → service → API client → fetch)

- Pros: Full frontend equivalent of hexagonal architecture, complete dependency inversion at every layer
- Cons: Very high ceremony, significant boilerplate for a React app at this stage, overkill until domain logic on the frontend becomes complex

## Rationale

The hook-based layer is the right starting point. Components depend on hooks, not on fetch or API details. This is the frontend equivalent of the backend handler depending on the domain, not on the database directly. An API client layer is the natural evolution when the number of endpoints grows and the cost of duplicated fetch logic across hooks becomes real. The full service layer remains a future option if frontend domain logic grows significantly in complexity.

## Consequences

- Components stay lean and focused on rendering
- Data fetching logic is reusable across components within a feature
- Loading, error, and empty states are handled consistently inside hooks
- Components are easier to test
