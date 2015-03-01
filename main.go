package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"

	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/mt2d2/forum/model"
)

const (
	databaseFile = "forums.db"
	limitPosts   = 10
	limitTopics  = 10
)

func convertToMarkdown(markdown string) template.HTML {
	unsafe := blackfriday.MarkdownCommon([]byte(markdown))

	policy := bluemonday.UGCPolicy()
	policy.AllowElements("video", "audio", "source")
	policy.AllowAttrs("controls").OnElements("video", "audio")
	policy.AllowAttrs("src").Matching(regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`)).Globally()

	html := policy.SanitizeBytes(unsafe)
	return template.HTML(html)
}

type breadCrumb struct{ URL, Title string }

type app struct {
	templates   *template.Template
	db          *sql.DB
	sessions    *sessions.CookieStore
	breadCrumbs []breadCrumb
}

func newApp() *app {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		panic("error opening database")
	}

	templates, err := template.New("").Funcs(template.FuncMap{"markDown": convertToMarkdown}).ParseFiles(
		"templates/header.html",
		"templates/footer.html",
		"templates/index.html",
		"templates/forum.html",
		"templates/topic.html",
		"templates/addPost.html",
		"templates/addTopic.html",
		"templates/register.html",
		"templates/login.html",
	)

	if err != nil {
		panic(err)
	}

	sessionStore := sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

	breadCrumbs := make([]breadCrumb, 0, 1)
	breadCrumbs = append(breadCrumbs, breadCrumb{"/", "Index"})
	return &app{templates, db, sessionStore, breadCrumbs}
}

func (app *app) destroy() {
	app.db.Close()
}

func (app *app) addBreadCrumb(url, title string) {
	app.breadCrumbs = append(app.breadCrumbs, breadCrumb{url, title})
}

func (app *app) useBreadCrumbs() *[]breadCrumb {
	ret := app.breadCrumbs
	app.breadCrumbs = app.breadCrumbs[:1]
	return &ret
}

func (app *app) addErrorFlashes(w http.ResponseWriter, r *http.Request, errs []error) {
	for _, err := range errs {
		app.addErrorFlash(w, r, err)
	}
}

func (app *app) addErrorFlash(w http.ResponseWriter, r *http.Request, error error) {
	app.addFlash(w, r, error.Error(), "error")
}

func (app *app) addSuccessFlash(w http.ResponseWriter, r *http.Request, str string) {
	app.addFlash(w, r, str, "success")
}

func (app *app) addFlash(w http.ResponseWriter, r *http.Request, content interface{}, key string) {
	session, _ := app.sessions.Get(r, "forumSession")
	session.AddFlash(content, key)
	session.Save(r, w)
}

func (app *app) renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	session, _ := app.sessions.Get(r, "forumSession")

	data["breadCrumbs"] = app.useBreadCrumbs()
	data["errorFlashes"] = session.Flashes("error")
	data["successFlashes"] = session.Flashes("success")

	if userID, ok := session.Values["user_id"].(int); ok {
		user, err := model.FindOneUserById(app.db, userID)
		if err == nil {
			data["user"] = user
		}
	}

	session.Save(r, w)

	err := app.templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) handleIndex(w http.ResponseWriter, req *http.Request) {
	forums, err := model.FindForums(app.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := make(map[string]interface{})
	results["forums"] = forums

	app.renderTemplate(w, req, "index", results)
}

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
