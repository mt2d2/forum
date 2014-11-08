package model

import "errors"
import "database/sql"
import "time"
import "strings"

type Post struct {
	Id        int
	Text      string
	Published time.Time
	TopicId   int
	UserId    int

	// relations
	User      *User
}

func NewPost() *Post {
	return &Post{-1, "", time.Now().UTC(), -1, -1,nil}
}

func ValidatePost(post *Post) (ok bool, errs []error) {
	errs = make([]error, 0)

	if strings.TrimSpace(post.Text) == "" {
		errs = append(errs, errors.New("Post must have some text."))
	}

	// todo, check for valid topic id with database query
	if post.TopicId == -1 {
		errs = append(errs, errors.New("Post must belong to a valid topic."))
	}

	// todo, check for valid user id 
	if post.UserId == -1 {
		errs = append(errs, errors.New("Post must belong to a valid user."))
	}

	return len(errs) == 0, errs
}

func SavePost(db *sql.DB, post *Post) error {
	_, err := db.Exec("INSERT INTO posts (id, text, published, topic_id, user_id) VALUES (NULL,?,?,?,?)", post.Text, post.Published, post.TopicId, post.UserId)
	return err
}

func FindPosts(db *sql.DB, reqId string) ([]Post, error) {
	rows, err := db.Query("SELECT posts.*, users.username FROM posts JOIN users ON posts.user_id = users.id WHERE topic_id=? ORDER BY datetime(published) ASC", reqId)
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
			userId    int
			username  string
		)

		err := rows.Scan(&id, &text, &published, &topicId, &userId, &username)
		if err != nil {
			return nil, err
		}

		posts = append(posts, Post{id, text, published, topicId, userId, 
						&User{-1, username, "", []byte{}, []byte{}}})
	}

	return posts, nil
}
