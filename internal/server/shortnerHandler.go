package server

import (
	"html/template"
	"math/rand"
	"my_project/internal/db"
	"net/http"
	"time"
)

type ShortURLPageData struct {
	LongURL  string
	ShortURL string
}

func (s *Server) shortnerHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	longURL := r.FormValue("long_url")
	customAlias := r.FormValue("custom_alias")
	ctx := r.Context()

	// Struct to pass data to template
	type PageData struct {
		LongURL     string
		CustomAlias string
		Error       string
	}

	// Check for empty long URL
	if longURL == "" {
		tmpl, _ := template.ParseFiles(
			"web/templates/base.html",
			"web/templates/header.html",
			"web/templates/footer.html",
			"web/templates/home.html",
		)
		data := PageData{
			LongURL:     longURL,
			CustomAlias: customAlias,
			Error:       "Please enter a URL to shorten",
		}
		err := tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, "failed to load template", http.StatusInternalServerError)
		}
		return
	}

	// Generate alias if empty
	if customAlias == "" {
		for {
			candidate := generateCustomAlias(5)
			exists, err := s.db.DB().AliasExists(ctx, candidate)
			if err != nil {
				tmpl, _ := template.ParseFiles("web/templates/base.html", "web/templates/header.html", "web/templates/footer.html", "web/templates/home.html")
				data := PageData{
					Error: "Database error, please try again",
				}
				err = tmpl.ExecuteTemplate(w, "base", data)
				if err != nil {
					http.Error(w, "failed to load template", http.StatusInternalServerError)
				}
				return
			}
			if !exists {
				customAlias = candidate
				break
			}
		}
	} else {
		exists, err := s.db.DB().AliasExists(ctx, customAlias)
		if err != nil {
			tmpl, _ := template.ParseFiles("web/templates/base.html", "web/templates/header.html", "web/templates/footer.html", "web/templates/home.html")
			data := PageData{
				Error: "Database error, please try again",
			}
			err = tmpl.ExecuteTemplate(w, "base", data)
			if err != nil {
				http.Error(w, "failed to load template", http.StatusInternalServerError)
			}
			return
		}
		if exists {
			tmpl, _ := template.ParseFiles("web/templates/base.html", "web/templates/header.html", "web/templates/footer.html", "web/templates/home.html")
			data := PageData{
				LongURL:     longURL,
				CustomAlias: customAlias,
				Error:       "Custom alias already taken",
			}
			err = tmpl.ExecuteTemplate(w, "base", data)
			if err != nil {
				http.Error(w, "failed to load template", http.StatusInternalServerError)
			}
			return
		}
	}

	// Create short URL
	shortUrl := "http://localhost:5000/r/" + customAlias
	createURLParams := db.CreateURLParams{
		OriginalUrl: longURL,
		ShortCode:   shortUrl,
		CustomAlias: customAlias,
		UserID:      2,
	}
	_, err := s.db.DB().CreateURL(ctx, createURLParams)
	if err != nil {
		tmpl, _ := template.ParseFiles("web/templates/base.html", "web/templates/header.html", "web/templates/footer.html", "web/templates/home.html")
		data := PageData{
			Error: "Failed to create short URL",
		}
		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, "failed to load template", http.StatusInternalServerError)
		}
		return
	}

	// Success: show shortened URL page
	tmpl, _ := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/footer.html",
		"web/templates/shortened.html",
	)
	data := ShortURLPageData{
		LongURL:  longURL,
		ShortURL: shortUrl,
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
	}
}


func generateCustomAlias(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	alias := make([]rune, n)
	for i := range alias {
		alias[i] = letters[r.Intn(len(letters))]
	}
	return string(alias)
}
