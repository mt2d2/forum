PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE forums(id INTEGER PRIMARY KEY, title varchar(255), description varchar(255));
INSERT INTO "forums" VALUES(1,'test','tester forum');
INSERT INTO "forums" VALUES(2,'forum zwei','eine Prüfung');
CREATE TABLE topics(id INTEGER PRIMARY KEY, title varchar(255), description varchar(255), forum_id integer, FOREIGN KEY(forum_id) REFERENCES forum(id));
INSERT INTO "topics" VALUES(1,'test topic','asdf asdf asdf',1);
INSERT INTO "topics" VALUES(2,'test topic','for forum 2: asdf asdf asdf',2);
INSERT INTO "topics" VALUES(3,'Aauto add','asdf asdf asdf !',2);
INSERT INTO "topics" VALUES(4,'rawr','just right',1);
CREATE TABLE users(id INTEGER PRIMARY KEY, username varchar(255), email varchar(255), password_hash blob);
INSERT INTO "users" VALUES(1,'test','test',X'24326124313024724573377564694B774B6546694C633349684D656365516C49684D46514E70306A796951784D757731514336374F6E4F476A635175');
INSERT INTO "users" VALUES(2,'tester','test@test.com',X'24326124313024552F31584E5167545054526D37346E456C49514739756B666F796A4B75472E6A554737653458644857334370646B676547516C4A6D');
CREATE TABLE posts(id INTEGER PRIMARY KEY, text TEXT, published TIMESTAMP, topic_id INTEGER, user_id INTEGER, FOREIGN KEY(topic_id) REFERENCES topic(id), FOREIGN KEY(user_id) REFERENCES user(id));
INSERT INTO "posts" VALUES(1,'test','2014-10-31 07:50:55.810912273',1,1);
INSERT INTO "posts" VALUES(2,'test2','2014-10-31 07:52:32.129118657',1,1);
INSERT INTO "posts" VALUES(3,'test3','2014-10-31 07:52:51.073031409',1,1);
INSERT INTO "posts" VALUES(4,'test4','2014-10-31 07:52:55.094815942',1,1);
INSERT INTO "posts" VALUES(5,'Hello, this is a simple test!','2014-10-31 07:54:45.416706454',2,1);
INSERT INTO "posts" VALUES(6,'asdf','2014-11-02 12:07:22.112559999',1,1);
INSERT INTO "posts" VALUES(7,'asdfasdfasdfasfasdfasdfasdfasfasdfasdfasdfasfasdfasdfasdfasfasdfasdfasdfasfasdfasdfasdfasfasdfasdfasdfasfasdfasdfasdfasfasdfasdfasdfasf','2014-11-02 12:07:35.479092061',2,1);
INSERT INTO "posts" VALUES(8,'asdfasdf','2014-11-02 12:07:47.598924176',1,1);
INSERT INTO "posts" VALUES(9,'Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
','2014-11-02 12:08:53.243066099',2,1);
INSERT INTO "posts" VALUES(10,'Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
','2014-11-02 12:31:38.452819455',2,1);
INSERT INTO "posts" VALUES(11,'Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
','2014-11-02 12:31:42.694287309',2,1);
INSERT INTO "posts" VALUES(12,'Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
','2014-11-02 12:31:46.971465623',2,1);
INSERT INTO "posts" VALUES(13,'asdf','2014-11-02 12:46:55.559568931',3,1);
INSERT INTO "posts" VALUES(14,'Another test!','2014-11-02 12:55:33.26669582',3,1);
INSERT INTO "posts" VALUES(15,'asdf','2014-11-03 06:18:08.768825362',1,1);
INSERT INTO "posts" VALUES(16,'asdfasdf','2014-11-03 06:19:30.406177986',1,1);
INSERT INTO "posts" VALUES(17,'blah','2014-11-03 06:20:29.902254599',2,1);
INSERT INTO "posts" VALUES(18,'heller','2014-11-03 06:21:45.20921242',3,1);
INSERT INTO "posts" VALUES(19,'blah blah','2014-11-03 06:22:43.276670489',2,1);
INSERT INTO "posts" VALUES(20,'yus','2014-11-03 06:26:05.636990782',3,1);
INSERT INTO "posts" VALUES(21,'asdf','2014-11-03 06:30:19.665975049',3,1);
INSERT INTO "posts" VALUES(22,'asdf','2014-11-03 06:30:49.605493535',3,1);
INSERT INTO "posts" VALUES(23,'meh
','2014-11-03 06:30:54.474603291',3,1);
INSERT INTO "posts" VALUES(24,'flash!','2014-11-03 06:34:56.892273745',1,1);
INSERT INTO "posts" VALUES(25,'flash for real!','2014-11-03 06:35:19.239830122',1,1);
INSERT INTO "posts" VALUES(26,'no really, flash','2014-11-03 06:36:30.634986366',1,1);
INSERT INTO "posts" VALUES(27,'how about another post?','2014-11-03 06:36:49.395334151',3,1);
INSERT INTO "posts" VALUES(28,'Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
','2014-11-03 06:39:00.954005023',2,1);
INSERT INTO "posts" VALUES(29,'','2014-11-04 05:56:58.608376074',2,1);
INSERT INTO "posts" VALUES(30,'blah blah blah blah','2014-11-04 06:08:47.772019858',4,1);
COMMIT;
