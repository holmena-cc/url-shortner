package server

import (
	"fmt"
	"html/template"
	"my_project/internal/db"
	"net/http"
)
type URLsPageData struct {
	UrlsCount   int
	TotalClicks int
	URLs        []db.ListURLsByUserWithClicksRow
	IsLoggedIn  bool
}

func (s *Server) urlsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int32)
	
	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/urls.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	ctx := r.Context()
	urls, err := s.db.DB().ListURLsByUserWithClicks(ctx, userID)
	if err != nil {
		http.Error(w, "failed to fetch urls", http.StatusInternalServerError)
		fmt.Println("db error:", err)
		return
	}
	totalClicks := 0
	for _, u := range urls {
		totalClicks += int(u.Clicks)
	}
	data := URLsPageData{
		UrlsCount:   len(urls),
		TotalClicks: totalClicks,
		URLs:        urls,
		IsLoggedIn:  true,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		fmt.Println("failed to render err:", err)
	}
}
