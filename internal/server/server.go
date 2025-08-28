package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"my_project/internal/database"
)

type Server struct {
	port      int
	templates *template.Template
	db        database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	srv := &Server{
		port: port,
		db:   database.New(),
	}
	if err := srv.LoadTemplates(); err != nil {
		panic(fmt.Sprintf("Failed to load templates: %v", err))
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
