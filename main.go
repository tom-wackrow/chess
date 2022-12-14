package main

import (
	"net/http"

	auth "chess-website/auth"
	routes "chess-website/routes"
	sockets "chess-website/sockets"
)


func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/dashboard", auth.RequireAuthenticatedUser(routes.Dashboard))
	http.HandleFunc("/chess/local", auth.RequireAuthenticatedUser(routes.LocalChess))
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/logout", routes.Logout)

	http.HandleFunc("/chess/stockfish", sockets.StockfishSocketHandler)
	http.HandleFunc("/chess/vsbot", routes.ChessVSBot)
	http.ListenAndServe(":80", nil)
}
