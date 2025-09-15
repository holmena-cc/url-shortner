package server

import (
	"fmt"
	"github.com/resend/resend-go/v2"
	"html/template"
	"net/http"
	"os"
)

type ContactPageData struct {
	IsLoggedIn bool
}

func (s *Server) contactHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/contact.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	_, ok := r.Context().Value(userIDKey).(int32)
	data := ContactPageData{
		IsLoggedIn: ok,
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		fmt.Println("failed to render err:", err)
	}
}

func (s *Server) contactFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	userName := r.FormValue("name")
	userEmail := r.FormValue("email")
	message := r.FormValue("message")
	if userName == "" || userEmail == "" || message == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	// send the email
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		fmt.Println("Missing RESEND_API_KEY")
		http.Error(w, "Server misconfigured", http.StatusInternalServerError)
		return
	}
	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"badisstiti11@gmail.com"},
		Subject: "New message from ShortyLink",
		Html: fmt.Sprintf(`<p><strong>Name:</strong> %s</p>
	                   <p><strong>Email:</strong> %s</p>
	                   <p><strong>Message:</strong><br>%s</p>`, userName, userEmail, message),
	}
	_, err = client.Emails.Send(params)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		fmt.Println("Email send error:", err)
		return
	}
	http.Redirect(w, r, "/thankyou", http.StatusSeeOther)

}
