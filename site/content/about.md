---
title: About
slug: about
---
# About S4G

S4G was built to make static site generation **simple and extensible**.

## Architecture

The generator is split into clean layers:

- **`core.DataSource`** — anything that produces `[]Content`
- **`core.Renderer`** — anything that consumes `[]Content` and writes output
- **`core.Generator`** — wires sources to renderers

## Adding a New Data Source

Create a struct that implements `core.DataSource` and register it in `cmd/build.go`.
No other files need to change.

## Future Plans

- JSON file source
- REST API source
- SQL (PostgreSQL / SQLite) source
- MongoDB source
- GraphQL query source
- Local dev server (`s4g serve`)
