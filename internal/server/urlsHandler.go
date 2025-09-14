package server

import (
	"fmt"
	"html/template"
	"my_project/internal/db"
	"net/http"
)

func (s *Server) urlsHandler(w http.ResponseWriter, r *http.Request) {
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
	var userId int32 = 2

	urls, err := s.db.DB().ListURLsByUserWithClicks(ctx, userId)
	if err != nil {
		http.Error(w, "failed to fetch urls", http.StatusInternalServerError)
		fmt.Println("db error:", err)
		return
	}
	totalClicks := 0
	for _, u := range urls {
		totalClicks += int(u.Clicks)
	}
	type URLsPageData struct {
		UrlsCount   int
		TotalClicks int
		URLs        []db.ListURLsByUserWithClicksRow
	}
	data := URLsPageData{
		UrlsCount:   len(urls),
		TotalClicks: totalClicks,
		URLs:        urls,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		fmt.Println("failed to render err:", err)
	}
}
