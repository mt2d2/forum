package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mt2d2/forum/model"
)

func numberOfTopicPages(topic model.Topic) int {
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

	numberOfPages := numberOfTopicPages(*topic)
	pageIndicies := make([]int, numberOfPages)
	for i := 0; i < numberOfPages; i++ {
		pageIndicies[i] = i + 1
	}
	currentPage := int(pageOffset + 1)

	posts, err := model.FindPosts(app.db, id, limitPosts, pageOffset*limitPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.addBreadCrumb("/forum/"+strconv.Itoa(topic.Forum.Id), topic.Forum.Title)
	app.addBreadCrumb("/topic/"+strconv.Itoa(topic.Id), topic.Title)
	if currentPage > 1 {
		app.addBreadCrumb("/topic/"+strconv.Itoa(topic.Id)+"/page/"+strconv.Itoa(currentPage), "page "+strconv.Itoa(currentPage))
	}

	results := make(map[string]interface{})
	results["topic"] = topic
	results["posts"] = posts
	results["pageIndicies"] = pageIndicies
	results["currentPage"] = currentPage

	app.renderTemplate(w, req, "topic", results)
}
