package server

import (
	"database/sql"
	"errors"
	"log"
	"my_project/internal/db"
	"net"
	"net/http"
	"strings"
)

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := strings.TrimPrefix(r.URL.Path, "/r/")
	if alias == "" {
		http.Error(w, "missing alias", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	url, err := s.db.DB().GetUrlByAlias(ctx, alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.render404(w, r)
			return
		}
		log.Println("redirectHandler: db error:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Extract real client IP (respects X-Forwarded-For from Render's proxy)
	ip := realIP(r)

	referrer := r.Referer()

	createVisitParams := db.CreateVisitParams{
		UrlID:     url.UrlID,
		IpAddress: ip,
		Referrer:  sql.NullString{String: referrer, Valid: referrer != ""},
		Country:   sql.NullString{Valid: false},
	}
	if _, err := s.db.DB().CreateVisit(ctx, createVisitParams); err != nil {
		log.Println("redirectHandler: failed to insert visit:", err)
	}

	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

// realIP extracts the client IP, honouring X-Forwarded-For set by Render's proxy.
func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can be a comma-separated list; the first entry is the client.
		if ip := strings.TrimSpace(strings.SplitN(xff, ",", 2)[0]); ip != "" {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// render404 serves a friendly 404 page.
func (s *Server) render404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	// Reuse the base layout; fall back to plain text if templates fail.
	http.Error(w, "404 — that short link doesn't exist.", http.StatusNotFound)
}