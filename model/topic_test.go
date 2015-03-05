package model

import (
	"math"
	"reflect"
	"strings"
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

func TestSaveTopic(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	newTopic := &Topic{
		Id:          -1,
		Title:       "NEW TITLE",
		Description: "NEW DESC",
		ForumId:     1,
		PostCount:   -1,
		Forum:       nil,
	}

	topicsForum1, err := FindTopics(db, "1", math.MaxInt64, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(topicsForum1) != 3 {
		t.Error("wrong number of topics")
	}

	err = SaveTopic(db, newTopic)
	if err != nil {
		t.Fatal(err)
	}

	newTopicsForum1, err := FindTopics(db, "1", math.MaxInt64, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(newTopicsForum1) != 4 {
		t.Error("wrong number of topics after insert")
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

func TestFindTopicsSmallLimit(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	topicsForum1, err := FindTopics(db, "1", 2, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(topicsForum1) != 2 {
		t.Error("wrong number of topics")
	}

	topicsForum2, err := FindTopics(db, "2", 2, 0)

	if len(topicsForum2) != 2 {
		t.Error("wrong number of topics")
	}
}

func TestFindTopicsSmallOffset(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	topicsForum1, err := FindTopics(db, "1", math.MaxInt64, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(topicsForum1) != 2 {
		t.Error("wrong number of topics")
	}

	if topicsForum1[0].Id != 4 {
		t.Error("first retrieved topic has wrong ID")
	}

	if topicsForum1[1].Id != 5 {
		t.Error("second retireved topics has wrong ID")
	}
}

func TestValidateTopic(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	topic := NewTopic()

	ok, errs := ValidateTopic(db, topic)
	if ok || len(errs) != 2 {
		t.Error("blank topic should not validate")
	}

	topic.Title = "Test"
	ok, errs = ValidateTopic(db, topic)
	if ok || len(errs) != 1 {
		t.Error("still missing a proper forum id")
	}

	topic.ForumId = 255
	ok, errs = ValidateTopic(db, topic)
	if ok || len(errs) != 1 {
		t.Error("invalid forum id")
	}

	topic.ForumId = 1
	ok, errs = ValidateTopic(db, topic)
	if !ok && len(errs) != 0 {
		t.Error("topic should now be valid")
	}

	topic.Title = strings.Repeat("a", 256)
	ok, errs = ValidateTopic(db, topic)
	if ok || len(errs) != 1 {
		t.Error("should not validate long title")
	}

	topic.Title = "Test"
	topic.Description = strings.Repeat("a", 256)
	ok, errs = ValidateTopic(db, topic)
	if ok || len(errs) != 1 {
		t.Error("should not validate long description")
	}

	topic.Description = ""
	topic.Title = strings.Repeat("a", 255)
	ok, errs = ValidateTopic(db, topic)
	if !ok && len(errs) != 0 {
		t.Error("255 title should be ok")
	}

	topic.Description = strings.Repeat("a", 255)
	ok, errs = ValidateTopic(db, topic)
	if !ok && len(errs) != 0 {
		t.Error("255 description should be ok")
	}
}
