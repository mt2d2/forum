package model

import (
	"math"
	"reflect"
	"testing"
)

func TestEmptyTopic(t *testing.T) {
	topic := NewTopic()
	if !reflect.DeepEqual(topic, &Topic{-1, "", "", -1, -1, nil}) {
		t.Error("topic not empty")
	}
}

func TestFindOneTopic(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	topic, err := FindOneTopic(db, "1")
	if err != nil {
		t.Fatal(err)
	}

	if topic.Id != 1 {
		t.Error("wrong id")
	}

	if topic.Forum == nil {
		t.Error("no forum relation")
	}

	if topic.ForumId != 1 ||
		(topic.Forum != nil && topic.ForumId != topic.Forum.Id) {
		t.Error("wrong forum")
	}

	if topic.Title != "test topic" {
		t.Error("wrong title")
	}

	if topic.Description != "asdf asdf asdf" {
		t.Error("wrong description")
	}

	if topic.PostCount != 11 {
		t.Error("wrong post count")
	}
}

func TestFindTopicsNoLimit(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	topicsForum1, err := FindTopics(db, "1", math.MaxInt64, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(topicsForum1) != 3 {
		t.Error("wrong number of topics")
	}

	topicsForum2, err := FindTopics(db, "2", math.MaxInt64, 0)

	if len(topicsForum2) != 2 {
		t.Error("wrong number of topics")
	}
}
