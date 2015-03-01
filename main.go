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

	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
)

const (
	databaseFile = "forums.db"
	limitPosts   = 10
	limitTopics  = 10
)

func backup() error {
	src, err := os.Open(databaseFile)
	defer src.Close()
	if err != nil {
		return errors.New("could not open database to backup")
	}

	err = os.MkdirAll("backup", 0755)
	if err != nil {
		return errors.New("could not create backup")
	}

	destFile := path.Join("backup", databaseFile+".gz")
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

var listen = flag.String("listen", "localhost:8080", "host and port to listen on")

func main() {
	flag.Parse()

	err := backup()
	if err != nil {
		panic(err)
	}
	log.Println("backup complete")

	app := newApp()
	defer app.destroy()
	log.Println("database opened")

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/", app.handleIndex)

	f := r.PathPrefix("/forum/{id:[0-9]+}").Subrouter()
	f.HandleFunc("/", app.handleForum).Methods("GET")
	f.HandleFunc("page/{page:[0-9]+}", app.handleForum).Methods("GET")
	f.HandleFunc("/add", app.handleLoginRequired(app.handleAddTopic, "/forum")).Methods("GET")
	f.HandleFunc("/add", app.handleLoginRequired(app.handleSaveTopic, "/forum")).Methods("POST")

	t := r.PathPrefix("/topic/{id:[0-9]+}").Subrouter()
	t.HandleFunc("/", app.handleTopic).Methods("GET")
	t.HandleFunc("/page/{page:[0-9]+}", app.handleTopic).Methods("GET")
	t.HandleFunc("/add", app.handleLoginRequired(app.handleAddPost, "/topic")).Methods("GET")
	t.HandleFunc("/add", app.handleLoginRequired(app.handleSavePost, "/topic")).Methods("POST")
	t.HandleFunc("/delete", app.handleLoginRequired(app.handleDeletePost, "/topic")).Methods("POST")

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
