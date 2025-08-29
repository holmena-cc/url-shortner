package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/home.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		fmt.Println("failed to render err:", err)
	}
}
