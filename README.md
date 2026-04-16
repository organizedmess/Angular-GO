# URL Shortener (Angular + Go)

Base project scaffold for a URL shortener using:
- Frontend: Angular
- Backend: Go (Gin + GORM + PostgreSQL)

## Project Structure

- `frontend/` Angular app
- `backend/` Go API server
- `docker-compose.yml` local PostgreSQL

## Run Backend

1. Start PostgreSQL:
   ```bash
   docker compose up -d
   ```
2. Start API:
   ```bash
   cd backend
   cp .env.example .env
   go run .
   ```

Backend runs on `http://localhost:8080`.

## Run Frontend

```bash
cd frontend
npm install
npm start
```

Frontend runs on `http://localhost:4200`.

## Available API Endpoints

- `POST /api/shorten`
- `GET /:code`
- `GET /api/stats/:code`
- `GET /api/health`
