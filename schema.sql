CREATE TABLE forums(id INTEGER PRIMARY KEY, title varchar(255), description varchar(255));
CREATE TABLE topics(id INTEGER PRIMARY KEY, title varchar(255), description varchar(255), forum_id integer, FOREIGN KEY(forum_id) REFERENCES forum(id));
CREATE TABLE users(id INTEGER PRIMARY KEY, username varchar(255), email varchar(255), password_hash blob);
CREATE TABLE posts(id INTEGER PRIMARY KEY, text TEXT, published TIMESTAMP, topic_id INTEGER, user_id INTEGER, FOREIGN KEY(topic_id) REFERENCES topic(id), FOREIGN KEY(user_id) REFERENCES user(id));
