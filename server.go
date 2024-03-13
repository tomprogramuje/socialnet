package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserServer struct {
	store UserStore
	http.Handler
}

type User struct {
	ID       int
	Username string
	Email    string
	Password string
	Squeaks  []string
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
	// Squeaks are Gopher's variant of tweets
	GetUserSqueaks(name string) ([]string, error)
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
		log.Println(err)
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
	user := r.PathValue("name")
	var payload User

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	squeak := string(payload.Squeaks[0])

	u.store.PostSqueak(user, squeak)
	w.WriteHeader(http.StatusAccepted)
}

func (u *UserServer) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload User

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	username := string(payload.Username)
	email := string(payload.Email)
	password := string(payload.Password)

	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed hashing the password", http.StatusInternalServerError)
	}

	u.store.CreateUser(username, email, string(encpw))
	w.WriteHeader(http.StatusAccepted)
}

func (u *UserServer) loginUser(w http.ResponseWriter, r *http.Request) {
	var payload User

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	username := string(payload.Username)
	password := string(payload.Password)

	if !u.verifyCredentials(username, password) {
		w.WriteHeader(http.StatusBadRequest) // check if it's ok
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (u *UserServer) verifyCredentials(username, password string) bool {
	user, _ := u.store.GetUserByUsername(username)
	bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))

	return true
}
