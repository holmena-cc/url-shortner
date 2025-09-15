package server

import (
	"fmt"
	"my_project/internal/db"
	"net/http"
)

func (s *Server) deleteURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	shortCode := r.FormValue("shortCode")
	userID := r.Context().Value(userIDKey).(int32)
	deleteURLParams := db.DeleteURLParams{
		ShortCode: shortCode,
		UserID:    userID,
	}
	err := s.db.DB().DeleteURL(r.Context(), deleteURLParams)
	if err != nil {
		http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		fmt.Println("delete error:", err)
		return
	}
	http.Redirect(w, r, "/urls", http.StatusSeeOther)
}
