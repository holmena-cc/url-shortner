package server

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux() // Register routes
	// mux.HandleFunc("/health",s.healthHandler)
	mux.HandleFunc("/", s.homeHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) LoadTemplates() error {
	// Parse all templates in the templates folder
	tmpl, err := template.ParseGlob(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
