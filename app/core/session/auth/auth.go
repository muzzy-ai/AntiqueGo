package auth

import (
	"AntiqueGo/app/models"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var store *sessions.CookieStore
var sessionUser = "user-session"

func SetSessionStore(sessionKey string) {
	store = sessions.NewCookieStore([]byte(sessionKey))
}

func GetSessionUser(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, sessionUser)
}

func IsLoggedIn(r *http.Request) bool {
	session, _ := store.Get(r, sessionUser)
	return session.Values["id"] != nil
}

func ComparePassword(password string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func MakePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func CurrentUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) *models.User {
	if !IsLoggedIn(r) {
		return nil
	}

	session, _ := store.Get(r, sessionUser)

	userModel := models.User{}
	user, err := userModel.FindByID(db, session.Values["id"].(string))
	if err != nil {
		session.Values["id"] = nil
		session.Save(r, w)
		return nil
	}

	return user
}