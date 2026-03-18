# ADR-005 — Use layered architecture for backend internal structure

**Status:** Accepted  
**Date:** 2026-03-18

## Context

The backend needs a defined internal structure that governs how layers depend on each other. The domain is moderately complex, it includes alert state machines, assignment logic, escalation policies, and cascade propagation which makes the choice of internal structure consequential.

## Decision

Use layered architecture for the backend internal structure.

## Options Considered

**Layered Architecture**

- Pros: Simple and widely understood, no interfaces required until needed, natural fit for early-stage development, fast to implement
- Cons: Domain can become coupled to persistence concerns, harder to test domain logic in isolation without infrastructure, swapping infrastructure requires touching domain code, can lead to bloated middle layers as complexity grows

**Hexagonal Architecture (Ports and Adapters)**

- Pros: Domain is completely isolated — it defines interfaces that infrastructure implements, easy to test domain logic without any infrastructure, clean boundary enforcement by design, infrastructure is swappable without touching domain
- Cons: Higher upfront overhead (interfaces and adapters required from the start), steeper learning curve, can feel over-engineered for simple CRUD operations early on

**Clean Architecture**

- Pros: Very explicit about what belongs where, good for large teams needing consistency
- Cons: High overhead, rigid layer definitions, significant boilerplate, typically overkill

## Rationale

Layered architecture is chosen as a deliberate starting point, not the permanent answer. The domain complexity of Orbital Command that includes alert state machines, escalation policies, and assignment logic, would justify hexagonal architecture. However, for the sake of learning, we will establish simple patterns first and refactored toward better patterns as complexity is felt firsthand. Starting with layered makes the pain points visible and makes the value of hexagonal concrete when the refactor arrives. Hexagonal remains the natural evolution of this codebase as domain complexity grows.

## Consequences

- Simple and fast to implement in early loops
- Handler directly orchestrates domain and repository calls
- Domain logic is harder to test in complete isolation from infrastructure
- As domain complexity grows, coupling between layers will become a pain point
