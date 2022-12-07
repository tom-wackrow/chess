package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

var users = map[string]string{
	"admin": "admin",
	"test":  "test",
}

var sessions = map[string]session{}

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

	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Username {
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
	if _, ok := users[creds.Username]; !ok {
		users[creds.Username] = creds.Password
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

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
		cookie, _ := r.Cookie("session")

		sessionToken := cookie.Value

		userSession, exists := sessions[sessionToken]

		if !exists {
			http.Redirect(w, r, "/login", 303)
			return
		}

		if userSession.isExpired() {
			delete(sessions, sessionToken)
			http.Redirect(w, r, "/login", 303)
			return
		}

		f(w, r)
	}
}