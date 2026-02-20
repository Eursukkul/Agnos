# Agnos Candidate Assignment - Backend

Gin-based Hospital Middleware API with Postgres, Docker Compose, and Nginx.

## Quick Start

```bash
docker compose up --build
```

Service endpoints via Nginx:
- `http://localhost:8088/healthz`
- `http://localhost:8088/staff/create`
- `http://localhost:8088/staff/login`
- `http://localhost:8088/patient/search` (JWT required)

## Run Locally

```bash
go mod tidy
go test ./...
DATABASE_URL='postgres://postgres:postgres@localhost:5432/agnos?sslmode=disable' JWT_SECRET='very-secret-key' go run ./cmd/server
```

## API Examples

### Create Staff

```bash
curl -X POST http://localhost:8088/staff/create \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123","hospital":"hospital-a"}'
```

### Staff Login

```bash
curl -X POST http://localhost:8088/staff/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123","hospital":"hospital-a"}'
```

### Patient Search

```bash
curl -X POST http://localhost:8088/patient/search \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{"national_id":"1234567890123"}'
```

## Notes

- `patient/search` only returns patients in the same hospital as the staff token.
- If `national_id` or `passport_id` is provided and patient is missing in DB, middleware calls Hospital A API (`GET /patient/search/{id}`), stores result, then searches again.
- `gender` is constrained to `M`/`F`.

## Deliverables

- Project Structure: `docs/project-structure.md`
- API Spec: `docs/api-spec.md`
- ER Diagram: `docs/er-diagram.md`
