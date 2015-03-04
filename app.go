package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mt2d2/forum/model"
	"github.com/russross/blackfriday"
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
		log.Panicln(err)
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
