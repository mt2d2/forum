package model

import "errors"
import "database/sql"
import "strconv"

type Forum struct {
	Id          int
	Title       string
	Description string
	TopicCount  int
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

	return Forum{id, title, description, 0}, nil
}

func topicCount(db *sql.DB, reqId string) (int, error) {
	var count int

	row := db.QueryRow("SELECT count(*) FROM topics WHERE forum_id = ?", reqId)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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

		topicCount, err := topicCount(db, strconv.Itoa(id))
		if err != nil {
			return nil, errors.New("could not could topics for forum with id " + strconv.Itoa(id))
		}

		forums = append(forums, Forum{id, title, description, topicCount})
	}

	return forums, nil
}
