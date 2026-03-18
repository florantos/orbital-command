# ADR-004 - Use a monorepo structure

**Status:** Accepted
**Date:** 2026-03-18

## Context

The project needs a repo structure that makes sense based on our team size and deployment schedule of our project.

## Decision

Use a monorepo structure.

## Options Considered

**Monorepo**

- Pros: All code in one place, atomic commits across the frontend and backend, shared tooling and configuration, no version sync issues
- Cons: Large repo over time, full CI run for backend and frontend on every change, more discipline to maintain boundaries between frontend and backend

**Polyrepo**

- Pros: Clear ownership boundaries per repo, CI only runs for what is changed, teams can work independently,
- Cons: Syncing changes across repos requires two PRs, two reviews, and two merges, cross repo refactoring is harder, version drift

## Rationale

Given that we are a single developer working on this project, monorepo makes more sense. The velocity and ease of deployment we get from a monorepo as a solo dev outweighs the slower ci deployment and clear ownership boundaries of a polyrepo.

## Consequences

- Simpler deployments
- Easier to implement features since the code is all in one place
- Frontend and Backend in a single PR instead of two means code reviewers have the full context
- Slower CI pipeline since it needs to always run backend and frontend checks
- Requires more discipline to keep backend and frontend separate
