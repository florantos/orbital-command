# ADR-003 - Use a monolithic deployment topology

**Status:** Accepted
**Date:** 2026-03-18

## Context

The project needs a deployment topology decision that is scalable based on our team size and scope of the project.

## Decision

Use a Monolithic deployment topology for our deployments.

## Options Considered

**Monolith**

- Pros: Simple to develop and deploy, no network overhead between project components, easy to debug, low operational complexity, faster iterations, transactions across the backend are straightforward
- Cons: You cannot scale components independently, tighter coupling risk as the codebase grows, deploying any change requires deploying the whole system, harder to use multi-technologies

**Microserviced**

- Pros: Independently deployable, independent scalability per service/component, technology flexibility per service/component, clear ownership boundaries, one service failing doesn't take the whole system offline
- Cons: Massive operational overhead, distributed transactions are complex and hard to manage, overkill for small teams

## Rationale

This is a solo-dev project, so we don't need to compartmentalize services and components for the sake of other teams or team members deploying independently from us. Microservices would be over-engineer for our case because we don't require a strict separation of services/components.

## Consequences

- Faster and simpler development
- A more easily debuggable system
- Transactions are simpler across the backend
- Risk of tight coupling as the codebase grows
- Single point of failure means entire system is offline if a process crashes
