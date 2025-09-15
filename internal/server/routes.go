package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public routes with optional auth
	mux.Handle("/", AuthMiddlewareOptional(http.HandlerFunc(s.homeHandler)))
	mux.Handle("/login", AuthMiddlewareOptional(http.HandlerFunc(s.loginHandler)))
	mux.Handle("/register", AuthMiddlewareOptional(http.HandlerFunc(s.registerHandler)))
	mux.Handle("/contact", AuthMiddlewareOptional(http.HandlerFunc(s.contactHandler)))
	mux.Handle("/api/contact", AuthMiddlewareOptional(http.HandlerFunc(s.contactFormHandler)))
	mux.Handle("/thankyou", AuthMiddlewareOptional(http.HandlerFunc(s.thankyouHandler)))
	mux.Handle("/logout", AuthMiddlewareOptional(http.HandlerFunc(s.logoutHandler)))
	mux.Handle("/r/", AuthMiddlewareOptional(http.HandlerFunc(s.redirectHandler)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("/health", AuthMiddlewareOptional(http.HandlerFunc(s.healthHandler)))

	// Protected routes (requires login)
	mux.Handle("/urls", AuthMiddleware(http.HandlerFunc(s.urlsHandler)))
	mux.Handle("/shortner", AuthMiddleware(http.HandlerFunc(s.shortnerHandler)))
	mux.Handle("/delete-url", AuthMiddleware(http.HandlerFunc(s.deleteURLHandler)))

	return mux
}
