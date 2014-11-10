package model

import "errors"
import "database/sql"

import "code.google.com/p/go.crypto/bcrypt"

type User struct {
	Id           int
	Username     string
	Email        string
	Password     []byte `schema:"-"`
	PasswordHash []byte `schema:"-"`
}

func (user *User) HashPassword() error {
	if len(user.Password) == 0 {
		return errors.New("User has no password.")
	}

	var err error
	user.PasswordHash, err = bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// clear old password
	for i := range user.Password {
		user.Password[i] = 0
	}

	return nil
}

func (user *User) CompareHashAndPassword(password *[]byte) error {
	if len(user.PasswordHash) == 0 {
		return errors.New("User has no password hash.")
	}

	if len(*password) == 0 {
		return errors.New("No password was given.")
	}

	err := bcrypt.CompareHashAndPassword(user.PasswordHash, *password)
	if err != nil {
		return err
	}

	// clear old password
	for i := range *password {
		(*password)[i] = 0
	}

	return nil
}

func ValidateUser(user *User) (ok bool, errs []error) {
	errs = make([]error, 0)

	if user.Username == "" {
		errs = append(errs, errors.New("Username must not be empty."))
	}

	// todo, check for unique username

	if len(user.Password) == 0 {
		errs = append(errs, errors.New("Password must not be empty."))
	}

	return len(errs) == 0, errs
}

func FindOneUser(db *sql.DB, reqId string) (User, error) {
	var (
		id           int
		username     string
		email        string
		passwordHash []byte
	)

	row := db.QueryRow("SELECT * FROM users WHERE username = ?", reqId)
	err := row.Scan(&id, &username, &email, &passwordHash)
	if err != nil {
		return User{}, errors.New("could not query for user with username " + reqId)
	}

	return User{id, username, email, []byte{}, passwordHash}, nil
}

func FindOneUserById(db *sql.DB, reqId int) (User, error) {
	var (
		id           int
		username     string
		email        string
		passwordHash []byte
	)

	row := db.QueryRow("SELECT * FROM users WHERE id = ?", reqId)
	err := row.Scan(&id, &username, &email, &passwordHash)
	if err != nil {
		return User{}, errors.New("could not query for user with id " + string(reqId))
	}

	return User{id, username, email, []byte{}, passwordHash}, nil
}

func NewUser() *User {
	return &User{-1, "", "", []byte{}, []byte{}}
}

func SaveUser(db *sql.DB, user *User) error {
	if len(user.PasswordHash) == 0 {
		return errors.New("Password must be hashed.")
	}

	_, err := db.Exec("INSERT INTO users (id, username, email, password_hash) VALUES (NULL,?,?,?)", user.Username, user.Email, user.PasswordHash)
	return err
}
