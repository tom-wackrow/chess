package chess

type ChessGame struct {
	ID            int    `json:"id"`
	FEN           string `json:"fen"`
	CurrentPlayer string `json:"currentPlayer"`
}

var GameList = []ChessGame{}

func CreateGame() ChessGame {
	game := ChessGame{
		ID:            len(GameList),
		FEN:           "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		CurrentPlayer: "w",
	}

	return game
}

func GetGameByID(id int) ChessGame {
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