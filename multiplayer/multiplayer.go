package multiplayer

import (
	"fmt"
	"net/http"
	"strconv"

	chess "chess-website/chess"
)



func MultiplayerCreate(w http.ResponseWriter, r *http.Request) { // /multiplayer/create
	game := chess.CreateGame()
	chess.GameList = append(chess.GameList, game)

	http.Redirect(w, r, fmt.Sprintf("/multiplayer/play/%v", game.ID), 303) // redirect to page for game that was created
}

func MultiplayerPlay(w http.ResponseWriter, r *http.Request) { // /multiplayer/play/:id
	id, _ := strconv.Atoi(r.URL.Path[len("/multiplayer/play/"):])
	if chess.GameExists(id) {
		game  := chess.GetGameByID(id)
		w.Write([]byte(strconv.Itoa(game.ID)))
		
	} else {
		http.Redirect(w, r, "/multiplayer/create", 303)
	}
}
