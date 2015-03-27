package model

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestEmptyPost(t *testing.T) {
	post := NewPost()
	if !reflect.DeepEqual(post, &Post{-1, "", post.Published, -1, -1, nil}) {
		t.Error("post not empty")
	}
}

func TestValidatePost(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	post := NewPost()

	ok, errs := ValidatePost(db, post)
	if ok || len(errs) != 3 {
		t.Error("blank post should not validate")
	}

	post.UserId = math.MaxInt64
	ok, errs = ValidatePost(db, post)
	if ok || len(errs) != 3 {
		t.Error("this is not a real user")
	}

	post.UserId = 1
	ok, errs = ValidatePost(db, post)
	if ok || len(errs) != 2 {
		t.Error("valid user, still missing topic and text")
	}

	post.Text = "Hello!"
	ok, errs = ValidatePost(db, post)
	if ok || len(errs) != 1 {
		t.Error("still missing topic!")
	}

	post.TopicId = math.MaxInt64
	ok, errs = ValidatePost(db, post)
	if ok || len(errs) != 1 {
		t.Error("should not validate invalid forum")
	}

	post.TopicId = 2
	ok, errs = ValidatePost(db, post)
	if !ok || len(errs) != 0 {
		t.Error("post should now validate")
	}
}

func TestWhitespacePost(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	whitespace := "\t\n\t\n\t\n    \t\n\t\n\t\n"
	post := &Post{1, whitespace, time.Now().UTC(), 1, 1, nil}
	ok, errs := ValidatePost(db, post)
	if ok || len(errs) != 1 {
		t.Error("whitespace is only invalid item")
	}
}
