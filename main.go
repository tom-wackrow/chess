package main

import (
	"net/http"

	routes "chess-website/routes"
)


func main() {
	http.HandleFunc("/dashboard", routes.Dashboard)
	http.ListenAndServe(":80", nil)
}
