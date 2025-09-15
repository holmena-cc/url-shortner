package server

import (
	"database/sql"
	"log"
	"my_project/internal/db"
	"net/http"
	"strings"
)

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := strings.TrimPrefix(r.URL.Path, "/r/")
	if alias == "" {
		http.Error(w, "missing custom alias", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	// 2. find original URL
	url, err := s.db.DB().GetUrlByAlias(ctx, alias)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	// 3. register the visit
	createVisitParams := db.CreateVisitParams{
		UrlID:     url.UrlID,
		IpAddress: "127.0.0.1",
		Referrer:  sql.NullString{String: "testing", Valid: true},
		Country:   sql.NullString{String: "XX", Valid: true},
	}
	_, err = s.db.DB().CreateVisit(ctx, createVisitParams)
	if err != nil {
		log.Println("failed to insert visit:", err)
	}
	// 4. redirect
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}
