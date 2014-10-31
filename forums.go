package main

import "html/template"
import "io"
import "net/http"
import "os"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

import "github.com/gorilla/mux"
import "github.com/gorilla/Schema"

import "forums/model"

const (
	DATABASE_FILE = "forums.db"
)

type App struct {
	templates *template.Template
	db        *sql.DB
}

func newApp() *App {
	db, err := sql.Open("sqlite3", DATABASE_FILE)
	if err != nil {
		panic("error opening database")
	}

	return &App{template.Must(template.ParseFiles(
		"templates/header.html",
		"templates/footer.html",
		"templates/index.html",
		"templates/forum.html",
		"templates/topic.html",
		"templates/addPost.html",
	)), db}
}

func (app *App) destroy() {
	app.db.Close()
}

func (app *App) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := app.templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) handleIndex(w http.ResponseWriter, req *http.Request) {
	forums, err := model.FindForums(app.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	app.renderTemplate(w, "index", forums)
}

func (app *App) handleForum(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	forum, err := model.FindOneForum(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	topics, err := model.FindTopics(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	results := make(map[string]interface{})
	results["forum"] = forum
	results["topics"] = topics

	app.renderTemplate(w, "forum", results)
}

func (app *App) handleTopic(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	topic, err := model.FindOneTopic(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	posts, err := model.FindPosts(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	results := make(map[string]interface{})
	results["topic"] = topic
	results["posts"] = posts

	app.renderTemplate(w, "topic", results)
}

func (app *App) handleAddPost(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	results := make(map[string]interface{})
	results["TopicId"] = id
	app.renderTemplate(w, "addPost", results)
}

func (app *App) handleSavePost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	post := model.NewPost()
    decoder := schema.NewDecoder()
    err := decoder.Decode(post, req.PostForm)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    }

	err = model.SavePost(app.db, post)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    }

	http.Redirect(w, req, "/topic/" + req.PostFormValue("TopicId"), 302)
}

func backup() {
	src, err := os.Open(DATABASE_FILE)
	defer src.Close()
	if err != nil {
		panic("could not open database to backup")
	}

	dest, err := os.Create("backup/" + DATABASE_FILE)
	defer dest.Close()
	if err != nil {
		panic("could not open backup/" + DATABASE_FILE)
	}

	io.Copy(dest, src)
}

func main() {
	backup()

	app := newApp()

	r := mux.NewRouter()
	r.HandleFunc("/", app.handleIndex)

	f := r.PathPrefix("/forum").Subrouter()
	f.HandleFunc("/{id:[0-9]+}", app.handleForum)

	t := r.PathPrefix("/topic").Subrouter()
	t.HandleFunc("/{id:[0-9]+}", app.handleTopic)
	t.HandleFunc("/{id:[0-9]+}/add", app.handleAddPost).Methods("GET")
	t.HandleFunc("/{id:[0-9]+}/add", app.handleSavePost).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

	app.destroy()
}
