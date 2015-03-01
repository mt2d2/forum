package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mt2d2/forum/model"
)

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
	}

	numberOfPages := int(math.Ceil(float64(forum.TopicCount) / float64(limitTopics)))
	pageIndicies := make([]int, numberOfPages)
	for i := 0; i < numberOfPages; i++ {
		pageIndicies[i] = i + 1
	}

	topics, err := model.FindTopics(app.db, id, limitTopics, pageOffset*limitTopics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.addBreadCrumb("/forum/"+strconv.Itoa(forum.Id), forum.Title)

	results := make(map[string]interface{})
	results["forum"] = forum
	results["topics"] = topics
	results["pageIndicies"] = pageIndicies
	results["currentPage"] = int(pageOffset + 1)

	app.renderTemplate(w, req, "forum", results)
}