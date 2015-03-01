package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
)

const (
	databaseFile = "forums.db"
	limitPosts   = 10
	limitTopics  = 10
)

func backup() {
	src, err := os.Open(databaseFile)
	defer src.Close()
	if err != nil {
		panic("could not open database to backup")
	}

	dest, err := os.Create("backup/" + databaseFile)
	defer dest.Close()
	if err != nil {
		panic("could not open backup/" + databaseFile)
	}

	io.Copy(dest, src)
}

func main() {
	backup()

	app := newApp()
	defer app.destroy()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/", app.handleIndex)

	f := r.PathPrefix("/forum").Subrouter()
	f.HandleFunc("/{id:[0-9]+}", app.handleForum).Methods("GET")
	f.HandleFunc("/{id:[0-9]+}/page/{page:[0-9]+}", app.handleForum).Methods("GET")
	f.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleAddTopic, "/forum")).Methods("GET")
	f.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleSaveTopic, "/forum")).Methods("POST")

	t := r.PathPrefix("/topic").Subrouter()
	t.HandleFunc("/{id:[0-9]+}", app.handleTopic).Methods("GET")
	t.HandleFunc("/{id:[0-9]+}/page/{page:[0-9]+}", app.handleTopic).Methods("GET")
	t.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleAddPost, "/topic")).Methods("GET")
	t.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleSavePost, "/topic")).Methods("POST")
	t.HandleFunc("/{id:[0-9]+}/delete", app.handleLoginRequired(app.handleDeletePost, "/topic")).Methods("POST")

	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/add", app.handleRegister).Methods("GET")
	u.HandleFunc("/add", app.saveRegister).Methods("POST")
	u.HandleFunc("/login", app.handleLogin).Methods("GET")
	u.HandleFunc("/login", app.saveLogin).Methods("POST")
	u.HandleFunc("/logout", app.handleLogout)

	http.Handle("/", httpgzip.NewHandler(r))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
