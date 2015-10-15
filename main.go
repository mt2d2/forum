package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/GeertJohan/go.rice"
	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
)

const (
	limitPosts  = 10
	limitTopics = 10
)

var listen = flag.String("listen", "localhost:8080", "host and port to listen on")
var db = flag.String("db", "forum.db", "sqlite3 database file")

func backup() error {
	src, err := os.Open(*db)
	defer src.Close()
	if err != nil {
		return errors.New("could not open database to backup")
	}

	backupPath := path.Join(filepath.Dir(*db), "backup")
	err = os.MkdirAll(backupPath, 0755)
	if err != nil {
		return errors.New("could not create backup")
	}

	destFile := path.Join(backupPath, filepath.Base(*db)+".gz")
	dest, err := os.Create(destFile)
	defer dest.Close()
	if err != nil {
		return err
	}

	gzipWriter := gzip.NewWriter(dest)
	_, err = io.Copy(gzipWriter, src)
	if err != nil {
		return err
	}

	return gzipWriter.Close()
}

func main() {
	flag.Parse()

	err := backup()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("backup complete")

	app := newApp()
	defer app.destroy()
	log.Println("database opened")

	r := mux.NewRouter()
	staticBox := rice.MustFindBox("static").HTTPBox()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox())))

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

	log.Printf("Serving on %s\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
