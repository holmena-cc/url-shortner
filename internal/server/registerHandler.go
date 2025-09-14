package server

import (
	"context"
	"fmt"
	"html/template"
	"my_project/internal/db"
	"net/http"
)
type LoginPageData struct {
	Error string
	Email string
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse templates once
	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/register.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		fmt.Println("template parse error:", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Render the login form
		if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
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
		ConfirmPassword := r.FormValue("ConfirmPassword")

		_, dbErr := s.db.DB().GetUserByEmail(context.Background(), email)
		// If a user with this email already exists
		if dbErr == nil {
			tmpl.ExecuteTemplate(w, "base", LoginPageData{
				Error: "❌ User with this email already exists",
				Email: "",
			})
			return
		}
		if password != ConfirmPassword {
			tmpl.ExecuteTemplate(w, "base", LoginPageData{
				Error: "❌ Passwords do not match",
				Email: email,

			})
			return
		}
		passwordHash, err := HashPassword(password)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		_, dbErr = s.db.DB().CreateUser(context.Background(), db.CreateUserParams{
			Email:        email,
			PasswordHash: passwordHash,
		})
		if dbErr != nil {
			http.Error(w, "failed to create account", http.StatusBadRequest)
			return
		}

		// Success
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
