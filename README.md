# ShortyLink

Modern URL Shortener.

## Preview

[video](https://github.com/user-attachments/assets/c26968cb-7161-42a8-bd50-1d4c6a2a702a)

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

| Variable              | Description                                                                 |
|------------------------|-----------------------------------------------------------------------------|
| `PORT`                | The port your application will run on (e.g. `5000`).                        |
| `APP_ENV`             | The application environment (`local`, `dev`, `staging`, or `production`).   |
| `BLUEPRINT_DB_HOST`   | Hostname of the PostgreSQL database (e.g. `psql_bp` if using Docker).        |
| `BLUEPRINT_DB_PORT`   | Port on which PostgreSQL is running (default is `5432`).                     |
| `BLUEPRINT_DB_DATABASE` | Name of the PostgreSQL database to connect to (e.g. `blueprint`).         |
| `BLUEPRINT_DB_USERNAME` | Username for database authentication (e.g. `melkey`).                     |
| `BLUEPRINT_DB_PASSWORD` | Password for database authentication.                                     |
| `BLUEPRINT_DB_SCHEMA` | Database schema to use (default is usually `public`).                        |
| `RESEND_API_KEY`      | API key for the [Resend](https://resend.com) email service.                  |
| `JWT_KEY`             | Secret key used for signing and verifying JWT tokens.                        |

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

- `build` → Build the Go executable as `main.exe` from `cmd/api/main.go`.

- `run` → Run the application directly with `go run`.

- `tailwind` → Compile Tailwind CSS from `input.css` to `output.css` in watch mode.

- `docker-run` → Build and start the Docker containers defined in `docker-compose.yml`.

- `docker-down` → Stop and remove all Docker containers defined in `docker-compose.yml`.

- `test` → Run all Go tests in the project.

- `itest` → Run integration tests for the database layer.

- `clean` → Remove the compiled Go binary `main`.

- `format` → Format all Go files and tidy the Go module dependencies.

## Future Improvements

There are several areas planned for enhancement in future iterations of the project:

- Implementation of a "Forgot Password" email feature  
- Email verification to prevent inactive or invalid addresses from filling the database  
- Brute force prevention by tracking failed login attempts per IP or account  
- Password change functionality and a user account management page  
- Shortened display name for the site when hosted (e.g. `sl.com`)  
- Utilization of existing Visits table columns for enhanced statistics
