package main

import (
	"net/http"

	routes "chess-website/routes"
)


func main() {
	http.HandleFunc("/dashboard", routes.Dashboard)
	http.HandleFunc("/login", routes.Login)
	http.ListenAndServe(":80", nil)
}
