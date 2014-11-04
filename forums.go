package main

import "html/template"
import "io"
import "net/http"
import "os"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

import "github.com/daaku/go.httpgzip"
import "github.com/gorilla/mux"
import "github.com/gorilla/Schema"
import "github.com/gorilla/securecookie"
import "github.com/gorilla/sessions"

import "forums/model"

const (
	DATABASE_FILE = "forums.db"
)

type App struct {
	templates *template.Template
	db        *sql.DB
	sessions  *sessions.CookieStore
}

func newApp() *App {
	db, err := sql.Open("sqlite3", DATABASE_FILE)
	if err != nil {
		panic("error opening database")
	}

	templates := template.Must(template.ParseFiles(
		"templates/header.html",
		"templates/footer.html",
		"templates/index.html",
		"templates/forum.html",
		"templates/topic.html",
		"templates/addPost.html",
		"templates/addTopic.html",
	))

	sessionStore := sessions.NewCookieStore(securecookie.GenerateRandomKey(256))

	return &App{templates, db, sessionStore}
}

func (app *App) destroy() {
	app.db.Close()
}

func (app *App) addErrorFlashes(w http.ResponseWriter, r *http.Request, errs []error) {
	for _, err := range errs {
		app.addErrorFlash(w, r, err.Error())
	}
}

func (app *App) addErrorFlash(w http.ResponseWriter, r *http.Request, error string) {
    session, _ := app.sessions.Get(r, "forumSession")
    session.AddFlash(error)
    session.Save(r, w)
}

func (app *App) renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
    session, _ := app.sessions.Get(r, "forumSession")
    data["flashes"] = session.Flashes()
    session.Save(r, w)

	err := app.templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) handleIndex(w http.ResponseWriter, req *http.Request) {
	forums, err := model.FindForums(app.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := make(map[string]interface{})
	results["forums"] = forums

	app.renderTemplate(w, req, "index", results)
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
		return
	}

	results := make(map[string]interface{})
	results["forum"] = forum
	results["topics"] = topics

	app.renderTemplate(w, req, "forum", results)
}

func (app *App) handleTopic(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	topic, err := model.FindOneTopic(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := model.FindPosts(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := make(map[string]interface{})
	results["topic"] = topic
	results["posts"] = posts

	app.renderTemplate(w, req, "topic", results)
}

func (app *App) handleAddTopic(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	results := make(map[string]interface{})
	results["ForumId"] = id
	app.renderTemplate(w, req, "addTopic", results)
}

func (app *App) handleSaveTopic(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	topic := model.NewTopic()
    decoder := schema.NewDecoder()
    err := decoder.Decode(topic, req.PostForm)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }

    ok, errors := model.ValidateTopic(topic)
    if !ok {
    	app.addErrorFlashes(w, req, errors)
    	http.Redirect(w, req, "/forum/" + req.PostFormValue("ForumId") + "/add", 302)
    	return
    }

	err = model.SaveTopic(app.db, topic)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    }

	http.Redirect(w, req, "/forum/" + req.PostFormValue("ForumId"), 302)
}

func (app *App) handleAddPost(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	results := make(map[string]interface{})
	results["TopicId"] = id
	app.renderTemplate(w, req, "addPost", results)
}

func (app *App) handleSavePost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	post := model.NewPost()
    decoder := schema.NewDecoder()
    err := decoder.Decode(post, req.PostForm)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }

    ok, errors := model.ValidatePost(post)
    if !ok {
    	app.addErrorFlashes(w, req, errors)
    	http.Redirect(w, req, "/topic/" + req.PostFormValue("TopicId") + "/add", 302)
    	return
    }

	err = model.SavePost(app.db, post)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	f.HandleFunc("/{id:[0-9]+}/add", app.handleAddTopic).Methods("GET")
	f.HandleFunc("/{id:[0-9]+}/add", app.handleSaveTopic).Methods("POST")

	t := r.PathPrefix("/topic").Subrouter()
	t.HandleFunc("/{id:[0-9]+}", app.handleTopic)
	t.HandleFunc("/{id:[0-9]+}/add", app.handleAddPost).Methods("GET")
	t.HandleFunc("/{id:[0-9]+}/add", app.handleSavePost).Methods("POST")

	http.Handle("/", httpgzip.NewHandler(r))
	http.ListenAndServe(":8080", nil)

	app.destroy()
}
