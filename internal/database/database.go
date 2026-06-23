package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"my_project/internal/db"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	Close() error
	DB() *db.Queries
}

type service struct {
	db      *sql.DB
	Queries *db.Queries
}

var dbInstance *service

// connString resolves the connection string.
// Priority: DATABASE_URL (Render / any managed Postgres)
// Fallback:  individual BLUEPRINT_DB_* vars (local Docker / legacy)
func connString() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	username := os.Getenv("BLUEPRINT_DB_USERNAME")
	password := os.Getenv("BLUEPRINT_DB_PASSWORD")
	database := os.Getenv("BLUEPRINT_DB_DATABASE")
	schema := os.Getenv("BLUEPRINT_DB_SCHEMA")

	if schema == "" {
		schema = "public"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		username, password, host, port, database, schema,
	)
}

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	sqlDB, err := sql.Open("pgx", connString())
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(sqlDB)

	dbInstance = &service{
		db:      sqlDB,
		Queries: queries,
	}
	return dbInstance
}

func (s *service) DB() *db.Queries {
	return s.Queries
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}