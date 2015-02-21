package model

import "errors"
import "database/sql"
import "strconv"

type Forum struct {
	Id          int
	Title       string
	Description string

	TopicCount int
	PostCount  int
}

func FindOneForum(db *sql.DB, reqId string) (Forum, error) {
	var (
		id          int
		title       string
		description string
	)

	row := db.QueryRow("SELECT * FROM forums WHERE id = ?", reqId)
	err := row.Scan(&id, &title, &description)
	if err != nil {
		return Forum{}, errors.New("could not query for forum with id " + reqId)
	}

	topicCount, postCount, err := topicAndPostCount(db, reqId)
	if err != nil {
		return Forum{}, err
	}

	return Forum{id, title, description, topicCount, postCount}, nil
}

func topicAndPostCount(db *sql.DB, reqId string) (int, int, error) {
	var topicCount int
	var postCount int

	row := db.QueryRow(`select
												count(distinct topics.id),
	 											count(posts.id)
											from forums
												left join topics on topics.forum_id = forums.id
												left join posts on posts.topic_id = topics.id
											where forums.id = ?`, reqId)

	err := row.Scan(&topicCount, &postCount)
	if err != nil {
		return 0, 0, err
	}

	return topicCount, postCount, nil
}

func FindForums(db *sql.DB) ([]Forum, error) {
	rows, err := db.Query("SELECT * FROM forums")
	defer rows.Close()
	if err != nil {
		return nil, errors.New("could not query for forums")
	}

	forums := make([]Forum, 0)
	for rows.Next() {
		var (
			id          int
			title       string
			description string
		)

		err := rows.Scan(&id, &title, &description)
		if err != nil {
			return nil, errors.New("could not process row")
		}

		topicCount, postCount, err := topicAndPostCount(db, strconv.Itoa(id))
		if err != nil {
			return nil, err
		}

		forums = append(forums, Forum{id, title, description, topicCount, postCount})
	}

	return forums, nil
}
