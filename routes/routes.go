package routes

import (
	"chess-website/auth"
	"fmt"
	"html/template"
	"net/http"

	"github.com/freeeve/uci"
)

// access "/dashboarD"
func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/dashboard.html")
	tmpl.Execute(w, "Hello World")
}

// access "/login"
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { // if request method is "POST" attempt to log user in
		if auth.AuthUser(w, r) {
			http.Redirect(w, r, "/dashboard", 303)
			return
		}
		http.Redirect(w, r, "/login", 303)
		return
	}
	// show login page to user
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/login.html")
	tmpl.Execute(w, nil)
}

// access "/register"
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method =="POST" { //if request method is "POST" attempt to create account
		creds := auth.GetCredentials(w, r)
		pwd2 := r.FormValue("password2")

		if _, err := auth.DB.GetEntryByUsername(creds.Username); err != nil {
			if pwd2 == creds.Password { // validate that password1 and password2 are equal
				auth.RegisterUser(w, r) // register user
				auth.AuthUser(w, r)
				http.Redirect(w, r, "/dashboard", 303) // redirect user to "/dashboard"
				return
			}
		}
		return
	}

	tmpl, _ := template.ParseFiles("templates/base.html", "templates/register.html")
	tmpl.Execute(w, nil)
}

// logout user
func Logout(w http.ResponseWriter, r *http.Request) {
	auth.LogoutUser(w, r)
	http.Redirect(w, r, "/login", 303)
}

// 'pass & play' chess
func LocalChess(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/chess/localChess.html")
	tmpl.Execute(w, nil)
}

// play chess against bot
func VSBot(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/base.html", "templates/chess/VSBot.html")
	tmpl.Execute(w, nil)
}

// API for playing against bot
func Stockfish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}


	fen := r.FormValue("fen")

	eng, err := uci.NewEngine("/stockfish.exe")
	if err != nil {
		fmt.Print(err)
		return
	}

	eng.SetOptions(uci.Options{
		Hash: 128,
		Ponder: false,
		OwnBook: true,
		MultiPV: 4,
	})

	eng.SetFEN(fen) // set engine board to one recieved from request

	resultOptions := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
	result, _ := eng.GoDepth(10, resultOptions)

	w.Write([]byte(result.BestMove)) // return best move
}
