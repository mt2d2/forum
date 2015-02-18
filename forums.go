package main

import "errors"
import "html/template"
import "io"
import "net/http"
import "os"
import "regexp"
import "strings"
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

import "github.com/microcosm-cc/bluemonday"
import "github.com/russross/blackfriday"
import "github.com/daaku/go.httpgzip"
import "github.com/gorilla/mux"
import "github.com/gorilla/Schema"
import "github.com/gorilla/securecookie"
import "github.com/gorilla/sessions"

import "github.com/mt2d2/forum/model"

const (
	DATABASE_FILE = "forums.db"
)

type App struct {
	templates *template.Template
	db        *sql.DB
	sessions  *sessions.CookieStore
}

func convertToMarkdown(markdown string) template.HTML {
	unsafe := blackfriday.MarkdownCommon([]byte(markdown))

	policy := bluemonday.UGCPolicy()
	policy.AllowElements("video", "audio", "source")
	policy.AllowAttrs("controls").On("video", "audio")
	policy.AllowAttrs("src").Matching(regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`)).Globally()

	html := policy.SanitizeBytes(unsafe)
	return template.HTML(html)
}

func newApp() *App {
	db, err := sql.Open("sqlite3", DATABASE_FILE)
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

	return &App{templates, db, sessionStore}
}

func (app *App) destroy() {
	app.db.Close()
}

func (app *App) addErrorFlashes(w http.ResponseWriter, r *http.Request, errs []error) {
	for _, err := range errs {
		app.addErrorFlash(w, r, err)
	}
}

func (app *App) addErrorFlash(w http.ResponseWriter, r *http.Request, error error) {
	app.addFlash(w, r, error.Error(), "error")
}

func (app *App) addSuccessFlash(w http.ResponseWriter, r *http.Request, str string) {
	app.addFlash(w, r, str, "success")
}

func (app *App) addFlash(w http.ResponseWriter, r *http.Request, content interface{}, key string) {
	session, _ := app.sessions.Get(r, "forumSession")
	session.AddFlash(content, key)
	session.Save(r, w)
}

func (app *App) renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	session, _ := app.sessions.Get(r, "forumSession")

	data["errorFlashes"] = session.Flashes("error")
	data["successFlashes"] = session.Flashes("success")

	if userId, ok := session.Values["user_id"].(int); ok {
		user, err := model.FindOneUserById(app.db, userId)
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

	ok, errors := model.ValidateTopic(app.db, topic)
	if !ok {
		app.addErrorFlashes(w, req, errors)
		http.Redirect(w, req, "/forum/"+req.PostFormValue("ForumId")+"/add", http.StatusFound)
		return
	}

	err = model.SaveTopic(app.db, topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, req, "/forum/"+req.PostFormValue("ForumId"), http.StatusFound)
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

	session, _ := app.sessions.Get(req, "forumSession")
	if userId, ok := session.Values["user_id"].(int); ok {
		post.UserId = userId
	}

	ok, errors := model.ValidatePost(app.db, post)
	if !ok {
		app.addErrorFlashes(w, req, errors)
		http.Redirect(w, req, "/topic/"+req.PostFormValue("TopicId")+"/add", http.StatusFound)
		return
	}

	err = model.SavePost(app.db, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/topic/"+req.PostFormValue("TopicId"), http.StatusFound)
}

func (app *App) handleRegister(w http.ResponseWriter, req *http.Request) {
	results := make(map[string]interface{})
	app.renderTemplate(w, req, "register", results)
}

func (app *App) saveRegister(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	user := model.NewUser()
	// manually grab password so we can convert to byte
	user.Username = req.PostFormValue("Username")
	user.Email = req.PostFormValue("Email")
	user.Password = []byte(req.PostFormValue("Password"))

	ok, errors := model.ValidateUser(user)
	if !ok {
		app.addErrorFlashes(w, req, errors)
		http.Redirect(w, req, "/user/add", http.StatusFound)
		return
	}

	err := user.HashPassword()
	if err != nil {
		app.addErrorFlash(w, req, err)
		http.Redirect(w, req, "/user/add", http.StatusFound)
		return
	}

	err = model.SaveUser(app.db, user)
	if err != nil {
		app.addErrorFlash(w, req, err)
		http.Redirect(w, req, "/user/add", http.StatusFound)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func (app *App) handleLogin(w http.ResponseWriter, req *http.Request) {
	results := make(map[string]interface{})
	results["Referer"] = req.Referer()
	app.renderTemplate(w, req, "login", results)
}

func (app *App) saveLogin(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	username := req.PostFormValue("Username")
	password := []byte(req.PostFormValue("Password"))

	if username == "" || len(password) == 0 {
		app.addErrorFlash(w, req, errors.New("Enter a username and password."))
		http.Redirect(w, req, "/user/login", http.StatusFound)
		return
	}

	invalidUserOrPassword := errors.New("Invalid username or password.")

	user, err := model.FindOneUser(app.db, username)
	if err != nil {
		app.addErrorFlash(w, req, invalidUserOrPassword)
		http.Redirect(w, req, "/user/login", http.StatusFound)
		return
	}

	err = user.CompareHashAndPassword(&password)
	if err != nil {
		app.addErrorFlash(w, req, invalidUserOrPassword)
		http.Redirect(w, req, "/user/login", http.StatusFound)
		return
	}

	session, _ := app.sessions.Get(req, "forumSession")
	session.Values["user_id"] = user.Id
	session.Save(req, w)

	app.addSuccessFlash(w, req, "Successfully logged in!")

	toRedirect := req.PostFormValue("Referer")
	if toRedirect == "" || strings.HasSuffix(toRedirect, "login") {
		toRedirect = "/"
	}

	http.Redirect(w, req, toRedirect, http.StatusFound)
}

func (app *App) handleLogout(w http.ResponseWriter, req *http.Request) {
	session, _ := app.sessions.Get(req, "forumSession")
	delete(session.Values, "user_id")
	session.Save(req, w)

	app.addSuccessFlash(w, req, "Successfully logged out.")

	toRedirect := req.Referer()
	if toRedirect == "" {
		toRedirect = "/"
	}

	http.Redirect(w, req, toRedirect, http.StatusFound)
}

func (app *App) handleLoginRequired(nextHandler func(http.ResponseWriter, *http.Request), pathToRedirect string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := app.sessions.Get(req, "forumSession")
		if _, ok := session.Values["user_id"]; !ok {
			newPath := pathToRedirect
			if id, ok := mux.Vars(req)["id"]; ok {
				newPath += "/" + id
			}

			app.addErrorFlash(w, req, errors.New("Must be logged in!"))
			http.Redirect(w, req, newPath, http.StatusFound)
			return
		}

		nextHandler(w, req)
	}
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
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/", app.handleIndex)

	f := r.PathPrefix("/forum").Subrouter()
	f.HandleFunc("/{id:[0-9]+}", app.handleForum)
	f.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleAddTopic, "/forum")).Methods("GET")
	f.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleSaveTopic, "/forum")).Methods("POST")

	t := r.PathPrefix("/topic").Subrouter()
	t.HandleFunc("/{id:[0-9]+}", app.handleTopic)
	t.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleAddPost, "/topic")).Methods("GET")
	t.HandleFunc("/{id:[0-9]+}/add", app.handleLoginRequired(app.handleSavePost, "/topic")).Methods("POST")

	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/add", app.handleRegister).Methods("GET")
	u.HandleFunc("/add", app.saveRegister).Methods("POST")
	u.HandleFunc("/login", app.handleLogin).Methods("GET")
	u.HandleFunc("/login", app.saveLogin).Methods("POST")
	u.HandleFunc("/logout", app.handleLogout)

	http.Handle("/", httpgzip.NewHandler(r))
	http.ListenAndServe(":8080", nil)

	app.destroy()
}
