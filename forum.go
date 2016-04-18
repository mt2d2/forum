package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/Schema"
	"github.com/gorilla/mux"
	"github.com/mt2d2/forum/model"
)

func numberOfForumPages(forum *model.Forum) int {
	return int(math.Ceil(float64(forum.TopicCount) / float64(limitTopics)))
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

func (app *app) handleForum(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	pageOffset := 0
	if page, ok := vars["page"]; ok {
		if val, err := strconv.Atoi(page); err == nil {
			pageOffset = val - 1
		}
	}

	forum, err := model.FindOneForum(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numberOfPages := numberOfForumPages(forum)
	pageIndicies := make([]int, numberOfPages)
	for i := 0; i < numberOfPages; i++ {
		pageIndicies[i] = i + 1
	}
	currentPage := int(pageOffset + 1)

	topics, err := model.FindTopics(app.db, id, limitTopics, pageOffset*limitTopics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.addBreadCrumb("/forum/"+strconv.Itoa(forum.Id), forum.Title)
	if currentPage > 1 {
		app.addBreadCrumb("forum/"+strconv.Itoa(forum.Id)+"/page/"+strconv.Itoa(currentPage), "page "+strconv.Itoa(currentPage))
	}

	results := make(map[string]interface{})
	results["forum"] = forum
	results["topics"] = topics
	results["pageIndicies"] = pageIndicies
	results["currentPage"] = currentPage

	app.renderTemplate(w, req, "forum", results)
}

func (app *app) handleAddTopic(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	forum, err := model.FindOneForum(app.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.addBreadCrumb("/forum/"+strconv.Itoa(forum.Id), forum.Title)

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
		return
	}

	http.Redirect(w, req, "/forum/"+req.PostFormValue("ForumId"), http.StatusFound)
}
