package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type loginPageData struct {
	Error      string
	IsLoggedIn bool
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse templates once
	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/login.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		fmt.Println("template parse error:", err)
		return
	}

	_, ok := r.Context().Value(userIDKey).(int32)
	data := loginPageData{
		IsLoggedIn: ok,
		Error:      "",
	}
	switch r.Method {
	case http.MethodGet:
		// Render the login form
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "failed to render page", http.StatusInternalServerError)
			fmt.Println("template execute error:", err)
		}

	case http.MethodPost:
		// Handle form submission
		if err := r.ParseForm(); err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, dbErr := s.db.DB().GetUserByEmail(context.Background(), email)
		if dbErr != nil {
			data.Error = "❌ Incorrect email or password"
			tmpl.ExecuteTemplate(w, "base", data)
			return
		}

		if err := CheckPassword(user.PasswordHash, password); err != nil {
			http.Error(w, "❌ Incorrect email or password", http.StatusUnauthorized)
			return
		}
		// Generate token
		token, _ := GenerateToken(user.UserID)
		// Set token as HttpOnly cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})
		// Success
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
