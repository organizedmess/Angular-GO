# URL Shortener (Angular + Go)

Base project scaffold for a URL shortener using:
- Frontend: Angular
- Backend: Go (Gin + GORM + PostgreSQL)

# Screenshots
<img width="2880" height="1566" alt="image" src="https://github.com/user-attachments/assets/0bb3c549-fab2-4b73-aff0-3c47a722f33f" />


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
