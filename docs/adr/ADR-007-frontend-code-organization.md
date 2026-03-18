# ADR-007 — Use feature-based code organization for the frontend

**Status:** Accepted  
**Date:** 2026-03-18

## Context

The frontend needs a code organization strategy that determines how files and components are grouped. This decision affects how easy it is to navigate the codebase, reason about features, and scale the UI as the system grows.

## Decision

Use feature-based organization for the frontend: `src/features/`. Each feature directory contains everything related to that feature: components, hooks, types, and styles.

## Options Considered

**Feature-based** (`src/features/modules`, `src/features/alerts`, `src/features/crew`)

- Pros: Everything for a feature lives together, easy to reason about one feature at a time, natural encapsulation, easy to delete or extract a feature, scales well as the number of features grows
- Cons: Shared components and hooks need a separate home (`src/components/`, `src/hooks/`), feature boundaries need to be thought through upfront, can lead to duplication if shared patterns aren't extracted

**By type** (`src/components/`, `src/hooks/`, `src/pages/`, `src/services/`)

- Pros: Familiar to most React developers, easy to find all components or all hooks at a glance, low overhead for small apps
- Cons: A change to one feature touches multiple directories, no natural encapsulation of feature concerns, doesn't scale well — a large `src/components/` folder becomes hard to navigate, coupling between unrelated components is easy to introduce accidentally

## Rationale

Feature-based organization is chosen from the start for the frontend, unlike the backend which starts flat-by-layer. The frontend will grow to include modules, alerts, crew, escalation policies, etc. each will have its own components, hooks, and state. Feature-based organization keeps each feature self-contained and makes it obvious where new code belongs.

## Consequences

- Each feature is self-contained
- Shared components and utilities need a deliberate home (`src/components/`, `src/lib/`)
- Deleting or extracting a feature is straightforward since everything lives in one place
- This organization is stable for the lifetime of the project
