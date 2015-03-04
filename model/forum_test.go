package model

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func mockForum1() *Forum {
	return &Forum{
		Id:          1,
		Title:       "test",
		Description: "tester forum",
		TopicCount:  2,
		PostCount:   12,
	}
}

func mockForum2() *Forum {
	return &Forum{
		Id:          2,
		Title:       "forum zwei",
		Description: "eine Prüfung",
		TopicCount:  2,
		PostCount:   18,
	}
}

func TestFindOneForum(t *testing.T) {
	forum1, err := FindOneForum(db, "1")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(forum1, *mockForum1()) {
		t.Error("forum does not equal mock forum")
	}

	forum2, err := FindOneForum(db, "2")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(forum2, *mockForum2()) {
		t.Errorf("forum does not equal mock forum")
	}
}

func TestFindForums(t *testing.T) {
	forums, err := FindForums(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(forums) != 2 {
		t.Error("expected 2 forums")
	}

	if !reflect.DeepEqual(forums[0], *mockForum1()) {
		t.Error("forum does not equal mock forum")
	}

	if !reflect.DeepEqual(forums[1], *mockForum2()) {
		t.Error("forum does not equal mock forum")
	}
}

func TestMain(m *testing.M) {
	// :memory: databases aren't shared amongst connections
	// https://groups.google.com/forum/#!topic/golang-nuts/AYZl1lNxCfA
	if ldb, err := sql.Open("sqlite3", "file:dummy.db?mode=memory&cache=shared"); err == nil {
		db = ldb
	}
	if _, err := db.Exec(MockupDB); err != nil {
		panic(err)
	}

	ret := m.Run()

	db.Close()
	os.Exit(ret)
}
