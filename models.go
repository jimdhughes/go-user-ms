package main

import "golang.org/x/crypto/bcrypt"

// User structure for accounts
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserSafe is the safe encoding of a User to be used for token grants and data transfer
// The users password should never be leaked beyond this service
type UserSafe struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// HashPassword sets the users password to the bcrypt hashed version of said password
func (u *User) HashPassword() error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(password)
	return nil
}

// CheckPassword compares a plain text password against the Hashed Password of the user
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
