package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
)

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

		user, dbErr := s.db.DB().GetUserByEmail(context.Background(), email)
		if dbErr != nil {
			tmpl.ExecuteTemplate(w, "base", LoginPageData{
				Error: "❌ Incorrect email or password",
			})
			return
		}

		if err := CheckPassword(user.PasswordHash, password); err != nil {
			http.Error(w, "❌ Incorrect email or password", http.StatusUnauthorized)
			return
		}

		// Success
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
