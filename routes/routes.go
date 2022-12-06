package routes

import (
	"html/template"
	"net/http"
)


func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	tmpl.Execute(w, "Hello World")
}
