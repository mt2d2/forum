package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/Schema"
	"github.com/gorilla/mux"
	"github.com/mt2d2/forum/model"
)

func numberOfTopicPages(topic *model.Topic) int {
	return int(math.Ceil(float64(topic.PostCount) / float64(limitPosts)))
}

func (app *app) handleTopic(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	pageOffset := 0
	if page, ok := vars["page"]; ok {
		if val, err := strconv.Atoi(page); err == nil {
			pageOffset = val - 1
		}
	}

	topic, err := model.FindOneTopic(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numberOfPages := numberOfTopicPages(topic)
	pageIndicies := make([]int, numberOfPages)
	for i := 0; i < numberOfPages; i++ {
		pageIndicies[i] = i + 1
	}

	posts, err := model.FindPosts(app.db, id, limitPosts, pageOffset*limitPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.addBreadCrumb("/forum/"+strconv.Itoa(topic.Forum.Id), topic.Forum.Title)
	app.addBreadCrumb("/topic/"+strconv.Itoa(topic.Id), topic.Title)

	results := make(map[string]interface{})
	results["topic"] = topic
	results["posts"] = posts
	results["pageIndicies"] = pageIndicies
	results["currentPage"] = int(pageOffset + 1)

	app.renderTemplate(w, req, "topic", results)
}

func (app *app) handleAddTopic(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	results := make(map[string]interface{})
	results["ForumId"] = id
	app.renderTemplate(w, req, "addTopic", results)
}

func (app *app) handleSaveTopic(w http.ResponseWriter, req *http.Request) {
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
