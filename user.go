package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mt2d2/forum/model"
)

func (app *app) handleRegister(w http.ResponseWriter, req *http.Request) {
	app.addBreadCrumb("/user/add", "Register")

	results := make(map[string]interface{})
	app.renderTemplate(w, req, "register", results)
}

func (app *app) saveRegister(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	user := model.NewUser()
	// manually grab password so we can convert to byte
	user.Username = req.PostFormValue("Username")
	user.Email = req.PostFormValue("Email")
	user.Password = []byte(req.PostFormValue("Password"))

	ok, errors := model.ValidateUser(app.db, user)
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

func (app *app) handleLogin(w http.ResponseWriter, req *http.Request) {
	app.addBreadCrumb("/user/login", "Login")
	results := make(map[string]interface{})
	results["Referer"] = req.Referer()
	app.renderTemplate(w, req, "login", results)
}

func (app *app) saveLogin(w http.ResponseWriter, req *http.Request) {
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

func (app *app) handleLogout(w http.ResponseWriter, req *http.Request) {
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

func (app *app) handleLoginRequired(nextHandler func(http.ResponseWriter, *http.Request), pathToRedirect string) func(http.ResponseWriter, *http.Request) {
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
