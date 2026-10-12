# AGENTS.md - serverplate Development Guide

serverplate is a Go web app that generates memorable, human-friendly server
names ("brave-mountain") through a web UI and a REST API. It is intended for
manual server naming — **not** for automated system provisioning.

Most things about this codebase (directory layout, tech stack, routes, env
vars, naming patterns) are self-evident from the code. This file only records
what is **not** easily inferable: intent, non-obvious constraints, and
workflow rules. If a fact below contradicts the code, the code wins.

## Commands (mise)

`mise.toml` is the build tool. Important tasks:

| Task | What it does |
|------|--------------|
| `mise run dev` | Hot-reload server (Air); regenerates templ on save |
| `mise run dev:frontend` | Vite dev server — run in a second terminal alongside `dev` |
| `mise run generate` | Regenerate templ + oapi-codegen output |
| `mise run build` | Production binary into `build/` |
| `mise run test` | `go test ./...` |
| `mise run lint` / `mise run format` | Style checks / auto-fix + templ fmt |
| `mise run db:migrate` | Apply pending migrations (dbmate) |
| `mise run db:reset` | **Drops the database**, recreates, migrates — dev only |
| `mise run db:seed` | Sync adjectives/nouns from files into DB |
| `mise run docker` | Build production image |

## Non-obvious architecture

### SQLite dual connection pools

The DB layer uses two separate pools against the same SQLite file to avoid
lock contention:

- **Write pool**: all INSERT/UPDATE/DELETE *and transactions*. `MaxOpenConns = 1`
  (SQLite serializes writes), DSN uses `_txlock=immediate`.
- **Read pool**: SELECTs only. `MaxOpenConns = max(4, runtime.NumCPU())`.

Every store method must explicitly pick the right pool: reads via
`db.Read()`, writes and transactions via `db.Write()` / `Write().WithTx()`.
See `internal/store/sqlitestore/` for the established pattern.

### Bucket lifecycle

- Buckets go Active → Archived → Deleted.
- Archived buckets are read-only: popping a name or updating them returns
  **409** (Problem Detail).
- A background task (`remove_archived_buckets`, hourly) deletes buckets
  archived more than **3 days** ago.
- Bucket cursor: each pop advances the cursor until the bucket's names are
  exhausted.

### Seed data

`db:seed` manages rows with `from_seed = 1`: words present in `db/seed/*.txt`
are added, missing ones are removed. Manually added words (`from_seed = 0`)
are preserved. Seed files are the source of truth for word counts — don't
copy counts around; read them from the DB or files.

## Generated code

`openapi.yaml` is the source of truth for the API contract. Adding or changing
an API endpoint means: edit the spec → `mise run generate` → implement the
generated method in `internal/server/api/api_handlers.go`. Method names match
the spec's `operationId`.

Files produced by code generation (`*_gen.go`, `*_templ.go`) and generated
assets are safe to read but **must never be edited by hand** — regenerate
instead.

To add a web endpoint, a template, or a background task: copy the closest
existing example. Good ones:

- API handlers: `internal/server/api/api_handlers.go`
- Web handlers + routes: `internal/server/*_handlers.go`, `internal/server/routes.go`
- Templates: any `internal/templates/*_page.templ` / `*_partial.templ`
- Background task: `internal/bg/remove_archived_buckets_task.go` + `runner.go`

## Conventions that aren't guessable

- **Form fields / query strings are snake_case** (`length_mode`,
  `filter_length_value`), read with `r.FormValue(...)`. This mirrors the REST
  API's snake_case JSON for consistency; keep both sides aligned.
- **Cron task names must be snake_case** — enforced at runtime via regex in
  `internal/bg/runner.go`, which *panics* on violation.
- **API errors**: 4xx errors return RFC 7807 Problem Details; 5xx return the
  generic `Error` schema. All API routes live under `/api/v1alpha1/`; use
  plural resource nouns.
- **Name format**: `{adjective}-{noun}`, validated RFC 1123 DNS-subdomain
  compliant (lowercase, hyphens, starts/ends alphanumeric, ≤63 chars):
  `^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`
- Templates use a `Page` suffix for full pages and `Partial` for fragments.

## Gotchas

- `mise run db:reset` wipes the database — don't run it on data you care about.
- Assets work differently per mode: dev (`ASSETS_MANIFEST_USE=false`) serves
  via the Vite dev server from the filesystem; prod builds `frontend/dist/`
  and embeds it with a manifest. If assets 404, check `ASSETS_MANIFEST_USE`,
  `ASSETS_ROOT_URL`, and that `frontend/dist/manifest.json` matches the mode.
- Air auto-regenerates templ code on save; if codegen is stale outside a dev
  session, run `mise run generate` manually.

## Where to look

| File | Purpose |
|------|---------|
| `openapi.yaml` | API contract (source of truth) |
| `internal/server/routes.go` | All web/API route registrations |
| `mise.toml` | Every task + tool versions |
| `db/seed/` | Adjective/noun word lists |
| `Dockerfile` | Deployment story (multi-stage, distroless) |

## Frontend aesthetics

<frontend_aesthetics>
You tend to converge toward generic, "on distribution" outputs. In frontend design, this creates what users call the "AI slop" aesthetic. Avoid this: make creative, distinctive frontends that surprise and delight. Focus on:

Typography: Choose fonts that are beautiful, unique, and interesting. Avoid generic fonts like Arial and Inter; opt instead for distinctive choices that elevate the frontend's aesthetics.

Color & Theme: Commit to a cohesive aesthetic. Use CSS variables for consistency. Dominant colors with sharp accents outperform timid, evenly-distributed palettes. Draw from IDE themes and cultural aesthetics for inspiration.

Motion: Use animations for effects and micro-interactions. Prioritize CSS-only solutions for HTML. Use Motion library for React when available. Focus on high-impact moments: one well-orchestrated page load with staggered reveals (animation-delay) creates more delight than scattered micro-interactions.

Backgrounds: Create atmosphere and depth rather than defaulting to solid colors. Layer CSS gradients, use geometric patterns, or add contextual effects that match the overall aesthetic.

Avoid generic AI-generated aesthetics:
- Overused font families (Inter, Roboto, Arial, system fonts)
- Clichéd color schemes (particularly purple gradients on white backgrounds)
- Predictable layouts and component patterns
- Cookie-cutter design that lacks context-specific character

Interpret creatively and make unexpected choices that feel genuinely designed for the context. Vary between light and dark themes, different fonts, different aesthetics. You still tend to converge on common choices (Space Grotesk, for example) across generations. Avoid this: it is critical that you think outside the box!
</frontend_aesthetics>
