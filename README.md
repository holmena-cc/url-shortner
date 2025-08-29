# ShortyLink

Modern URL Shortener.

## Preview

Here is a quick preview [video]() of the website.

- **Login Screen:**

<img width="1000" alt="login" src="https://github.com/user-attachments/assets/0f78778f-c710-4fae-8632-a4a1ef319b38" /> <br>

- **Home Screen:**

<img width="1000" alt="home" src="https://github.com/user-attachments/assets/d4bd039a-0140-4ef9-b952-f2ffac365b3f" />

## Tech Stack

- **Backend:** Go
- **Database:** PostgreSQL
- **Frontend:** Tailwind CSS (requires Node.js for compilation)
- **Containerization:** Docker

## Features

- User authentication.
- Quickly shorten and share URLs.
- Create custom aliases for links.
- Track and manage your shortened URLs.
- Secure passwords hashed with bcrypt.

## Environment Variables

- `PORT` – Port for the Go server (e.g. 5000)
- `APP_ENV` – App environment (local, production, etc.)
- `BLUEPRINT_DB_HOST` – DB container host (e.g. psql_bp)
- `BLUEPRINT_DB_PORT` – DB port (e.g. 5432)
- `BLUEPRINT_DB_DATABASE` – Database name (e.g. blueprint)
- `BLUEPRINT_DB_USERNAME` – DB username
- `BLUEPRINT_DB_PASSWORD` – DB password
- `BLUEPRINT_DB_SCHEMA` – DB schema (e.g. public)

## Database Note

PostgreSQL executes initialization SQL files only the first time a database is created in a volume.  
To apply new migrations or updates after that, you must run the SQL files manually inside the container.

1. Open an interactive PostgreSQL session:
```bash
docker exec -it <container_name> psql -U <username> -d <database>
```
2. Run the SQL migration file:
```sql
\i /path/to/migration.sql
```

You can also run queries and inspect tables directly in this session for testing or debugging.

## Makefile

- `build`  
  Build the Go executable as `main.exe` from `cmd/api/main.go`.

- `run`  
  Run the application directly with `go run`.

- `tailwind`  
  Compile Tailwind CSS from `input.css` to `output.css` in watch mode.

- `docker-run`  
  Build and start the Docker containers defined in `docker-compose.yml`.

- `docker-down`  
  Stop and remove all Docker containers defined in `docker-compose.yml`.

- `test`  
  Run all Go tests in the project.

- `itest`  
  Run integration tests for the database layer.

- `clean`  
  Remove the compiled Go binary `main`.

- `format`  
  Format all Go files and tidy the Go module dependencies.

