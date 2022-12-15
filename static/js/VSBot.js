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

    var body = new FormData();
    body.append("fen", game.fen());
    let response = fetch("/stockfish", {
        method: "POST",
        body: body,
    }).then(async (response) => {
        const data = await response.text();
        console.log(data);
        console.log(data.substring(0, 2));
        console.log(data.substring(2,));
        game.move({
        from: data.substring(0, 2),
        to: data.substring(2,),
        promotion: "q",
    });

    board.position(game.fen())
    });

    
}
