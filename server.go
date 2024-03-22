package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserServer struct {
	store UserStore
	http.Handler
}

type User struct {
	ID        int          `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	Squeaks   []SqueakPost `json:"squeaks"`
	CreatedAt time.Time    `json:"createdAt"`
}

// Squeaks are Gopher's variant of tweets
type SqueakPost struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUserServer(store UserStore) *UserServer {
	u := new(UserServer)

	u.store = store

	router := http.NewServeMux()
	router.Handle("/userbase", http.HandlerFunc(u.userbaseHandler))
	router.Handle("GET /users/{name}", http.HandlerFunc(u.showSqueaks))
	router.Handle("POST /users/{name}", http.HandlerFunc(u.saveSqueak))
	router.Handle("/register", http.HandlerFunc(u.registerUser))
	router.Handle("/login", http.HandlerFunc(u.loginUser))

	u.Handler = router

	return u
}

type UserStore interface {
	GetUserSqueaks(name string) ([]SqueakPost, error)
	PostSqueak(name, squeak string) (int, error)
	GetUserbase() ([]User, error)
	CreateUser(name, email, password string) (int, error)
	GetUserByUsername(username string) (*User, error)
}

const jsonContentType = "application/json"

func (u *UserServer) userbaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	userbase, err := u.store.GetUserbase()
	if err != nil {
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(userbase); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (u *UserServer) showSqueaks(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("name")
	squeaks, err := u.store.GetUserSqueaks(user)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)

	if err := json.NewEncoder(w).Encode(squeaks); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (u *UserServer) saveSqueak(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("name")
	var squeak SqueakPost

	if err := json.NewDecoder(r.Body).Decode(&squeak); err != nil {
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	_, err := u.store.PostSqueak(username, squeak.Text)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
}

func (u *UserServer) registerUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	existingUser, err := u.store.GetUserByUsername(user.Username)
	if err == nil {
		if existingUser.Username == user.Username {
			http.Error(w, "username already taken", http.StatusConflict)
		}
		if existingUser.Email == user.Email {
			http.Error(w, "email already taken", http.StatusConflict)
		}
		return
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed hashing the password", http.StatusInternalServerError)
		return
	}

	_, err = u.store.CreateUser(user.Username, user.Email, string(encpw))
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating user: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (u *UserServer) loginUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	username := string(user.Username)
	password := string(user.Password)

	success := u.verifyCredentials(username, password)
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		http.Error(w, "failed to create token", http.StatusBadRequest)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		MaxAge:   int(time.Hour * 24 * 30),
		Secure:   false,
		HttpOnly: true,
		SameSite: 2,
	})

	w.WriteHeader(http.StatusAccepted)
}

func (u *UserServer) verifyCredentials(username, password string) bool {
	user, err := u.store.GetUserByUsername(username)
	if err != nil {
		log.Println("verifyCredentials: %w", err)
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err)
		return false
	}

	return true
}
