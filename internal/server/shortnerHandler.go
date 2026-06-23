package server

import (
	"html/template"
	"math/rand"
	"my_project/internal/db"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type PageData struct {
	LongURL     string
	CustomAlias string
	Error       string
	ShortURL    string
	IsLoggedIn  bool
}

var aliasRe = regexp.MustCompile(`^[a-zA-Z0-9-]{3,30}$`)

func (s *Server) shortnerHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int32)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	longURL := r.FormValue("long_url")
	customAlias := r.FormValue("custom_alias")
	ctx := r.Context()

	homeTmpl, _ := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/footer.html",
		"web/templates/home.html",
	)

	data := PageData{
		LongURL:     longURL,
		CustomAlias: customAlias,
		IsLoggedIn:  true,
	}

	renderError := func(msg string) {
		data.Error = msg
		if err := homeTmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "failed to load template", http.StatusInternalServerError)
		}
	}

	// Validate long URL
	if longURL == "" {
		renderError("Please enter a URL to shorten")
		return
	}
	parsed, err := url.ParseRequestURI(longURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		renderError("Please enter a valid URL starting with http:// or https://")
		return
	}

	// Validate or generate alias
	if customAlias == "" {
		for {
			candidate := generateCustomAlias(6)
			exists, err := s.db.DB().AliasExists(ctx, candidate)
			if err != nil {
				renderError("Database error, please try again")
				return
			}
			if !exists {
				customAlias = candidate
				break
			}
		}
	} else {
		if !aliasRe.MatchString(customAlias) {
			renderError("Alias must be 3–30 characters, letters, numbers, and hyphens only")
			return
		}
		exists, err := s.db.DB().AliasExists(ctx, customAlias)
		if err != nil {
			renderError("Database error, please try again")
			return
		}
		if exists {
			renderError("Custom alias already taken")
			return
		}
	}

	// Build short URL from env (required in production)
	baseURL := os.Getenv("APP_BASE_URL")
	shortURL := baseURL + "/r/" + customAlias

	createURLParams := db.CreateURLParams{
		OriginalUrl: longURL,
		ShortCode:   shortURL,
		CustomAlias: customAlias,
		UserID:      userID,
	}
	if _, err := s.db.DB().CreateURL(ctx, createURLParams); err != nil {
		renderError("Failed to create short URL")
		return
	}

	// Success page
	successTmpl, _ := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/header.html",
		"web/templates/footer.html",
		"web/templates/shortened.html",
	)
	data.ShortURL = shortURL
	if err := successTmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
	}
}

// generateCustomAlias uses the package-level rand (auto-seeded in Go 1.20+).
func generateCustomAlias(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alias := make([]byte, n)
	for i := range alias {
		alias[i] = letters[rand.Intn(len(letters))]
	}
	return string(alias)
}