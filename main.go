package main

import (
	"net/http"

	auth "chess-website/auth"
	multiplayer "chess-website/multiplayer"
	routes "chess-website/routes"
)


func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", auth.RequireAuthenticatedUser(routes.Dashboard))
	http.HandleFunc("/chess", auth.RequireAuthenticatedUser(routes.LocalChess))
	http.HandleFunc("/vsbot", auth.RequireAuthenticatedUser(routes.VSBot))
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/logout", routes.Logout)
	http.HandleFunc("/stockfish", routes.Stockfish)
	// http.HandleFunc("/socket", chess.SocketHandler)
	http.HandleFunc("/multiplayer/create", multiplayer.MultiplayerCreate)
	http.HandleFunc("/multiplayer/play/", multiplayer.MultiplayerPlay)

	http.ListenAndServe(":80", nil)
}
