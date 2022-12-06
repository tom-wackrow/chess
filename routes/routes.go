package routes

import (
	"chess-website/auth"
	"html/template"
	"net/http"
)


func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	tmpl.Execute(w, "Hello World")
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if auth.AuthUser(w, r) {
			http.Redirect(w, r, "/dashboard", 303)
			return
		}
		http.Redirect(w, r, "/login", 303)
		return
	}
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/login.html")
	tmpl.Execute(w, nil)
}