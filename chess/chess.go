package chess

import "github.com/gorilla/websocket"

// for storing currently ongoing games
type ChessGame struct {
	ID            int    `json:"id"`
	FEN           string `json:"fen"`
	CurrentPlayer string `json:"currentPlayer"`
	WhitePlayer     Player `json:"-"`
	BlackPlayer     Player `json:"-"`
}

type Player struct {
	Username string
	Conn *websocket.Conn
}

// list of all occuring games
var GameList = []ChessGame{}

// initialise a new chess game
func CreateGame() ChessGame {
	game := ChessGame{
		ID:            len(GameList),
		FEN:           "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		CurrentPlayer: "w",
		WhitePlayer: Player{
			Username: "",
			Conn: nil,
		},
		BlackPlayer: Player{
			Username: "",
			Conn: nil,
		},
	}

	return game
}

// return a game from GameList using its id
func GetGameByID(id int) (ChessGame) {
	for _, game := range GameList {
		if game.ID == id {
			return game
		}
	}
	return ChessGame{}
}

func GameExists(id int) bool {
	for _, game := range GameList {
		if game.ID == id {
			return true
		}
	}
	return false
}