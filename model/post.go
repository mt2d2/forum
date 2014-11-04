package model

import "errors"
import "database/sql"
import "time"

type Post struct {
	Id        int
	Text      string
	Published time.Time
	TopicId   int
}

func NewPost() *Post {
	return &Post{-1, "", time.Now().UTC(), -1}
}

func ValidatePost(post *Post) (ok bool, []errors errs) {
	ok = true
	errs = nil


	return
}

func SavePost(db *sql.DB, post *Post) error {
	_, err := db.Exec("INSERT INTO posts (id, text, published, topic_id) VALUES (NULL,?,?,?)", post.Text, post.Published, post.TopicId)
	return err
}

func FindPosts(db *sql.DB, reqId string) ([]Post, error) {
	rows, err := db.Query("SELECT * FROM posts WHERE topic_id=? ORDER BY datetime(published) ASC", reqId)
	defer rows.Close()
	if err != nil {
		return nil, errors.New("could not query for posts for topic " + reqId)
	}

	posts := make([]Post, 0)
	for rows.Next() {
		var (
			id        int
			text      string
			published time.Time
			topicId   int
		)

		err := rows.Scan(&id, &text, &published, &topicId)
		if err != nil {
			return nil, err
		}

		posts = append(posts, Post{id, text, published, topicId})
	}

	return posts, nil
}
