package main

import (
	"testing"
)

func TestUserPassword(t *testing.T) {
	password := "testpassword"
	u := User{
		Email:    "testEmail@domain.com",
		Password: password,
	}
	u.HashPassword()
	if u.Password == password {
		// Password failed to hash
		t.Fail()
	}
}

func TestUserPasswordCompareCorrectPassword(t *testing.T) {
	password := "testpassword"
	u := User{
		Email:    "testEmail@domain.com",
		Password: password,
	}

	u.HashPassword()
	err := u.CheckPassword(password)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserPasswordCompareIncorrectPassword(t *testing.T) {
	password := "testpassword"
	u := User{
		Email:    "testEmail@domain.com",
		Password: password,
	}
	password = "nope!"
	u.HashPassword()
	err := u.CheckPassword(password)
	if err == nil {
		t.Fatal(err)
	}
}
