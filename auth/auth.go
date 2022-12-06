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