# Orbital Command

![CI](https://github.com/florantos/orbital-command/actions/workflows/ci.yml/badge.svg)

A space station systems control platform for monitoring modules, managing crew, and responding to alerts in real time.

## Prerequisites

- [Docker](https://www.docker.com/get-started)

## Getting Started

Clone the repository and create your environment file:

```bash
git clone https://github.com/florantos/orbital-command
cd orbital-command
cp .env.example .env
```

Start the stack:

```bash
docker compose up
```

The application will be available at:

- Frontend: http://localhost:5173
- API: http://localhost:8080
- Health check: http://localhost:8080/health

## Running Tests

**Backend:**

```bash
cd backend
go test ./...
```

**Frontend:**

```bash
cd frontend
pnpm lint
pnpm build
```

## License

MIT
