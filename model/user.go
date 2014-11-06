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
		user.Password[i] = 0;
	}

	return nil
}

func NewUser() *User {
	return &User{-1, "", "", []byte{}, []byte{}, }
}

func SaveUser(db *sql.DB, user *User) error {
	if len(user.PasswordHash) == 0 {
		return errors.New("Password must be hashed.")
	}

	_, err := db.Exec("INSERT INTO users (id, username, email, password_hash) VALUES (NULL,?,?,?)", user.Username, user.Email, user.PasswordHash)
	return err
}
