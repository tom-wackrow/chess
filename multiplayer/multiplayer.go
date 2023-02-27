package multiplayer

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	chess "chess-website/chess"
)



func MultiplayerCreate(w http.ResponseWriter, r *http.Request) { // /multiplayer/create
	game := chess.CreateGame()
	chess.GameList = append(chess.GameList, game)

	http.Redirect(w, r, fmt.Sprintf("/multiplayer/play/%v", game.ID), 303) // redirect to page for game' that was created
}

func MultiplayerPlay(w http.ResponseWriter, r *http.Request) { // /multiplayer/play/:id
	id, _ := strconv.Atoi(r.URL.Path[len("/multiplayer/play/"):])
	if chess.GameExists(id) {
		// game := chess.GetGameByID(id)
		// stuff to play game here

		// go HandleGame(game)
		tmpl, _ := template.ParseFiles("templates/base.html", "templates/chess/multiplayer.html")
		tmpl.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/multiplayer/create", 303)
	}
}

func MultiplayerSocketHandler(game chess.ChessGame) {
	for game.WhitePlayer.Conn == nil {} // when game is created, user that created the game is set as black
	// when white websocket connection is no longer empty, the function will continue

	
}