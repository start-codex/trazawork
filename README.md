# Taskcore

Taskcore is an open source project management platform for teams. Organize work on Kanban and Scrum-style boards, track issues, and ship faster.

## Goals

- Simple enough to start in minutes, flexible enough to grow with your team.
- Configurable per project: workflows, issue types, and boards adapt to how you work.
- Lightweight and self-hostable, inspired by the Gitea/Forgejo approach.

## Principles

- **Open source first** — clear contribution path, modular architecture, no vendor lock-in.
- **Project-level configuration** — each project defines its own statuses, issue types, and boards.
- **Minimal defaults** — ships with `To Do`, `In Progress`, `Done` out of the box.
- **Explicit over magic** — raw SQL, no ORM, no hidden layers.

## Tech stack

- **Backend:** Go — monolith, `database/sql` + `sqlx`, explicit SQL queries.
- **Frontend:** SvelteKit + shadcn-svelte.
- **Database:** PostgreSQL.
- **Deployment:** Single Docker image, `docker compose` for local dev.

## Roadmap

### Phase 1 — MVP (current)
- Workspaces and team membership.
- Projects with configurable statuses.
- Boards with columns mapped to statuses.
- Issues with title, description, status, priority, assignee, and due date.
- Drag and drop between columns.

### Phase 2 — Dev-centric
- Configurable issue types per project (Epic, Story, Task, Subtask).
- Parent/child relationships between issues.
- Sprint planning and backlog view.

### Phase 3 — Cross-industry
- Project templates (engineering, marketing, support, operations, legal).
- Basic automations (e.g. notify on status change).
- Metrics and reports.

## Getting started

```bash
# Clone and start
git clone https://github.com/start-codex/taskcore
cd taskcore
docker compose up --build

# App runs at http://localhost:8080
```

Create your first user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "name": "Your Name", "password": "yourpassword"}'
```

## API

All responses follow the envelope format:

```json
{ "status": 200, "data": {} }
{ "status": 400, "error": "description" }
```

## Documentation

- Product scope: [docs/01-product-scope.md](docs/01-product-scope.md)
- Architecture: [docs/02-architecture.md](docs/02-architecture.md)
- Data model: [docs/03-data-model.md](docs/03-data-model.md)
- Go conventions: [docs/04-go-conventions.md](docs/04-go-conventions.md)

## License

AGPL-3.0 — see [LICENSE](LICENSE) for details.
