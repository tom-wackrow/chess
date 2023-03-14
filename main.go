package main

import (
	"net/http"

	auth "chess-website/auth"
	multiplayer "chess-website/multiplayer"
	routes "chess-website/routes"
)


func main() {
	// serve static files such as javascript files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	// add routes that users can access
	http.HandleFunc("/", auth.RequireAuthenticatedUser(routes.Dashboard))
	http.HandleFunc("/chess", auth.RequireAuthenticatedUser(routes.LocalChess))
	http.HandleFunc("/vsbot", auth.RequireAuthenticatedUser(routes.VSBot))
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/logout", routes.Logout)
	http.HandleFunc("/stockfish", routes.Stockfish)
	// http.HandleFunc("/socket", chess.SocketHandler)
	http.HandleFunc("/multiplayer/create", auth.RequireAuthenticatedUser(multiplayer.MultiplayerCreate))
	http.HandleFunc("/multiplayer/play/", auth.RequireAuthenticatedUser(multiplayer.MultiplayerPlay))

	http.ListenAndServe(":80", nil)
}
