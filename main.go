package main

import (
	"net/http"

	auth "chess-website/auth"
	routes "chess-website/routes"
)


func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/dashboard", auth.RequireAuthenticatedUser(routes.Dashboard))
	http.HandleFunc("/chess", auth.RequireAuthenticatedUser(routes.LocalChess))
	http.HandleFunc("/vsbot", routes.VSBot)
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/logout", routes.Logout)
	http.HandleFunc("/stockfish", routes.Stockfish)
	http.ListenAndServe(":80", nil)
}
