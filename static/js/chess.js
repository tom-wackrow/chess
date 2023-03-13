var config = {
    draggable: true,
    position: "start",
    onDragStart: onDragStart,
    onDrop: onDrop,
}

var board = Chessboard("board", config)
var game = new Chess()


function onDragStart (source, piece, position, orientation) {
    // if the game is over, do not pick up any pieces
    if (game.game_over()) return false

    // only pick up pieces for the side whose turn it is
    if ((game.turn() === 'w' && piece.search(/^b/) !== -1) ||
        (game.turn() === 'b' && piece.search(/^w/) !== -1)) {
    return false
    }
}

function onDrop(source, target) {
    // move will be null if move is not legal
    var move = game.move({
        from: source,
        to: target,
        promotion: "q", // for now always promote to queen
    })
    
    if (move === null) return "snapback" // if move is not legal, snap back to previous board state

    board.position(game.fen())
}
