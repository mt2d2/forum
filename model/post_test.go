package model

import (
	"reflect"
	"testing"
)

func TestEmptyPost(t *testing.T) {
	post := NewPost()
	if !reflect.DeepEqual(post, &Post{-1, "", post.Published, -1, -1, nil}) {
		t.Error("topic not empty")
	}
}
