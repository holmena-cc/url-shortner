package server

import (
	"fmt"
	"html/template"
	"net/http"
)

type HomePageData struct {
	LongURL     string
	CustomAlias string
	Error       string
	IsLoggedIn  bool
}

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
	_, ok := r.Context().Value(userIDKey).(int32)
	data := HomePageData{
		LongURL:     "",
		CustomAlias: "",
		Error:       "",
		IsLoggedIn:  ok,
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		fmt.Println("failed to render err:", err)
	}
}
