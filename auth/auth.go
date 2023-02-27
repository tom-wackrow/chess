package auth

import (
	"chess-website/database"
	"net/http"
	"time"

	"crypto/sha256"

	"github.com/google/uuid"
)

var DB *database.Database
var	EMPTY_ENTRY database.UserInfoEntry = database.UserInfoEntry{}

var sessions = map[string]session{}

func init() {
	DB, _ = database.InitDB()
}

type session struct {
	username string
	expiry   time.Time
}

type credentials struct {
	Username, Password string
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func GetCredentials(w http.ResponseWriter, r *http.Request) credentials{
	return credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
}


func AuthUser(w http.ResponseWriter, r *http.Request) bool {
	creds := GetCredentials(w, r)

	expectedEntry, err := DB.GetEntryByUsername(creds.Username)

	if err != nil {
		return false
	}

	actualPasswordHash := sha256.Sum256([]byte(creds.Password))
	if string(actualPasswordHash[:]) != expectedEntry.PasswordHash {
		return false
	}

	
	newSession := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[newSession] = session{
		username: creds.Username,
		expiry: expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: newSession,
		Expires: expiresAt,
	})
	return true
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	creds := GetCredentials(w, r)
	if _, err := DB.GetEntryByUsername(creds.Username); err != nil {
		DB.Insert(database.UserInfoEntry{
			Username: creds.Username, 
			PasswordHash: creds.Password,
		})
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}
	sessionToken := cookie.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: "",
		Expires: time.Now(),
	})
}

func RequireAuthenticatedUser(f http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")

		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}
		sessionToken := cookie.Value

		userSession, exists := sessions[sessionToken]

		if !exists {
			http.Redirect(w, r, "/login", 303)
			return
		}

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

		f(w, r)
	}
}

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
		delete(sessions, oldSessionToken)

		http.SetCookie(w, &http.Cookie{
			Name: "session",
			Value: newSessionToken,
			Expires: ExpiresAt,
		})

		f(w, r)
	}
}