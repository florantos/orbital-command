# ADR-011 — Use CSS Modules with design tokens for styling

**Status:** Accepted  
**Date:** 2026-03-18

## Context

The frontend needs a styling strategy that supports a complex, theme-heavy sci-fi aesthetic with consistent design tokens, component-level style isolation, and animated state transitions. The choice affects how styles are organized, how conflicts are prevented at scale, and how the design system is maintained.

## Decision

Use CSS Modules for style isolation combined with CSS custom properties (design tokens) for theming and consistency across components.

## Options Considered

**CSS Modules + design tokens**

- Pros: Style isolation at the component level means no class name collisions at scale, design tokens via CSS custom properties provide a single source of truth for colors, spacing, and typography, works with standard CSS, tokens are accessible in CSS and JavaScript, no runtime overhead, works naturally with the sci-fi theme's heavy use of color and animation
- Cons: More files per component (separate `.module.css` file), sharing styles across components requires explicit token usage or shared modules, no utility classes for rapid prototyping

**Utility-first (Tailwind)**

- Pros: Rapid prototyping, no context switching between files, consistent spacing and sizing scales built in, large ecosystem and community
- Cons: Class names become long and hard to read for complex components, customizing a sci-fi aesthetic requires significant Tailwind configuration, utility classes don't express intent e.g. `text-green-400` doesn't communicate "operational health state", HTML becomes cluttered, no pre-built Tailwind compiler available in this environment

**CSS-in-JS (styled-components, emotion)**

- Pros: Styles colocated with components, dynamic styles based on props are straightforward, full JavaScript power in styles
- Cons: Runtime overhead (styles are generated at runtime), increases bundle size, SSR complexity, harder to use CSS custom properties and animations naturally, adds a significant dependency

## Rationale

CSS Modules solve class name collisions without introducing runtime overhead or new syntax. Design tokens via CSS custom properties are the right fit for this project's sci-fi aesthetic: health state colors, severity indicators, surface colors, and spacing all need to be defined once and referenced everywhere. A token like `--color-health-operational` communicates intent in a way that `text-green-400` never could. Tailwind's utility-first approach works well for standard business UIs but fights against a heavily themed, animation-rich design system. CSS-in-JS adds runtime cost and complexity for no meaningful benefit over CSS Modules given the project's requirements.

## Consequences

- Style isolation at the component level
- Design tokens defined once in `index.css` and referenced everywhere
- Changing the theme requires updating one file
- Semantic token names communicate design intent (`--color-health-critical` vs a raw color value)
- No runtime overhead
- Separate `.module.css` file per component
- Animations and transitions are natural in standard CSS
