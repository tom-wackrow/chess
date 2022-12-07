package main

import (
	"net/http"

	routes "chess-website/routes"
)


func main() {
	http.HandleFunc("/dashboard", routes.Dashboard)
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/logout", routes.Logout)
	http.ListenAndServe(":80", nil)
}
