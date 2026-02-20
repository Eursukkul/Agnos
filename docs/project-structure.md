# Project Structure

```text
.
├── cmd/server/main.go
├── internal
│   ├── config
│   ├── db
│   ├── his
│   ├── http
│   ├── middleware
│   ├── model
│   ├── repository
│   └── service
├── db/init/001_init.sql
├── nginx/default.conf
├── docker-compose.yml
├── Dockerfile
└── docs
```

## Layering

- `http`: transport handlers and routing
- `service`: business logic and policy
- `repository`: persistence access (Postgres)
- `his`: external Hospital A API integration client
- `middleware`: JWT auth and hospital scoping
