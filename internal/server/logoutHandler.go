package server

import (
	"net/http"
	"time"
)

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Overwrite cookie with empty value and past expiry date
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
