package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/Schema"
	"github.com/gorilla/mux"
	"github.com/mt2d2/forum/model"
)

func (app *app) handleAddPost(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	results := make(map[string]interface{})
	results["TopicId"] = id
	app.renderTemplate(w, req, "addPost", results)
}

func (app *app) handleSavePost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	post := model.NewPost()
	decoder := schema.NewDecoder()
	err := decoder.Decode(post, req.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := app.sessions.Get(req, "forumSession")
	if userID, ok := session.Values["user_id"].(int); ok {
		post.UserId = userID
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

	topic, err := model.FindOneTopic(app.db, strconv.Itoa(post.TopicId))

	http.Redirect(w, req, "/topic/"+req.PostFormValue("TopicId")+"/page/"+strconv.Itoa(numberOfTopicPages(topic)), http.StatusFound)
}

func (app *app) handleDeletePost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	session, _ := app.sessions.Get(req, "forumSession")
	if userID, ok := session.Values["user_id"].(int); ok {
		user, err := model.FindOneUserById(app.db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		post, err := model.FindOnePost(app.db, req.PostFormValue("PostId"))
		if err != nil {
			app.addErrorFlash(w, req, err)
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}

		if user.Id != post.User.Id {
			app.addErrorFlash(w, req, errors.New("You can only delete your own posts!"))
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}

		model.DeletePost(app.db, post.Id)
		http.Redirect(w, req, "/topic/"+req.PostFormValue("TopicId"), http.StatusFound)
	} else {
		app.addErrorFlash(w, req, errors.New("Must be logged in!"))
		http.Redirect(w, req, "/", http.StatusFound)
	}
}
