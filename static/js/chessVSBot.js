

window.addEventListener("load", function () {
    var socket = new WebSocket("ws://localhost:80/chess/stockfish")

    socket.addEventListener("open", (event) => {
        socket.send("hello world")
    })

    socket.addEventListener("message", (event) => {
        console.log(event.data)
    })
})