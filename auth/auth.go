package auth

import (
	"chess-website/database"
	"net/http"
	"time"

	"crypto/sha256"

	"github.com/google/uuid"
)

// reference to database
var DB *database.Database

// empty entry for comparisons with data fetched from database
var	EMPTY_ENTRY database.UserInfoEntry = database.UserInfoEntry{}

// store user sessions in a dictionary
var sessions = map[string]session{}

func init() {
	// initialise reference to database when the module is loaded
	DB, _ = database.InitDB()
}

type session struct {
	username string
	expiry   time.Time
}

type credentials struct {
	Username, Password string
}

// helper function to determine if a session is expired yet (the default expiry time i set is 2 minutes)
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// retrieve credentials from POST request from user on "/login"
func GetCredentials(w http.ResponseWriter, r *http.Request) credentials {
	return credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
}

// takes a POST request and attempts to log in the user
func AuthUser(w http.ResponseWriter, r *http.Request) bool {
	creds := GetCredentials(w, r)

	// get expected password hash for username from database
	expectedEntry, err := DB.GetEntryByUsername(creds.Username)

	if err != nil {
		return false
	}

	// hash the user-entered password and compare it to the expected password
	actualPasswordHash := sha256.Sum256([]byte(creds.Password))
	if string(actualPasswordHash[:]) != expectedEntry.PasswordHash {
		return false
	}

	// create new session id and expiry time
	newSession := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// add session to store
	sessions[newSession] = session{
		username: creds.Username,
		expiry: expiresAt,
	}

	// set users browser session cookie
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: newSession,
		Expires: expiresAt,
	})

	// user was logged in correctly
	return true
}

// create account and add entry to database
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	creds := GetCredentials(w, r)

	// the only error GetEntryByUsername returns is a missing entry
	// so if an error is returned then we know that we can add the new
	// account to the database as there is no entry with that username
	if _, err := DB.GetEntryByUsername(creds.Username); err != nil {
		DB.Insert(database.UserInfoEntry{
			Username: creds.Username, 
			PasswordHash: creds.Password,
		})
	}
}

// log user out of website
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	sessionToken := cookie.Value

	delete(sessions, sessionToken) // remove session from store

	// remove session cookie from users browser
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: "",
		Expires: time.Now(),
	})
}

// middleware for requiring certain pages to only have logged in users
// accessing them
func RequireAuthenticatedUser(f http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")

		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}
		sessionToken := cookie.Value

		userSession, exists := sessions[sessionToken]


		if !exists { // user is not logged in / attempted to access page with false session token
			http.Redirect(w, r, "/login", 303)
			return
		}

		// if the session is expired remove it from the database and
		// reauthenticate user
		if userSession.isExpired() {
			delete(sessions, sessionToken)
			http.SetCookie(w, &http.Cookie{
				Name: "session",
				Value: "",
				Expires: time.Now(),
			})
			http.Redirect(w, r, "/login", 303)
			return
		}

		// create new session
		newSessionToken := uuid.NewString()
		ExpiresAt := time.Now().Add(120 * time.Second)

		sessions[newSessionToken] = session{
			username: sessions[sessionToken].username,
			expiry: ExpiresAt,
		}
		delete(sessions, sessionToken)

		http.SetCookie(w, &http.Cookie{
			Name: "session",
			Value: newSessionToken,
			Expires: ExpiresAt,
		})

		f(w, r) // load original page
	}
}

// whenever a user accesses a page, their session token will be refreshed
func RefreshSession(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")

		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}

		oldSessionToken := cookie.Value
		newSessionToken := uuid.NewString()
		ExpiresAt := time.Now().Add(120 * time.Second)

		sessions[newSessionToken] = session{
			username: sessions[oldSessionToken].username,
			expiry: ExpiresAt,
		}
		
		// delete old token
		delete(sessions, oldSessionToken)

		// set browser cookie to new token
		http.SetCookie(w, &http.Cookie{
			Name: "session",
			Value: newSessionToken,
			Expires: ExpiresAt,
		})

		f(w, r)
	}
}